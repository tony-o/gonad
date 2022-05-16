package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type EITHER int

const (
	LEFT EITHER = iota
	RIGHT
)

var EITHERMAP map[EITHER]string = map[EITHER]string{
	LEFT:  "Left",
	RIGHT: "Right",
}

type Either[L any, R any] struct {
	left     L
	right    R
	is_right bool
}

type EitherIfc[L any, R any] interface {
	Which() EITHER
	Right() R
	Left() L
}

func Left[L any, R any](left L) Either[L, R] {
	return Either[L, R]{left: left, is_right: false}
}

func Right[L any, R any](right R) Either[L, R] {
	return Either[L, R]{right: right, is_right: true}
}

func (e Either[L, R]) Which() EITHER {
	if e.is_right {
		return RIGHT
	}
	return LEFT
}
func (e Either[L, R]) Left() L {
	return e.left
}

func (e Either[L, R]) Right() R {
	return e.right
}

func Transform[L, M, R any](fn1 func(L) Either[L, R], fn2 func(Either[L, R]) Either[M, R]) func(L) Either[M, R] {
	return func(l L) Either[M, R] {
		e := fn1(l)
		if e.Which() == LEFT {
			return fn2(e)
		}
		return Right[M, R](e.Right())
	}
}

func TransformWith1e1[L, M, R, A any](fn1 func(A) Either[L, R], fn2 func(L, A) Either[M, R]) func(_, _ A) Either[M, R] {
	return func(a1, a2 A) Either[M, R] {
		e := fn1(a1)
		if e.Which() == LEFT {
			return fn2(e.Left(), a2)
		}
		return Right[M, R](e.Right())
	}
}
func Map[L, R any](rightHalt bool, fns ...func(Either[L, R]) Either[L, R]) func(Either[L, R]) Either[L, R] {
	return func(e Either[L, R]) Either[L, R] {
		if e.Which() == LEFT {
			for _, fn := range fns {
				v1 := fn(e)
				if v1.Which() == RIGHT && rightHalt {
					return Right[L, R](v1.Right())
				}
			}
		}
		return Right[L, R](e.Right())
	}
}

type FilesOperator Either[[]*os.File, error]

func (fo FilesOperator) Which() EITHER {
	return Either[[]*os.File, error](fo).Which()
}
func (fo FilesOperator) Left() []*os.File {
	return Either[[]*os.File, error](fo).Left()
}
func (fo FilesOperator) Right() error {
	return Either[[]*os.File, error](fo).Right()
}

var _ EitherIfc[[]*os.File, error] = FilesOperator{}

func OpenFiles(ss ...string) FilesOperator {
	fs := []*os.File{}
	for _, s := range ss {
		r, err := os.Open(s)
		if err != nil {
			for _, c := range fs {
				c.Close()
			}
			return FilesOperator(Right[[]*os.File, error](err))
		}
		fs = append(fs, r)
	}
	return FilesOperator(Left[[]*os.File, error](fs))
}

func (fo FilesOperator) CloseFiles() Either[int, error] {
	if fo.Which() == RIGHT {
		return Right[int, error](fo.Right())
	}
	fs := fo.Left()
	for _, f := range fs {
		err := f.Close()
		if err != nil {
			return Right[int, error](errors.New(fmt.Sprint(err)))
		}
	}
	return Left[int, error](0)
}

func (fo FilesOperator) CopyFiles() FilesOperator {
	if fo.Which() == RIGHT {
		return fo
	}
	fs := fo.Left()
	if len(fs)%2 != 0 {
		return FilesOperator(Right[[]*os.File, error](errors.New("Invalid copy length")))
	}
	for i := 0; i < len(fs); i += 2 {
		src := fs[i]
		dest := fs[i+1]
		_, err := io.Copy(src, dest)
		if err != nil {
			return FilesOperator(Right[[]*os.File, error](err))
		}
	}
	return fo
}

func main() {
	tests := [][]string{
		{"/tmp/x", "/tmp/y"},
		{"/tmp/x"},
		{"/tmp/doesnt-exist"},
	}
	for _, t := range tests {
		fmt.Printf("t=%v\n", t)
		result := OpenFiles(t...).CopyFiles().CloseFiles()
		if result.Which() == LEFT {
			fmt.Printf("SUCCESS=%v\n", result.Left())
		} else {
			fmt.Printf("FAILURE=%v\n", result.Right())
		}
	}
}
