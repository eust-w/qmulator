package parser

import (
	"fmt"
	"math"
	"math/cmplx"
)

// executeCU1 执行受控U1门（受控相位门）
func (s *GateStatement) executeCU1(ctx *ExecutionContext) error {
	if len(s.Controls) != 1 || len(s.Targets) != 1 {
		return fmt.Errorf("CU1 gate requires exactly one control and one target qubit")
	}

	if len(s.Params) != 1 {
		return fmt.Errorf("CU1 gate requires exactly one parameter (lambda)")
	}

	// 解析参数
	lambda, err := parseSymbolicExpression(s.Params[0])
	if err != nil {
		return fmt.Errorf("invalid lambda parameter for CU1 gate: %s (%v)", s.Params[0], err)
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

	// 获取控制比特和目标比特
	controlQubit := controlReg.Qubits[control.Index]
	targetQubit := targetReg.Qubits[target.Index]

	// 如果控制比特处于|1⟩态的概率非零，则应用U1门到目标比特
	if cmplx.Abs(controlQubit.Beta) > 0 {
		// U1(lambda)门: |0⟩ -> |0⟩, |1⟩ -> e^(i*lambda)|1⟩
		phase := complex(math.Cos(lambda), math.Sin(lambda))
		targetQubit.Beta = phase * targetQubit.Beta
	}

	return nil
}

// executeCU2 执行受控U2门
func (s *GateStatement) executeCU2(ctx *ExecutionContext) error {
	if len(s.Controls) != 1 || len(s.Targets) != 1 {
		return fmt.Errorf("CU2 gate requires exactly one control and one target qubit")
	}

	if len(s.Params) != 2 {
		return fmt.Errorf("CU2 gate requires exactly two parameters (phi, lambda)")
	}

	// 解析参数
	phi, err := parseSymbolicExpression(s.Params[0])
	if err != nil {
		return fmt.Errorf("invalid phi parameter for CU2 gate: %s (%v)", s.Params[0], err)
	}

	lambda, err := parseSymbolicExpression(s.Params[1])
	if err != nil {
		return fmt.Errorf("invalid lambda parameter for CU2 gate: %s (%v)", s.Params[1], err)
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

	// 获取控制比特和目标比特
	controlQubit := controlReg.Qubits[control.Index]
	targetQubit := targetReg.Qubits[target.Index]

	// 如果控制比特处于|1⟩态的概率非零，则应用U2门到目标比特
	if cmplx.Abs(controlQubit.Beta) > 0 {
		// U2(phi,lambda)门: 
		// |0⟩ -> (|0⟩ + e^(i*lambda)|1⟩)/sqrt(2)
		// |1⟩ -> (e^(i*phi)|0⟩ - e^(i*(phi+lambda))|1⟩)/sqrt(2)
		alpha := targetQubit.Alpha
		beta := targetQubit.Beta
		
		// 计算相位因子
		eiPhi := complex(math.Cos(phi), math.Sin(phi))
		eiLambda := complex(math.Cos(lambda), math.Sin(lambda))
		eiPhiPlusLambda := complex(math.Cos(phi+lambda), math.Sin(phi+lambda))
		
		// 应用U2门
		targetQubit.Alpha = (alpha + eiLambda*beta) / complex(math.Sqrt(2), 0)
		targetQubit.Beta = (eiPhi*alpha - eiPhiPlusLambda*beta) / complex(math.Sqrt(2), 0)
	}

	return nil
}

// executeCU3 执行受控U3门
func (s *GateStatement) executeCU3(ctx *ExecutionContext) error {
	if len(s.Controls) != 1 || len(s.Targets) != 1 {
		return fmt.Errorf("CU3 gate requires exactly one control and one target qubit")
	}

	if len(s.Params) != 3 {
		return fmt.Errorf("CU3 gate requires exactly three parameters (theta, phi, lambda)")
	}

	// 解析参数
	theta, err := parseSymbolicExpression(s.Params[0])
	if err != nil {
		return fmt.Errorf("invalid theta parameter for CU3 gate: %s (%v)", s.Params[0], err)
	}

	phi, err := parseSymbolicExpression(s.Params[1])
	if err != nil {
		return fmt.Errorf("invalid phi parameter for CU3 gate: %s (%v)", s.Params[1], err)
	}

	lambda, err := parseSymbolicExpression(s.Params[2])
	if err != nil {
		return fmt.Errorf("invalid lambda parameter for CU3 gate: %s (%v)", s.Params[2], err)
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

	// 获取控制比特和目标比特
	controlQubit := controlReg.Qubits[control.Index]
	targetQubit := targetReg.Qubits[target.Index]

	// 如果控制比特处于|1⟩态的概率非零，则应用U3门到目标比特
	if cmplx.Abs(controlQubit.Beta) > 0 {
		// U3(theta,phi,lambda)门: 
		// |0⟩ -> cos(theta/2)|0⟩ + e^(i*lambda)sin(theta/2)|1⟩
		// |1⟩ -> e^(i*phi)sin(theta/2)|0⟩ - e^(i*(phi+lambda))cos(theta/2)|1⟩
		alpha := targetQubit.Alpha
		beta := targetQubit.Beta
		
		// 计算三角函数值
		cosTheta2 := math.Cos(theta / 2)
		sinTheta2 := math.Sin(theta / 2)
		
		// 计算相位因子
		eiPhi := complex(math.Cos(phi), math.Sin(phi))
		eiLambda := complex(math.Cos(lambda), math.Sin(lambda))
		eiPhiPlusLambda := complex(math.Cos(phi+lambda), math.Sin(phi+lambda))
		
		// 应用U3门
		targetQubit.Alpha = complex(cosTheta2, 0)*alpha + complex(-sinTheta2, 0)*eiPhiPlusLambda*beta
		targetQubit.Beta = complex(sinTheta2, 0)*eiLambda*alpha + complex(cosTheta2, 0)*eiPhi*beta
	}

	return nil
}

// executeCRZ 执行受控RZ门
func (s *GateStatement) executeCRZ(ctx *ExecutionContext) error {
	if len(s.Controls) != 1 || len(s.Targets) != 1 {
		return fmt.Errorf("CRZ gate requires exactly one control and one target qubit")
	}

	if len(s.Params) != 1 {
		return fmt.Errorf("CRZ gate requires exactly one parameter (theta)")
	}

	// 解析参数
	theta, err := parseSymbolicExpression(s.Params[0])
	if err != nil {
		return fmt.Errorf("invalid theta parameter for CRZ gate: %s (%v)", s.Params[0], err)
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

	// 获取控制比特和目标比特
	controlQubit := controlReg.Qubits[control.Index]
	targetQubit := targetReg.Qubits[target.Index]

	// 如果控制比特处于|1⟩态的概率非零，则应用RZ门到目标比特
	if cmplx.Abs(controlQubit.Beta) > 0 {
		// RZ(theta)门: 
		// |0⟩ -> e^(-i*theta/2)|0⟩
		// |1⟩ -> e^(i*theta/2)|1⟩
		phase1 := complex(math.Cos(-theta/2), math.Sin(-theta/2))
		phase2 := complex(math.Cos(theta/2), math.Sin(theta/2))
		
		targetQubit.Alpha = phase1 * targetQubit.Alpha
		targetQubit.Beta = phase2 * targetQubit.Beta
	}

	return nil
}

// executeSWAP 执行SWAP门
func (s *GateStatement) executeSWAP(ctx *ExecutionContext) error {
	if len(s.Targets) != 2 {
		return fmt.Errorf("SWAP gate requires exactly two target qubits")
	}

	target1 := s.Targets[0]
	target2 := s.Targets[1]

	reg1, exists := ctx.Registers[target1.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target1.Register)
	}

	if target1.Index < 0 || target1.Index >= reg1.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target1.Index, target1.Register)
	}

	reg2, exists := ctx.Registers[target2.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target2.Register)
	}

	if target2.Index < 0 || target2.Index >= reg2.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target2.Index, target2.Register)
	}

	// 获取两个目标比特
	qubit1 := reg1.Qubits[target1.Index]
	qubit2 := reg2.Qubits[target2.Index]

	// 交换两个量子比特的状态
	alpha1 := qubit1.Alpha
	beta1 := qubit1.Beta
	
	qubit1.Alpha = qubit2.Alpha
	qubit1.Beta = qubit2.Beta
	
	qubit2.Alpha = alpha1
	qubit2.Beta = beta1

	return nil
}

