package dag

import (
	"testing"
)

func TestBasicEdgeHashcode(t *testing.T) {
	e1 := BasicEdge(myint(1), myint(2))
	e2 := BasicEdge(myint(1), myint(2))
	if e1.Hashcode() != e2.Hashcode() {
		t.Fatalf("bad")
	}
}

type test struct {
	Value string
}

func (t test) Hashcode() string {
	return t.Value
}
func TestBasicEdgeHashcode_pointer(t *testing.T) {

	v1, v2 := &test{"foo"}, &test{"bar"}
	e1 := BasicEdge(v1, v2)
	e2 := BasicEdge(v1, v2)
	if e1.Hashcode() != e2.Hashcode() {
		t.Fatalf("bad")
	}
}
