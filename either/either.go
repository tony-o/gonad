package either

type WHICH int64

const (
	LEFT WHICH = iota
	RIGHT
)

type Either[L, R any] struct {
	isLeft bool
	left   L
	right  R
}

type EitherInterface[L, R any] interface {
	Which() WHICH
	Left() L
	Right() R
}

func Left[L, R any](l L) Either[L, R] {
	return Either[L, R]{
		isLeft: true,
		left:   l,
	}
}
func Right[L, R any](r R) Either[L, R] {
	return Either[L, R]{
		right: r,
	}
}

func (e Either[L, R]) Which() WHICH {
	if e.isLeft {
		return LEFT
	}
	return RIGHT
}

func (e Either[L, R]) Left() L {
	return e.left
}

func (e Either[L, R]) Right() R {
	return e.right
}