// executeCSWAP 执行受控SWAP门（Fredkin门）
func (s *GateStatement) executeCSWAP(ctx *ExecutionContext) error {
	if len(s.Controls) != 1 || len(s.Targets) != 2 {
		return fmt.Errorf("CSWAP gate requires exactly one control and two target qubits")
	}

	control := s.Controls[0]
	target1 := s.Targets[0]
	target2 := s.Targets[1]

	controlReg, exists := ctx.Registers[control.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", control.Register)
	}

	if control.Index < 0 || control.Index >= controlReg.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", control.Index, control.Register)
	}

	reg1, exists := ctx.Registers[target1.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target1.Register)
	}

	if target1.Index < 0 || target1.Index >= reg1.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target1.Index, target1.Register)
	}

	reg2, exists := ctx.Registers[target2.Register]
	if !exists {
		return fmt.Errorf("quantum register %s not found", target2.Register)
	}

	if target2.Index < 0 || target2.Index >= reg2.Size {
		return fmt.Errorf("qubit index %d out of range for register %s", target2.Index, target2.Register)
	}

	// 获取控制比特和两个目标比特
	controlQubit := controlReg.Qubits[control.Index]
	qubit1 := reg1.Qubits[target1.Index]
	qubit2 := reg2.Qubits[target2.Index]

	// 如果控制比特处于|1⟩态的概率非零，则交换两个目标比特的状态
	if cmplx.Abs(controlQubit.Beta) > 0 {
		// 交换两个量子比特的状态
		alpha1 := qubit1.Alpha
		beta1 := qubit1.Beta
		
		qubit1.Alpha = qubit2.Alpha
		qubit1.Beta = qubit2.Beta
		
		qubit2.Alpha = alpha1
		qubit2.Beta = beta1
	}

	return nil
}
