package storage

import (
	"bytes"
	"container/heap"
)

type MultiIterator struct {
	start []byte
	end   []byte
	// Acts as priority queue based on sort order of current key in each iterator
	iterators     []KVIterator
	iteratorOrder map[KVIterator]int
	lessComp      int
}

// MultiIterator iterates in order over a series o
func NewMultiIterator(reverse bool, iterators ...KVIterator) *MultiIterator {
	// reuse backing array
	lessComp := -1
	if reverse {
		lessComp = 1
	}
	mi := &MultiIterator{
		iterators:     iterators,
		iteratorOrder: make(map[KVIterator]int),
		lessComp:      lessComp,
	}
	mi.init()
	return mi
}

func (mi *MultiIterator) init() {
	validIterators := mi.iterators[:0]
	for i, it := range mi.iterators {
		mi.iteratorOrder[it] = i
		if it.Valid() {
			validIterators = append(validIterators, it)
			start, end := it.Domain()
			if mi.start == nil || bytes.Compare(start, mi.start) == mi.lessComp {
				mi.start = start
			}
			if mi.end == nil || bytes.Compare(mi.end, end) == mi.lessComp {
				mi.end = end
			}
		} else {
			// Not clear if this is necessary - fairly sure it is permitted so can't hurt
			it.Close()
		}
	}
	mi.iterators = validIterators
	heap.Init(mi)
}

// sort.Interface implementation
func (mi *MultiIterator) Len() int {
	return len(mi.iterators)
}

func (mi *MultiIterator) Less(i, j int) bool {
	comp := bytes.Compare(mi.iterators[i].Key(), mi.iterators[j].Key())
	return comp == mi.lessComp || (comp == 0 && mi.iteratorOrder[mi.iterators[i]] < mi.iteratorOrder[mi.iterators[j]])
}

func (mi *MultiIterator) Swap(i, j int) {
	mi.iterators[i], mi.iterators[j] = mi.iterators[j], mi.iterators[i]
}

func (mi *MultiIterator) Push(x interface{}) {
	mi.iterators = append(mi.iterators, x.(KVIterator))
}

func (mi *MultiIterator) Pop() interface{} {
	n := len(mi.iterators) - 1
	it := mi.iterators[n]
	mi.iterators = mi.iterators[:n]
	return it
}

func (mi *MultiIterator) Domain() ([]byte, []byte) {
	return mi.start, mi.end
}

func (mi *MultiIterator) Valid() bool {
	return len(mi.iterators) > 0
}

func (mi *MultiIterator) Next() {
	// Always advance the lowest iterator - the same one we serve the KV pair from
	it := heap.Pop(mi).(KVIterator)
	it.Next()
	if it.Valid() {
		heap.Push(mi, it)
	}
}

func (mi *MultiIterator) Key() []byte {
	return mi.Peek().Key()
}

func (mi *MultiIterator) Value() []byte {
	return mi.Peek().Value()
}

func (mi *MultiIterator) Peek() KVIterator {
	return mi.iterators[0]
}

func (mi *MultiIterator) Close() {
	// Close any remaining valid iterators
	for _, it := range mi.iterators {
		it.Close()
	}
}