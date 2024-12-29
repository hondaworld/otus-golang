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
	First  *ListItem
	Last   *ListItem
	Length int
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	listItem := ListItem{Value: v}

	if l.Len() == 0 {
		l.Last = &listItem
	} else {
		listItem.Next = l.First
		l.First.Prev = &listItem
	}

	l.First = &listItem
	l.Length++

	return &listItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	listItem := ListItem{Value: v}

	if l.Len() == 0 {
		l.First = &listItem
	} else {
		listItem.Prev = l.Last
		l.Last.Next = &listItem
	}

	l.Last = &listItem
	l.Length++

	return &listItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.First = i.Next

		if l.First != nil {
			l.First.Prev = nil
		}
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.Last = i.Prev

		if l.Last != nil {
			l.Last.Next = nil
		}
	}

	l.Length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev != nil {
		if i.Next != nil {
			i.Next.Prev = i.Prev
		} else {
			l.Last = i.Prev
			l.Last.Next = nil
		}

		l.First.Prev = i
		i.Prev = nil
		i.Next = l.First
		l.First = i
	}
}

func (l *list) Len() int {
	return l.Length
}

func NewList() List {
	return new(list)
}
