package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
    listLen int
    firstItem *ListItem
    lastItem *ListItem
}

func (l *list) Len() int {
    return l.listLen
}

func (l *list) Front() *ListItem {
    if l.listLen > 0 {
        return l.firstItem
    }
    return nil
}

func (l *list) Back() *ListItem {
    if l.listLen > 0 {
        return l.lastItem
    }
    return nil
}

func (l *list) MoveToFront(item *ListItem) {
    if (item == nil) || (item == l.firstItem) {
        return
    }

    if item.Prev != nil {
        item.Prev.Next = item.Next
    }
    if item.Next != nil {
        item.Next.Prev = item.Prev
    }
    if item == l.lastItem {
        l.lastItem = item.Prev
    }
    if l.firstItem != nil {
        l.firstItem.Prev = item
    }
    item.Prev = nil
    item.Next = l.firstItem
    l.firstItem = item
}

func (l *list) PushBack(v interface{}) *ListItem {
    item := new(ListItem)
    item.Value = v

    if l.lastItem == nil {
        l.lastItem = item
        l.firstItem = item
        item.Next = nil
        item.Prev = nil
    } else {
        l.lastItem.Next = item
        item.Prev = l.lastItem
        l.lastItem = item
    }
    l.listLen++
    return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
    item := new(ListItem)
    item.Value = v
    if l.firstItem == nil {
        l.lastItem = item
        l.firstItem = item
        item.Next = nil
        item.Prev = nil
    } else {
        l.firstItem.Prev = item
        item.Next = l.firstItem
        l.firstItem = item
    }
    l.listLen++
    return l.firstItem
}

func (l *list) Remove(item *ListItem) {
    if item == nil {
        return
    }
    if item.Prev != nil {
        item.Prev.Next = item.Next
    }
    if item.Next != nil {
        item.Next.Prev = item.Prev
    }
    if item == l.firstItem {
        l.firstItem = item.Next
    }
    if item == l.lastItem {
        l.lastItem = item.Prev
    }
    l.listLen--
    return
}

func NewList() List {
	return new(list)
}
