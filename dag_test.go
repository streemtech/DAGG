package dag

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestAcyclicGraphRoot(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Connect(BasicEdge[myint](myint(3), myint(2)))
	g.Connect(BasicEdge[myint](myint(3), myint(1)))

	if root, err := g.Root(); err != nil {
		t.Fatalf("err: %s", err)
	} else if root != myint(3) {
		t.Fatalf("bad: %#v", root)
	}
}

func TestAcyclicGraphRoot_cycle(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Connect(BasicEdge[myint](myint(1), myint(2)))
	g.Connect(BasicEdge[myint](myint(2), myint(3)))
	g.Connect(BasicEdge[myint](myint(3), myint(1)))

	if _, err := g.Root(); err == nil {
		t.Fatal("should error")
	}
}

func TestAcyclicGraphRoot_multiple(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Connect(BasicEdge[myint](myint(3), myint(2)))

	if _, err := g.Root(); err == nil {
		t.Fatal("should error")
	}
}

func TestAcyclicGraphTransReduction(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Connect(BasicEdge[myint](myint(1), myint(2)))
	g.Connect(BasicEdge[myint](myint(1), myint(3)))
	g.Connect(BasicEdge[myint](myint(2), myint(3)))
	g.TransitiveReduction()

	actual := strings.TrimSpace(g.String())
	expected := strings.TrimSpace(testGraphTransReductionStr)
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestAcyclicGraphTransReduction_more(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Add(myint(4))
	g.Connect(BasicEdge[myint](myint(1), myint(2)))
	g.Connect(BasicEdge[myint](myint(1), myint(3)))
	g.Connect(BasicEdge[myint](myint(1), myint(4)))
	g.Connect(BasicEdge[myint](myint(2), myint(3)))
	g.Connect(BasicEdge[myint](myint(2), myint(4)))
	g.Connect(BasicEdge[myint](myint(3), myint(4)))
	g.TransitiveReduction()

	actual := strings.TrimSpace(g.String())
	expected := strings.TrimSpace(testGraphTransReductionMoreStr)
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

// use this to simulate slow sort operations
type counter struct {
	Name  string
	Calls int64
}

func (s *counter) Hashcode() string {
	return s.Name
}

func (s *counter) String() string {
	s.Calls++
	return s.Name
}

// Make sure we can reduce a sizable, fully-connected graph.
func TestAcyclicGraphTransReduction_fullyConnected(t *testing.T) {
	var g AcyclicGraph[*counter]

	const nodeCount = 200
	nodes := make([]*counter, nodeCount)
	for i := 0; i < nodeCount; i++ {
		nodes[i] = &counter{Name: strconv.Itoa(i)}
	}

	// Add them all to the graph
	for _, n := range nodes {
		g.Add(n)
	}

	// connect them all
	for i := range nodes {
		for j := range nodes {
			if i == j {
				continue
			}
			g.Connect(BasicEdge[*counter](nodes[i], nodes[j]))
		}
	}

	g.TransitiveReduction()

	vertexNameCalls := int64(0)
	for _, n := range nodes {
		vertexNameCalls += n.Calls
	}

	switch {
	case vertexNameCalls > 2*nodeCount:
		// Make calling it more the 2x per node fatal.
		// If we were sorting this would give us roughly ln(n)(n^3) calls, or
		// >59000000 calls for 200 vertices.
		t.Fatalf("VertexName called %d times", vertexNameCalls)
	case vertexNameCalls > 0:
		// we don't expect any calls, but a change here isn't necessarily fatal
		t.Logf("WARNING: VertexName called %d times", vertexNameCalls)
	}
}

func TestAcyclicGraphValidate(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Connect(BasicEdge[myint](myint(3), myint(2)))
	g.Connect(BasicEdge[myint](myint(3), myint(1)))

	if err := g.Validate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestAcyclicGraphValidate_cycle(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Connect(BasicEdge[myint](myint(3), myint(2)))
	g.Connect(BasicEdge[myint](myint(3), myint(1)))
	g.Connect(BasicEdge[myint](myint(1), myint(2)))
	g.Connect(BasicEdge[myint](myint(2), myint(1)))

	if err := g.Validate(); err == nil {
		t.Fatal("should error")
	}
}

func TestAcyclicGraphValidate_cycleSelf(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Connect(BasicEdge[myint](myint(1), myint(1)))

	if err := g.Validate(); err == nil {
		t.Fatal("should error")
	}
}

func TestAcyclicGraphAncestors(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Add(myint(4))
	g.Add(myint(5))
	g.Connect(BasicEdge[myint](myint(0), myint(1)))
	g.Connect(BasicEdge[myint](myint(1), myint(2)))
	g.Connect(BasicEdge[myint](myint(2), myint(3)))
	g.Connect(BasicEdge[myint](myint(3), myint(4)))
	g.Connect(BasicEdge[myint](myint(4), myint(5)))

	actual, err := g.Ancestors(myint(2))
	if err != nil {
		t.Fatalf("err: %#v", err)
	}

	expected := []Vertex[myint]{myint(3), myint(4), myint(5)}

	if actual.Len() != len(expected) {
		t.Fatalf("bad length! expected %#v to have len %d", actual, len(expected))
	}

	for _, e := range expected {
		if !actual.Include(e) {
			t.Fatalf("expected: %#v to include: %#v", expected, actual)
		}
	}
}

func TestAcyclicGraphDescendents(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Add(myint(4))
	g.Add(myint(5))
	g.Connect(BasicEdge[myint](myint(0), myint(1)))
	g.Connect(BasicEdge[myint](myint(1), myint(2)))
	g.Connect(BasicEdge[myint](myint(2), myint(3)))
	g.Connect(BasicEdge[myint](myint(3), myint(4)))
	g.Connect(BasicEdge[myint](myint(4), myint(5)))

	actual, err := g.Descendents(myint(2))
	if err != nil {
		t.Fatalf("err: %#v", err)
	}

	expected := []Vertex[myint]{myint(0), myint(1)}

	if actual.Len() != len(expected) {
		t.Fatalf("bad length! expected %#v to have len %d", actual, len(expected))
	}

	for _, e := range expected {
		if !actual.Include(e) {
			t.Fatalf("expected: %#v to include: %#v", expected, actual)
		}
	}
}

// func TestAcyclicGraphWalk(t *testing.T) {
// 	var g AcyclicGraph[myint]
// 	g.Add(myint(1))
// 	g.Add(myint(2))
// 	g.Add(myint(3))
// 	g.Connect(BasicEdge[myint](myint(3), myint(2)))
// 	g.Connect(BasicEdge[myint](myint(3), myint(1)))

// 	var visits []Vertex[myint]
// 	var lock sync.Mutex
// 	g.Walk(func(v Vertex[myint]) {
// 		lock.Lock()
// 		defer lock.Unlock()
// 		visits = append(visits, v)
// 	})
// 	t.Log("test")

// 	expected := [][]Vertex[myint]{
// 		{myint(1), myint(2), myint(3)},
// 		{myint(2), myint(1), myint(3)},
// 	}
// 	time.Sleep(time.Second)
// 	for _, e := range expected {
// 		if reflect.DeepEqual(visits, e) {
// 			return
// 		}
// 	}

// 	t.Fatalf("bad: %#v", visits)
// }

// func TestAcyclicGraphWalk_error(t *testing.T) {
// 	var g AcyclicGraph[myint]
// 	g.Add(myint(1))
// 	g.Add(myint(2))
// 	g.Add(myint(3))
// 	g.Add(myint(4))
// 	g.Connect(BasicEdge[myint](myint(4), myint(3)))
// 	g.Connect(BasicEdge[myint](myint(3), myint(2)))
// 	g.Connect(BasicEdge[myint](myint(2), myint(1)))

// 	var visits []Vertex[myint]
// 	var lock sync.Mutex
// 	g.Walk(func(v Vertex[myint]) {
// 		lock.Lock()
// 		defer lock.Unlock()

// 		var diags Diagnostics

// 		if v == 2 {
// 			diags = diags.Append(fmt.Errorf("error"))
// 			return diags
// 		}

// 		visits = append(visits, v)
// 		return diags
// 	})
// 	if err == nil {
// 		t.Fatal("should error")
// 	}

// 	expected := []Vertex{1}
// 	if !reflect.DeepEqual(visits, expected) {
// 		t.Errorf("wrong visits\ngot:  %#v\nwant: %#v", visits, expected)
// 	}

// }

// func BenchmarkDAG(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		count := 150
// 		b.StopTimer()
// 		g := &AcyclicGraph[mystr]{}

// 		// create 4 layers of fully connected nodes
// 		// layer A
// 		for i := 0; i < count; i++ {
// 			g.Add(mystr(fmt.Sprintf("A%d", i)))
// 		}

// 		// layer B
// 		for i := 0; i < count; i++ {
// 			B := fmt.Sprintf("B%d", i)
// 			g.Add(mystr(B))
// 			for j := 0; j < count; j++ {
// 				g.Connect(BasicEdge[mystr](mystr(B), mystr(fmt.Sprintf("A%d", j))))
// 			}
// 		}

// 		// layer C
// 		for i := 0; i < count; i++ {
// 			c := fmt.Sprintf("C%d", i)
// 			g.Add(mystr(c))
// 			for j := 0; j < count; j++ {
// 				// connect them to previous layers so we have something that requires reduction
// 				g.Connect(BasicEdge[mystr](mystr(c), mystr(fmt.Sprintf("A%d", j))))
// 				g.Connect(BasicEdge[mystr](mystr(c), mystr(fmt.Sprintf("B%d", j))))
// 			}
// 		}

// 		// layer D
// 		for i := 0; i < count; i++ {
// 			d := fmt.Sprintf("D%d", i)
// 			g.Add(mystr(d))
// 			for j := 0; j < count; j++ {
// 				g.Connect(BasicEdge[mystr](mystr(d), mystr(fmt.Sprintf("A%d", j))))
// 				g.Connect(BasicEdge[mystr](mystr(d), mystr(fmt.Sprintf("B%d", j))))
// 				g.Connect(BasicEdge[mystr](mystr(d), mystr(fmt.Sprintf("C%d", j))))
// 			}
// 		}

// 		b.StartTimer()
// 		// Find dependencies for every node
// 		for _, v := range g.Vertices() {
// 			_, err := g.Ancestors(v)
// 			if err != nil {
// 				b.Fatal(err)
// 			}
// 		}

// 		// reduce the final graph
// 		g.TransitiveReduction()
// 	}
// }

func TestAcyclicGraph_ReverseDepthFirstWalk_WithRemoval(t *testing.T) {
	var g AcyclicGraph[myint]
	g.Add(myint(1))
	g.Add(myint(2))
	g.Add(myint(3))
	g.Connect(BasicEdge[myint](myint(3), myint(2)))
	g.Connect(BasicEdge[myint](myint(2), myint(1)))

	var visits []Vertex[myint]
	var lock sync.Mutex
	err := g.SortedReverseDepthFirstWalk([]Vertex[myint]{myint(1)}, func(v Vertex[myint], d int) error {
		lock.Lock()
		defer lock.Unlock()
		visits = append(visits, v)
		g.Remove(v)
		return nil
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []Vertex[myint]{myint(1), myint(2), myint(3)}
	if !reflect.DeepEqual(visits, expected) {
		t.Fatalf("expected: %#v, got: %#v", expected, visits)
	}
}

const testGraphTransReductionStr = `
1
  2
2
  3
3
`

const testGraphTransReductionMoreStr = `
1
  2
2
  3
3
  4
4
`
