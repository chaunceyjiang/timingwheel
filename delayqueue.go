package timingwheel

import "container/heap"

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



