package parser

// QASMVersion 表示QASM版本
type QASMVersion string

const (
	// QASM1 表示QASM 1.0版本
	QASM1 QASMVersion = "1.0"
	// QASM2 表示QASM 2.0版本
	QASM2 QASMVersion = "2.0"
)

// Program 表示一个QASM程序
type Program struct {
	Version QASMVersion
	Includes []string
	Statements []Statement
}

// Statement 表示一个QASM语句
type Statement interface {
	Execute(ctx *ExecutionContext) error
}

// ExecutionContext 表示执行上下文
type ExecutionContext struct {
	Registers          map[string]*QuantumRegister
	ClassicalRegisters map[string]*ClassicalRegister
	Gates              map[string]Gate
	GateDefs           map[string]*GateDefinition
}

// NewExecutionContext 创建一个新的执行上下文
func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		Registers:          make(map[string]*QuantumRegister),
		ClassicalRegisters: make(map[string]*ClassicalRegister),
		Gates:              make(map[string]Gate),
		GateDefs:           make(map[string]*GateDefinition),
	}
}

// QuantumRegister 表示量子寄存器
type QuantumRegister struct {
	Name  string
	Size  int
	Qubits []*Qubit
}

// ClassicalRegister 表示经典寄存器
type ClassicalRegister struct {
	Name  string
	Size  int
	Bits  []bool
}

// Qubit 表示量子比特
type Qubit struct {
	// 量子态使用复数表示
	Alpha complex128 // |0⟩的振幅
	Beta  complex128 // |1⟩的振幅
}

// Gate 表示量子门
type Gate interface {
	Apply(qubits []*Qubit) error
}

// ParseResult 表示解析结果
type ParseResult struct {
	Program *Program
	Errors  []error
}
