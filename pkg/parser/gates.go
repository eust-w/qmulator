package parser

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
)

// executeIdentity 执行Identity门（单位门）
func (s *GateStatement) executeIdentity(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("identity gate requires exactly one target qubit")
	}

	// 单位门不改变量子态，所以这里不需要做任何操作
	return nil
}

// executePhase 执行相位门（S, T等）
func (s *GateStatement) executePhase(ctx *ExecutionContext, angle float64) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("phase gate requires exactly one target qubit")
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
	
	// 相位门: |0⟩ -> |0⟩, |1⟩ -> e^(i*angle)|1⟩
	phase := complex(math.Cos(angle), math.Sin(angle))
	qubit.Beta = phase * qubit.Beta

	return nil
}

// executeCZ 执行CZ门
func (s *GateStatement) executeCZ(ctx *ExecutionContext) error {
	if len(s.Controls) != 1 || len(s.Targets) != 1 {
		return fmt.Errorf("CZ gate requires exactly one control and one target qubit")
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

	// 简化实现，CZ门在控制比特为|1⟩时对目标比特应用Z门
	controlQubit := controlReg.Qubits[control.Index]
	targetQubit := targetReg.Qubits[target.Index]

	// 如果控制比特处于|1⟩态的概率非零，则应用Z门到目标比特
	if cmplx.Abs(controlQubit.Beta) > 0 {
		// Z门: |0⟩ -> |0⟩, |1⟩ -> -|1⟩
		targetQubit.Beta = -targetQubit.Beta
	}

	return nil
}

// executeCY 执行CY门
func (s *GateStatement) executeCY(ctx *ExecutionContext) error {
	if len(s.Controls) != 1 || len(s.Targets) != 1 {
		return fmt.Errorf("CY gate requires exactly one control and one target qubit")
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

	// 简化实现，CY门在控制比特为|1⟩时对目标比特应用Y门
	controlQubit := controlReg.Qubits[control.Index]
	targetQubit := targetReg.Qubits[target.Index]

	// 如果控制比特处于|1⟩态的概率非零，则应用Y门到目标比特
	if cmplx.Abs(controlQubit.Beta) > 0 {
		// Y门: |0⟩ -> i|1⟩, |1⟩ -> -i|0⟩
		alpha := targetQubit.Alpha
		beta := targetQubit.Beta
		
		targetQubit.Alpha = -complex(0, 1) * beta
		targetQubit.Beta = complex(0, 1) * alpha
	}

	return nil
}

// executeCH 执行CH门
func (s *GateStatement) executeCH(ctx *ExecutionContext) error {
	if len(s.Controls) != 1 || len(s.Targets) != 1 {
		return fmt.Errorf("CH gate requires exactly one control and one target qubit")
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

	// 简化实现，CH门在控制比特为|1⟩时对目标比特应用H门
	controlQubit := controlReg.Qubits[control.Index]
	targetQubit := targetReg.Qubits[target.Index]

	// 如果控制比特处于|1⟩态的概率非零，则应用H门到目标比特
	if cmplx.Abs(controlQubit.Beta) > 0 {
		// Hadamard门: |0⟩ -> (|0⟩ + |1⟩)/√2, |1⟩ -> (|0⟩ - |1⟩)/√2
		alpha := targetQubit.Alpha
		beta := targetQubit.Beta
		
		targetQubit.Alpha = (alpha + beta) / complex(math.Sqrt(2), 0)
		targetQubit.Beta = (alpha - beta) / complex(math.Sqrt(2), 0)
	}

	return nil
}

// executeU1 执行U1门
func (s *GateStatement) executeU1(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("U1 gate requires exactly one target qubit")
	}

	if len(s.Params) != 1 {
		return fmt.Errorf("U1 gate requires exactly one parameter")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	// 解析参数
	lambda, err := strconv.ParseFloat(s.Params[0], 64)
	if err != nil {
		return fmt.Errorf("invalid parameter for U1 gate: %s", s.Params[0])
	}

	qubit := reg.Qubits[target.Index]
	
	// U1(λ)门: |0⟩ -> |0⟩, |1⟩ -> e^(iλ)|1⟩
	phase := complex(math.Cos(lambda), math.Sin(lambda))
	qubit.Beta = phase * qubit.Beta

	return nil
}

// executeU2 执行U2门
func (s *GateStatement) executeU2(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("U2 gate requires exactly one target qubit")
	}

	if len(s.Params) != 2 {
		return fmt.Errorf("U2 gate requires exactly two parameters")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	// 解析参数
	phi, err := strconv.ParseFloat(s.Params[0], 64)
	if err != nil {
		return fmt.Errorf("invalid phi parameter for U2 gate: %s", s.Params[0])
	}

	lambda, err := strconv.ParseFloat(s.Params[1], 64)
	if err != nil {
		return fmt.Errorf("invalid lambda parameter for U2 gate: %s", s.Params[1])
	}

	qubit := reg.Qubits[target.Index]
	
	// U2(φ,λ)门: |0⟩ -> (|0⟩ + e^(iλ)|1⟩)/√2, |1⟩ -> (e^(iφ)|0⟩ - e^(i(φ+λ))|1⟩)/√2
	alpha := qubit.Alpha
	beta := qubit.Beta
	
	// 计算相位因子
	eiPhi := complex(math.Cos(phi), math.Sin(phi))
	eiLambda := complex(math.Cos(lambda), math.Sin(lambda))
	eiPhiPlusLambda := complex(math.Cos(phi+lambda), math.Sin(phi+lambda))
	
	// 应用U2门
	qubit.Alpha = (alpha + eiLambda*beta) / complex(math.Sqrt(2), 0)
	qubit.Beta = (eiPhi*alpha - eiPhiPlusLambda*beta) / complex(math.Sqrt(2), 0)

	return nil
}

// executeU3 执行U3门
func (s *GateStatement) executeU3(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("U3 gate requires exactly one target qubit")
	}

	if len(s.Params) != 3 {
		return fmt.Errorf("U3 gate requires exactly three parameters")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	// 解析参数
	theta, err := strconv.ParseFloat(s.Params[0], 64)
	if err != nil {
		return fmt.Errorf("invalid theta parameter for U3 gate: %s", s.Params[0])
	}

	phi, err := strconv.ParseFloat(s.Params[1], 64)
	if err != nil {
		return fmt.Errorf("invalid phi parameter for U3 gate: %s", s.Params[1])
	}

	lambda, err := strconv.ParseFloat(s.Params[2], 64)
	if err != nil {
		return fmt.Errorf("invalid lambda parameter for U3 gate: %s", s.Params[2])
	}

	qubit := reg.Qubits[target.Index]
	
	// U3(θ,φ,λ)门: 
	// |0⟩ -> cos(θ/2)|0⟩ + e^(iλ)sin(θ/2)|1⟩
	// |1⟩ -> e^(iφ)sin(θ/2)|0⟩ - e^(i(φ+λ))cos(θ/2)|1⟩
	alpha := qubit.Alpha
	beta := qubit.Beta
	
	// 计算三角函数值
	cosTheta2 := math.Cos(theta / 2)
	sinTheta2 := math.Sin(theta / 2)
	
	// 计算相位因子
	eiPhi := complex(math.Cos(phi), math.Sin(phi))
	eiLambda := complex(math.Cos(lambda), math.Sin(lambda))
	eiPhiPlusLambda := complex(math.Cos(phi+lambda), math.Sin(phi+lambda))
	
	// 应用U3门
	qubit.Alpha = complex(cosTheta2, 0)*alpha + complex(-sinTheta2, 0)*eiPhiPlusLambda*beta
	qubit.Beta = complex(sinTheta2, 0)*eiLambda*alpha + complex(cosTheta2, 0)*eiPhi*beta

	return nil
}

// executeU 执行U门（通用单量子比特门）
func (s *GateStatement) executeU(ctx *ExecutionContext) error {
	if len(s.Targets) != 1 {
		return fmt.Errorf("U gate requires exactly one target qubit")
	}

	if len(s.Params) != 3 {
		return fmt.Errorf("U gate requires exactly three parameters")
	}

	target := s.Targets[0]
	reg, exists := ctx.Registers[target.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target.Register)
	}

	if target.Index < 0 || target.Index >= reg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target.Index, target.Register)
	}

	// 解析参数
	theta, err := strconv.ParseFloat(s.Params[0], 64)
	if err != nil {
		return fmt.Errorf("invalid theta parameter for U gate: %s", s.Params[0])
	}

	phi, err := strconv.ParseFloat(s.Params[1], 64)
	if err != nil {
		return fmt.Errorf("invalid phi parameter for U gate: %s", s.Params[1])
	}

	lambda, err := strconv.ParseFloat(s.Params[2], 64)
	if err != nil {
		return fmt.Errorf("invalid lambda parameter for U gate: %s", s.Params[2])
	}

	qubit := reg.Qubits[target.Index]
	
	// U(θ,φ,λ)门与U3门相同: 
	// |0⟩ -> cos(θ/2)|0⟩ + e^(iλ)sin(θ/2)|1⟩
	// |1⟩ -> e^(iφ)sin(θ/2)|0⟩ - e^(i(φ+λ))cos(θ/2)|1⟩
	alpha := qubit.Alpha
	beta := qubit.Beta
	
	// 计算三角函数值
	cosTheta2 := math.Cos(theta / 2)
	sinTheta2 := math.Sin(theta / 2)
	
	// 计算相位因子
	eiPhi := complex(math.Cos(phi), math.Sin(phi))
	eiLambda := complex(math.Cos(lambda), math.Sin(lambda))
	eiPhiPlusLambda := complex(math.Cos(phi+lambda), math.Sin(phi+lambda))
	
	// 应用U门
	qubit.Alpha = complex(cosTheta2, 0)*alpha + complex(-sinTheta2, 0)*eiPhiPlusLambda*beta
	qubit.Beta = complex(sinTheta2, 0)*eiLambda*alpha + complex(cosTheta2, 0)*eiPhi*beta

	return nil
}
