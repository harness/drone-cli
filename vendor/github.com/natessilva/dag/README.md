# dag

[![GoDoc](https://godoc.org/github.com/natessilva/dag?status.svg)](https://godoc.org/github.com/natessilva/dag)

dag.Runner is a mechanism to orchestrate goroutines and the order in which they run, using the semantics of a directed acyclic graph.

Create a zero value dag.Runner, add vertices (functions) to it, add edges or dependencies between vertices, and finally invoke Run. Run will run all of the vertices in parallel topological order returning after either all vertices complete, or an error gets returned.

## Example

```go
var r dag.Runner

r.AddVertex("one", func() error {
    fmt.Println("one and two will run in parallel before three")
    return nil
})
r.AddVertex("two", func() error {
    fmt.Println("one and two will run in parallel before three")
    return nil
})
r.AddVertex("three", func() error {
    fmt.Println("three will run before four")
    return errors.New("three is broken")
})
r.AddVertex("four", func() error {
    fmt.Println("four will never run")
    return nil
})

r.AddEdge("one", "three")
r.AddEdge("two", "three")

r.AddEdge("three", "four")

fmt.Printf("the runner terminated with: %v\n", r.Run())
```