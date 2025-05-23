package parser

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"strconv"
)

// QubitRef 表示对量子比特的引用
type QubitRef struct {
	Register string
	Index    int
}

// QregStatement 表示量子寄存器声明语句
type QregStatement struct {
	Name string
	Size int
}

// Execute 执行量子寄存器声明
func (s *QregStatement) Execute(ctx *ExecutionContext) error {
	if _, exists := ctx.Registers[s.Name]; exists {
		return fmt.Errorf("quantum register %s already defined", s.Name)
	}

	qubits := make([]*Qubit, s.Size)
	for i := 0; i < s.Size; i++ {
		// 初始化为|0⟩态
		qubits[i] = &Qubit{
			Alpha: complex(1.0, 0.0),
			Beta:  complex(0.0, 0.0),
		}
	}

	ctx.Registers[s.Name] = &QuantumRegister{
		Name:  s.Name,
		Size:  s.Size,
		Qubits: qubits,
	}

	return nil
}

// CregStatement 表示经典寄存器声明语句
type CregStatement struct {
	Name string
	Size int
}

// Execute 执行经典寄存器声明
func (s *CregStatement) Execute(ctx *ExecutionContext) error {
	if _, exists := ctx.ClassicalRegisters[s.Name]; exists {
		return fmt.Errorf("classical register %s already defined", s.Name)
	}

	ctx.ClassicalRegisters[s.Name] = &ClassicalRegister{
		Name: s.Name,
		Size: s.Size,
		Bits: make([]bool, s.Size),
	}

	return nil
}

// GateStatement 表示门操作语句
type GateStatement struct {
	Name      string
	Controls  []QubitRef
	Targets   []QubitRef
	Parameter float64
	Params    []string    // 用于参数化门操作，如U门的三个参数
}

// Execute 执行门操作
func (s *GateStatement) Execute(ctx *ExecutionContext) error {
	switch s.Name {
	case "h":
		return s.executeHadamard(ctx)
	case "x":
		return s.executePauliX(ctx)
	case "y":
		return s.executePauliY(ctx)
	case "z":
		return s.executePauliZ(ctx)
	case "id":
		return s.executeIdentity(ctx)
	case "s":
		return s.executePhase(ctx, math.Pi/2)
	case "sdg":
		return s.executePhase(ctx, -math.Pi/2)
	case "t":
		return s.executePhase(ctx, math.Pi/4)
	case "tdg":
		return s.executePhase(ctx, -math.Pi/4)
	case "cx":
		return s.executeCNOT(ctx)
	case "cz":
		return s.executeCZ(ctx)
	case "cy":
		return s.executeCY(ctx)
	case "ch":
		return s.executeCH(ctx)
	case "ccx":
		return s.executeToffoli(ctx)
	case "rx":
		return s.executeRX(ctx)
	case "ry":
		return s.executeRY(ctx)
	case "rz":
		return s.executeRZ(ctx)
	case "u1":
		return s.executeU1(ctx)
	case "u2":
		return s.executeU2(ctx)
	case "u3":
		return s.executeU3(ctx)
	case "U":
		return s.executeU(ctx)
	// 新添加的量子门
	case "cu1":
		return s.executeCU1(ctx)
	case "cu2":
		return s.executeCU2(ctx)
	case "cu3":
		return s.executeCU3(ctx)
	case "crz":
		return s.executeCRZ(ctx)
	case "swap":
		return s.executeSWAP(ctx)
	case "cswap":
		return s.executeCSWAP(ctx)
	case "fredkin":
		return s.executeCSWAP(ctx) // fredkin是cswap的别名
	default:
		// 检查是否是自定义门
		if def, exists := ctx.GateDefs[s.Name]; exists {
			// 创建自定义门调用
			customGate := &CustomGateStatement{
				Name:       s.Name,
				Parameters: make([]float64, len(s.Params)),
				Qubits:     append(s.Targets, s.Controls...),
				Definition: def,
			}
			
			// 将参数字符串转换为浮点数
			for i, paramStr := range s.Params {
				param, err := strconv.ParseFloat(paramStr, 64)
				if err != nil {
					return fmt.Errorf("invalid parameter value: %s", paramStr)
				}
				customGate.Parameters[i] = param
			}
			
			return customGate.Execute(ctx)
		}
		
		return fmt.Errorf("unsupported gate: %s", s.Name)
	}
}

