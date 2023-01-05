package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first *ListItem
	last  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) insertAfter(node *ListItem, newNode *ListItem) {
	newNode.Prev = node
	if node.Next == nil {
		newNode.Next = nil
		l.last = newNode
	} else {
		newNode.Next = node.Next
		node.Next.Prev = newNode
	}
	node.Next = newNode
}

func (l *list) insertBefore(node *ListItem, newNode *ListItem) {
	newNode.Next = node
	if node.Prev == nil {
		newNode.Prev = nil
		l.first = newNode
	} else {
		newNode.Prev = node.Prev
		node.Prev.Next = newNode
	}
	node.Prev = newNode
}

func (l *list) insertFront(node *ListItem) {
	if l.Front() == nil {
		l.first = node
		l.last = node
		node.Prev = nil
		node.Next = nil
	} else {
		l.insertBefore(l.first, node)
	}
}

func (l *list) PushFront(v interface{}) *ListItem {
	var node ListItem
	node.Value = v
	l.insertFront(&node)
	l.len++
	return &node
}

func (l *list) PushBack(v interface{}) *ListItem {
	var node ListItem
	node.Value = v

	if l.last == nil {
		l.insertFront(&node)
	} else {
		l.insertAfter(l.last, &node)
	}

	l.len++
	return &node
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.first = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.insertFront(i)
	l.len++
}

func (l *list) Clear() {
	l.len = 0
	l.first = nil
	l.last = nil
}

func NewList() List {
	return new(list)
}
