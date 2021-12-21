package dagg

import (
	"fmt"
)

// Edge represents an edge in the graph, with a source and target vertex.
type Edge[T Hashable] interface {
	Source() T
	Target() T

	Hashable
}

// BasicEdge returns an Edge implementation that simply tracks the source
// and target given as-is.
func BasicEdge[T Hashable](source, target T) Edge[T] {
	return &basicEdge[T]{Src: source, Trgt: target}
}

// basicEdge is a basic implementation of Edge that has the source and
// target vertex.
type basicEdge[T Hashable] struct {
	Src, Trgt T
}

func (e *basicEdge[T]) Hashcode() string {
	return fmt.Sprintf("%s-%s", e.Src.Hashcode(), e.Trgt.Hashcode())
}

func (e *basicEdge[T]) Source() T {
	return e.Src
}

func (e *basicEdge[T]) Target() T {
	return e.Trgt
}
