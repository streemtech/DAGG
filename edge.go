package dag

import (
	"fmt"
)

// Edge represents an edge in the graph, with a source and target vertex.
type Edge[T Hashable] interface {
	Source() Vertex[T]
	Target() Vertex[T]

	Hashable
}

// BasicEdge returns an Edge implementation that simply tracks the source
// and target given as-is.
func BasicEdge[T Hashable](source, target Vertex[T]) Edge[T] {
	return &basicEdge[T]{S: source, T: target}
}

// basicEdge is a basic implementation of Edge that has the source and
// target vertex.
type basicEdge[T Hashable] struct {
	S, T Vertex[T]
}

func (e *basicEdge[T]) Hashcode() string {
	return fmt.Sprintf("%s-%s", e.S.Hashcode(), e.T.Hashcode())
}

func (e *basicEdge[T]) Source() Vertex[T] {
	return e.S
}

func (e *basicEdge[T]) Target() Vertex[T] {
	return e.T
}
