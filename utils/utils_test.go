package utils

import (
	"testing"
)

type FooBar struct {
	Foo []int
}

func TestDeepCopy(t *testing.T) {
	src := &FooBar{
		Foo: make([]int, 1),
	}

	src.Foo = append(src.Foo, 1)

	var dst FooBar
	if err := DeepCopy(src, &dst); err != nil {
		t.Error(err)
	}

	src.Foo[0] = 2
	src.Foo = append(src.Foo, 2)

	if len(src.Foo) == len(dst.Foo) || src.Foo[0] == dst.Foo[0] {
		t.Error("Deep copy failed")
	}
}
