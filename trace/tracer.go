// Interfaces allows us to define an API without
// restriction on the implementation details.

package trace

import (
	"fmt"
	"io"
)

// Tracer is the interface that describes an object capable of
// tracing events through code.
type Tracer interface {
	// using capital T in Tracer means we intend it to be publicly visible
	//...interface means the method is accepting zero or more arguments
	// of any type
	Trace(...interface{})
}

// New Function
// using the -cover flag provides coverage and display how much
// of the code was touched during test execution.
func New(w io.Writer) Tracer {
	// return nil
	return &tracer{out: w}
}

// tracer writes to an io.Writer.
type tracer struct {
	// Trace output is written to out
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// defined a struct call nilTracer
type nilTracer struct{}

// defined a function that inherits from nilTracer struct
// which does nothing
func (t *nilTracer) Trace(a ...interface{}) {}

// Off creates a Tracer that will ignore calls to Trace.
func Off() Tracer {
	return &nilTracer{}
}
