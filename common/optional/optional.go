package optional

type Optional[T any] struct {
	value    T
	hasValue bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value:    value,
		hasValue: true,
	}
}

func None[T any]() Optional[T] {
	return Optional[T]{hasValue: false}
}

func (o Optional[T]) Get() T {
	return o.value
}

func (o Optional[T]) IsEmpty() bool {
	return !o.hasValue
}

func (o Optional[T]) IsPresent() bool {
	return o.hasValue
}
