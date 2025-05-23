package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

// parseCU1GateStatement 解析CU1门语句
func parseCU1GateStatement(line string) (*GateStatement, error) {
	// cu1(lambda) control, target;
	// 例如: cu1(pi/2) q[0], q[1];
	re := regexp.MustCompile(`cu1\s*\(([^)]+)\)\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 6 {
		return nil, fmt.Errorf("invalid CU1 gate statement: %s", line)
	}

	lambda := matches[1]
	controlReg := matches[2]
	controlIdx, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid control qubit index: %s", matches[3])
	}

	targetReg := matches[4]
	targetIdx, err := strconv.Atoi(matches[5])
	if err != nil {
		return nil, fmt.Errorf("invalid target qubit index: %s", matches[5])
	}

	return &GateStatement{
		Name:      "cu1",
		Controls:  []QubitRef{{Register: controlReg, Index: controlIdx}},
		Targets:   []QubitRef{{Register: targetReg, Index: targetIdx}},
		Parameter: 0.0, // 参数将在执行时解析
		Params:    []string{lambda},
	}, nil
}

// parseCU2GateStatement 解析CU2门语句
func parseCU2GateStatement(line string) (*GateStatement, error) {
	// cu2(phi,lambda) control, target;
	// 例如: cu2(pi/2,pi/4) q[0], q[1];
	re := regexp.MustCompile(`cu2\s*\(([^,]+),([^)]+)\)\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 7 {
		return nil, fmt.Errorf("invalid CU2 gate statement: %s", line)
	}

	phi := matches[1]
	lambda := matches[2]
	controlReg := matches[3]
	controlIdx, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("invalid control qubit index: %s", matches[4])
	}

	targetReg := matches[5]
	targetIdx, err := strconv.Atoi(matches[6])
	if err != nil {
		return nil, fmt.Errorf("invalid target qubit index: %s", matches[6])
	}

	return &GateStatement{
		Name:      "cu2",
		Controls:  []QubitRef{{Register: controlReg, Index: controlIdx}},
		Targets:   []QubitRef{{Register: targetReg, Index: targetIdx}},
		Parameter: 0.0, // 参数将在执行时解析
		Params:    []string{phi, lambda},
	}, nil
}

// parseCU3GateStatement 解析CU3门语句
func parseCU3GateStatement(line string) (*GateStatement, error) {
	// cu3(theta,phi,lambda) control, target;
	// 例如: cu3(pi/2,pi/4,pi/8) q[0], q[1];
	re := regexp.MustCompile(`cu3\s*\(([^,]+),([^,]+),([^)]+)\)\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 8 {
		return nil, fmt.Errorf("invalid CU3 gate statement: %s", line)
	}

	theta := matches[1]
	phi := matches[2]
	lambda := matches[3]
	controlReg := matches[4]
	controlIdx, err := strconv.Atoi(matches[5])
	if err != nil {
		return nil, fmt.Errorf("invalid control qubit index: %s", matches[5])
	}

	targetReg := matches[6]
	targetIdx, err := strconv.Atoi(matches[7])
	if err != nil {
		return nil, fmt.Errorf("invalid target qubit index: %s", matches[7])
	}

	return &GateStatement{
		Name:      "cu3",
		Controls:  []QubitRef{{Register: controlReg, Index: controlIdx}},
		Targets:   []QubitRef{{Register: targetReg, Index: targetIdx}},
		Parameter: 0.0, // 参数将在执行时解析
		Params:    []string{theta, phi, lambda},
	}, nil
}

// parseCRZGateStatement 解析CRZ门语句
func parseCRZGateStatement(line string) (*GateStatement, error) {
	// crz(lambda) control, target;
	// 例如: crz(pi/2) q[0], q[1];
	re := regexp.MustCompile(`crz\s*\(([^)]+)\)\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 6 {
		return nil, fmt.Errorf("invalid CRZ gate statement: %s", line)
	}

	lambda := matches[1]
	controlReg := matches[2]
	controlIdx, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid control qubit index: %s", matches[3])
	}

	targetReg := matches[4]
	targetIdx, err := strconv.Atoi(matches[5])
	if err != nil {
		return nil, fmt.Errorf("invalid target qubit index: %s", matches[5])
	}

	return &GateStatement{
		Name:      "crz",
		Controls:  []QubitRef{{Register: controlReg, Index: controlIdx}},
		Targets:   []QubitRef{{Register: targetReg, Index: targetIdx}},
		Parameter: 0.0, // 参数将在执行时解析
		Params:    []string{lambda},
	}, nil
}

// parseSWAPGateStatement 解析SWAP门语句
func parseSWAPGateStatement(line string) (*GateStatement, error) {
	// swap q[0], q[1];
	re := regexp.MustCompile(`swap\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid SWAP gate statement: %s", line)
	}

	target1Reg := matches[1]
	target1Idx, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid first target qubit index: %s", matches[2])
	}

	target2Reg := matches[3]
	target2Idx, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("invalid second target qubit index: %s", matches[4])
	}

	return &GateStatement{
		Name:      "swap",
		Controls:  []QubitRef{},
		Targets:   []QubitRef{{Register: target1Reg, Index: target1Idx}, {Register: target2Reg, Index: target2Idx}},
		Parameter: 0.0,
	}, nil
}

// parseCSWAPGateStatement 解析CSWAP门语句
func parseCSWAPGateStatement(line string) (*GateStatement, error) {
	// cswap control, target1, target2;
	// 例如: cswap q[0], q[1], q[2];
	re := regexp.MustCompile(`cswap\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 7 {
		return nil, fmt.Errorf("invalid CSWAP gate statement: %s", line)
	}

	controlReg := matches[1]
	controlIdx, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid control qubit index: %s", matches[2])
	}

	target1Reg := matches[3]
	target1Idx, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("invalid first target qubit index: %s", matches[4])
	}

	target2Reg := matches[5]
	target2Idx, err := strconv.Atoi(matches[6])
	if err != nil {
		return nil, fmt.Errorf("invalid second target qubit index: %s", matches[6])
	}

	return &GateStatement{
		Name:      "cswap",
		Controls:  []QubitRef{{Register: controlReg, Index: controlIdx}},
		Targets:   []QubitRef{{Register: target1Reg, Index: target1Idx}, {Register: target2Reg, Index: target2Idx}},
		Parameter: 0.0,
	}, nil
}

// parseFredkinGateStatement 解析Fredkin门语句（CSWAP的别名）
func parseFredkinGateStatement(line string) (*GateStatement, error) {
	// fredkin control, target1, target2;
	// 例如: fredkin q[0], q[1], q[2];
	re := regexp.MustCompile(`fredkin\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 7 {
		return nil, fmt.Errorf("invalid Fredkin gate statement: %s", line)
	}

	controlReg := matches[1]
	controlIdx, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid control qubit index: %s", matches[2])
	}

	target1Reg := matches[3]
	target1Idx, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("invalid first target qubit index: %s", matches[4])
	}

	target2Reg := matches[5]
	target2Idx, err := strconv.Atoi(matches[6])
	if err != nil {
		return nil, fmt.Errorf("invalid second target qubit index: %s", matches[6])
	}

	return &GateStatement{
		Name:      "fredkin",
		Controls:  []QubitRef{{Register: controlReg, Index: controlIdx}},
		Targets:   []QubitRef{{Register: target1Reg, Index: target1Idx}, {Register: target2Reg, Index: target2Idx}},
		Parameter: 0.0,
	}, nil
}
