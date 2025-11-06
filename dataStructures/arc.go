package datastructures

type Arc[T any] struct {
	val     T
	getChan chan chan T
	setChan chan T
}

func NewArc[T any](initial T) *Arc[T] {
	arc := &Arc[T]{
		val:     initial,
		getChan: make(chan chan T),
		setChan: make(chan T),
	}

	go arc.manageValue()
	return arc
}

func (a *Arc[T]) manageValue() {
	for {
		select {
		case getRequest := <-a.getChan:
			getRequest <- a.val
		case newVal := <-a.setChan:
			a.val = newVal
		}
	}
}

func (a *Arc[T]) Get() T {
	respChan := make(chan T)
	a.getChan <- respChan
	return <-respChan
}

func (a *Arc[T]) Set(newVal T) {
	a.setChan <- newVal
}
