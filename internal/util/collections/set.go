package collections

type Void struct{}

type Set[T comparable] map[T]Void

func NewSet[T comparable](items ...T) Set[T] {
	set := make(Set[T], len(items))
	for _, v := range items {
		set[v] = Void{}
	}

	return set
}

func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Add(item T) {
	s[item] = Void{}
}

// SetFromResult wraps tuple of values slice and error into set.
func SetFromResult[T comparable](result []T, err error) (Set[T], error) {
	if err != nil {
		return nil, err
	}

	return NewSet(result...), nil
}
