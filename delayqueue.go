package timingwheel

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

type item struct {
	value    interface{}
	priority int64
	index    int64
}

// priorityQueue 优先级队列
type priorityQueue []*item

func newPriorityQueue(size int64) priorityQueue {
	return make(priorityQueue, 0, size)

}

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = int64(i)
	pq[j].index = int64(j)
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	c := cap(*pq)

	// 如果新加的元素已经超过了初始化数组的大小,那么进行扩容
	if n+1 > c {
		npq := make(priorityQueue, n, 2*c)
		copy(npq, *pq)
		*pq = npq // 将当前pq 指向新的list
	}

	*pq = (*pq)[0 : n+1]
	i := x.(*item)
	i.index = int64(n)
	(*pq)[n] = i

}

func (pq *priorityQueue) Pop() interface{} {
	n := len(*pq)
	c := cap(*pq)

	// 移除元素后,元素的长度已经小于容器的一半,释放后不用的内存
	if n < (c/2) && c > 25 {
		npq := make(priorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq // 将当前pq 指向新的list
	}
	i := (*pq)[n-1]
	i.index = -1
	*pq = (*pq)[0 : n-1]
	return i
}

func (pq *priorityQueue) PeekAndShift(max int64) (*item, int64) {

	if pq.Len() == 0 {
		return nil, 0
	}

	e := (*pq)[0]

	if e.priority > max {
		return nil, e.priority - max // 还没有到预期的值,这里返回 delta
	}
	heap.Remove(pq, 0)
	return e, 0
}

// delayQueue 延迟队列
type delayQueue struct {
	C  chan interface{} // 用来发送数据
	mu sync.Mutex
	pq priorityQueue // 队列

	pending int32 // 标记 是否 在等待

	wakeupC chan struct{} // 用来唤醒当前队列,来处理数据

}

func newDelayQueue(size int64) *delayQueue {
	return &delayQueue{
		C:       make(chan interface{}),
		pq:      newPriorityQueue(size),
		wakeupC: make(chan struct{}),
	}
}

// offer 放入数据
func (dq *delayQueue) offer(elem interface{}, expiration int64) {
	e := &item{
		value:    elem,
		priority: expiration,
	}

	dq.mu.Lock()
	heap.Push(&dq.pq, e)
	dq.mu.Unlock()

	//
	if atomic.CompareAndSwapInt32(&dq.pending, 1, 0) {
		// 唤醒队列
		dq.wakeupC <- struct{}{}
	}

}

func (dq *delayQueue) poll(exitC chan struct{}, nowF func() int64) {
	for {
		// 获取当前时间
		now := nowF()

		dq.mu.Lock()
		e, delta := dq.pq.PeekAndShift(now)
		dq.mu.Unlock()
		if e == nil {
			// 当前 队列没有到期,则阻塞队列,只有被唤醒,标记阻塞
			atomic.StoreInt32(&dq.pending, 1)
		}
		// 获取元素失败,则阻塞当前操作,直到被唤醒
		if e == nil {
			if delta == 0 {
				// 队列没有任何元素
				select {
				case <-dq.wakeupC: // 阻塞,直到被唤醒
					continue
				case <-exitC:
					goto exit
				}
			} else if delta > 0 {
				select {
				case <-dq.wakeupC:
					continue
				case <-exitC:
					goto exit
				// 延迟 delta 时间 然后处理 唤醒 队列
				case <-time.After(time.Duration(delta) * time.Millisecond):

					if atomic.SwapInt32(&dq.pending, 0) == 0 {
						// 为什么这里等于0 因为如果新加的元素 第一次poll 没有超时, 而没有进行第二次poll 时,那么这里应该阻塞,直到被 poll
						//FIXME 这里也可以直接发送出去,而不用继续continue
						<-dq.wakeupC
					}

					continue // 被唤醒后,去队列中取值
				}
			}

		}
		select {
		case dq.C <- e.value: // 将当到期元素,发送出去
		case <-exitC:
			goto exit

		}
	}

exit:
	atomic.StoreInt32(&dq.pending, 0)
}
