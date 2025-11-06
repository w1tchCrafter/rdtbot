package datastructures

type Queue[T any] struct {
	data []T
}

func (q *Queue[T]) Push(data T) {
	q.data = append(q.data, data)
}

func (q *Queue[T]) Pop() T {
	var ret T

	if len(q.data) > 0 {
		ret = q.data[0]
	}

	if len(q.data) >= 1 {
		q.data = q.data[1:]
	}

	return ret
}
