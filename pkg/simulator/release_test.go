package simulator

import (
	"github.com/eust-w/qmulator/pkg/parser"
	"testing"
)

func TestBundledBellExampleParsesAndRuns(t *testing.T) {
	s := NewSimulator(parser.QASM2)
	r, e := s.RunFile("../../examples/bell_state.qasm")
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Errors) > 0 {
		t.Fatalf("parser errors: %v", r.Errors)
	}
	if len(r.ClassicalRegisters["c"]) != 2 {
		t.Fatal("missing register result")
	}
	// This checks execution plumbing, not entanglement correctness.
}
