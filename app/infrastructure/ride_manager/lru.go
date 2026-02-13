package ride_manager

type LRUCache struct {
	store map[int]*node
	head  *node
	tail  *node
}

type node struct {
	Prev  *node
	Next  *node
	Value int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		store: make(map[int]*node, 3100),
		head:  nil,
		tail:  nil,
	}
}

func (this *LRUCache) Get(key int) int {
	newHead := this.store[key]

	// change the position
	this.changeHead(newHead)

	return newHead.Value
}

func (this *LRUCache) Put(key int, value int) {
	newHead := &node{
		Prev: nil,
		Next: nil,
	}

	// this mean both head & tail is nil
	if this.head == nil {
		this.head = newHead
		this.tail = newHead

		return
	}

	this.changeHead(newHead)
}

func (this *LRUCache) changeHead(newHead *node) {
	prevHead := this.head
	this.head = newHead
	newHead.Prev = prevHead
	prevHead.Next = newHead

	newHead.Next = nil
}

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */
