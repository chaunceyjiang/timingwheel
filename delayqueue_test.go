package timingwheel

import (
	"container/heap"
	"fmt"
	"testing"
)

func Test_newPriorityQueue(t *testing.T) {

	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := newPriorityQueue(10)
	i := 0
	for value, priority := range items {
		pq = append(pq,&item{
			value:    value,
			priority: int64(priority),
			index:    int64(i),
		})
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	e := &item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, e)


	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		e := heap.Pop(&pq).(*item)
		fmt.Printf("%.2d:%s ", e.priority, e.value)
	}

}