// executeHadamard 执行Hadamard门
func (s *GateStatement) executeHadamard(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("hadamard gate requires exactly one target qubit")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	qubit := reg.Qubits[target.Index]
	
	// Hadamard变换: |0⟩ -> (|0⟩ + |1⟩)/√2, |1⟩ -> (|0⟩ - |1⟩)/√2
	alpha := qubit.Alpha
	beta := qubit.Beta
	
	qubit.Alpha = (alpha + beta) / complex(math.Sqrt(2), 0)
	qubit.Beta = (alpha - beta) / complex(math.Sqrt(2), 0)

	return nil
}

// executePauliX 执行Pauli-X门（NOT门）
func (s *GateStatement) executePauliX(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("Pauli-X gate requires exactly one target qubit")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	qubit := reg.Qubits[target.Index]
	
	// Pauli-X变换: |0⟩ -> |1⟩, |1⟩ -> |0⟩
	alpha := qubit.Alpha
	beta := qubit.Beta
	
	qubit.Alpha = beta
	qubit.Beta = alpha

	return nil
}

// executePauliY 执行Pauli-Y门
func (s *GateStatement) executePauliY(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("Pauli-Y gate requires exactly one target qubit")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	qubit := reg.Qubits[target.Index]
	
	// Pauli-Y变换: |0⟩ -> i|1⟩, |1⟩ -> -i|0⟩
	alpha := qubit.Alpha
	beta := qubit.Beta
	
	qubit.Alpha = -complex(0, 1) * beta
	qubit.Beta = complex(0, 1) * alpha

	return nil
}

// executePauliZ 执行Pauli-Z门
func (s *GateStatement) executePauliZ(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("Pauli-Z gate requires exactly one target qubit")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	qubit := reg.Qubits[target.Index]
	
	// Pauli-Z变换: |0⟩ -> |0⟩, |1⟩ -> -|1⟩
	qubit.Beta = -qubit.Beta

	return nil
}

// executeCNOT 执行CNOT门
func (s *GateStatement) executeCNOT(ctx *ExecutionContext) error {
	if len(s.Controls) != 1 || len(s.Targets) != 1 {
		return fmt.Errorf("CNOT gate requires exactly one control and one target qubit")
	}

	control := s.Controls[0]
	target := s.Targets[0]

	controlReg, exists := ctx.Registers[control.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", control.Register)
	}

	if control.Index < 0 || control.Index >= controlReg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", control.Index, control.Register)
	}

	targetReg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= targetReg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	// 这里简化实现，实际上CNOT门在多量子比特系统中需要考虑纠缠态
	// 这个简化版本只在控制比特为|1⟩时翻转目标比特
	controlQubit := controlReg.Qubits[control.Index]
	targetQubit := targetReg.Qubits[target.Index]

	// 如果控制比特处于|1⟩态的概率非零，则应用X门到目标比特
	if cmplx.Abs(controlQubit.Beta) > 0 {
		// 交换目标比特的Alpha和Beta
		alpha := targetQubit.Alpha
		beta := targetQubit.Beta
		
		targetQubit.Alpha = beta
		targetQubit.Beta = alpha
	}

	return nil
}

// executeToffoli 执行Toffoli门（CCNOT）
func (s *GateStatement) executeToffoli(ctx *ExecutionContext) error {
	if len(s.Controls) != 2 || len(s.Targets) != 1 {
		return fmt.Errorf("Toffoli gate requires exactly two controls and one target qubit")
	}

	control1 := s.Controls[0]
	control2 := s.Controls[1]
	target := s.Targets[0]

	control1Reg, exists := ctx.Registers[control1.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", control1.Register)
	}

	if control1.Index < 0 || control1.Index >= control1Reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", control1.Index, control1.Register)
	}

	control2Reg, exists := ctx.Registers[control2.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", control2.Register)
	}

	if control2.Index < 0 || control2.Index >= control2Reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", control2.Index, control2.Register)
	}

	targetReg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= targetReg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	// 简化实现，只在两个控制比特都为|1⟩时翻转目标比特
	control1Qubit := control1Reg.Qubits[control1.Index]
	control2Qubit := control2Reg.Qubits[control2.Index]
	targetQubit := targetReg.Qubits[target.Index]

	// 如果两个控制比特都处于|1⟩态的概率非零，则应用X门到目标比特
	if cmplx.Abs(control1Qubit.Beta) > 0 && cmplx.Abs(control2Qubit.Beta) > 0 {
		// 交换目标比特的Alpha和Beta
		alpha := targetQubit.Alpha
		beta := targetQubit.Beta
		
		targetQubit.Alpha = beta
		targetQubit.Beta = alpha
	}

	return nil
}

