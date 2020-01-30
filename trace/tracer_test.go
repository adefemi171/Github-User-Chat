package trace

import (
	"bytes"
	"testing"
)

// Writing a unit test and executed whe tests is running.
// run go test in terminal
func TestNew(t *testing.T) {
	// t.Error("No test written yet!")
	// Updating the TestNew function
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Return from Function New should not be nil")
	} else {
		tracer.Trace("Trace Package found")
		if buf.String() != "Trace Package found \n" {
			t.Errorf("Trace should not write '%s'.", buf.String())
			// tracer.Trace("Trace Package Found")
		}
	}
}

// calls the off function to get a silent tracer before
// making a call to Trace
func TestOff(t *testing.T) {
	var silentTracer Tracer = off()
	silentTracer.Trace("Nothing")
}