// executeRX 执行RX门（绕X轴旋转）
func (s *GateStatement) executeRX(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("RX gate requires exactly one target qubit")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	qubit := reg.Qubits[target.Index]
	
	// RX(θ)变换矩阵:
	// [cos(θ/2), -i*sin(θ/2)]
	// [-i*sin(θ/2), cos(θ/2)]
	
	theta := s.Parameter
	cos := math.Cos(theta / 2)
	sin := math.Sin(theta / 2)
	
	alpha := qubit.Alpha
	beta := qubit.Beta
	
	qubit.Alpha = complex(cos, 0) * alpha + complex(0, -sin) * beta
	qubit.Beta = complex(0, -sin) * alpha + complex(cos, 0) * beta

	return nil
}

// executeRY 执行RY门（绕Y轴旋转）
func (s *GateStatement) executeRY(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("RY gate requires exactly one target qubit")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	qubit := reg.Qubits[target.Index]
	
	// RY(θ)变换矩阵:
	// [cos(θ/2), -sin(θ/2)]
	// [sin(θ/2), cos(θ/2)]
	
	theta := s.Parameter
	cos := math.Cos(theta / 2)
	sin := math.Sin(theta / 2)
	
	alpha := qubit.Alpha
	beta := qubit.Beta
	
	qubit.Alpha = complex(cos, 0) * alpha + complex(-sin, 0) * beta
	qubit.Beta = complex(sin, 0) * alpha + complex(cos, 0) * beta

	return nil
}

// executeRZ 执行RZ门（绕Z轴旋转）
func (s *GateStatement) executeRZ(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("RZ gate requires exactly one target qubit")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	qubit := reg.Qubits[target.Index]
	
	// RZ(θ)变换矩阵:
	// [e^(-iθ/2), 0]
	// [0, e^(iθ/2)]
	
	theta := s.Parameter
	phase1 := complex(math.Cos(-theta/2), math.Sin(-theta/2))
	phase2 := complex(math.Cos(theta/2), math.Sin(theta/2))
	
	qubit.Alpha = phase1 * qubit.Alpha
	qubit.Beta = phase2 * qubit.Beta

	return nil
}

// MeasureStatement 表示测量语句
type MeasureStatement struct {
	QregName  string
	CregName  string
	QbitIndex int // -1表示整个寄存器
	CbitIndex int // -1表示整个寄存器
}

// Execute 执行测量操作
func (s *MeasureStatement) Execute(ctx *ExecutionContext) error {
	qreg, exists := ctx.Registers[s.QregName]
	if !exists {
		return fmt.Errorf("quantum register %s not found", s.QregName)
	}

	creg, exists := ctx.ClassicalRegisters[s.CregName]
	if !exists {
		return fmt.Errorf("classical register %s not found", s.CregName)
	}

	// 测量整个寄存器
	if s.QbitIndex == -1 && s.CbitIndex == -1 {
		if qreg.Size != creg.Size {
			return fmt.Errorf("register size mismatch: quantum register %s has size %d, classical register %s has size %d",
				s.QregName, qreg.Size, s.CregName, creg.Size)
		}

		for i := 0; i < qreg.Size; i++ {
			result, err := measureQubit(qreg.Qubits[i])
			if err != nil {
				return err
			}
			creg.Bits[i] = result
		}
		return nil
	}

	// 测量单个量子比特
	if s.QbitIndex >= 0 && s.QbitIndex < qreg.Size && s.CbitIndex >= 0 && s.CbitIndex < creg.Size {
		result, err := measureQubit(qreg.Qubits[s.QbitIndex])
		if err != nil {
			return err
		}
		creg.Bits[s.CbitIndex] = result
		return nil
	}

	return fmt.Errorf("invalid qubit or classical bit index")
}

// measureQubit 测量单个量子比特
func measureQubit(qubit *Qubit) (bool, error) {
	// 计算测量为|1⟩的概率
	prob1 := cmplx.Abs(qubit.Beta) * cmplx.Abs(qubit.Beta)
	
	// 生成随机数决定测量结果
	r := rand.Float64()
	
	var result bool
	if r < prob1 {
		// 测量结果为|1⟩
		result = true
		qubit.Alpha = complex(0, 0)
		qubit.Beta = complex(1, 0)
	} else {
		// 测量结果为|0⟩
		result = false
		qubit.Alpha = complex(1, 0)
		qubit.Beta = complex(0, 0)
	}
	
	return result, nil
}
