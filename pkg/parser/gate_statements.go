package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// parseUGateStatement 解析U门语句
func parseUGateStatement(line string) (*GateStatement, error) {
	// U(θ,φ,λ) q;
	re := regexp.MustCompile(`U\(([^,]+),([^,]+),([^)]+)\)\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 6 {
		return nil, fmt.Errorf("invalid U gate statement: %s", line)
	}

	theta := matches[1]
	phi := matches[2]
	lambda := matches[3]
	regName := matches[4]
	qubitIdx, err := strconv.Atoi(matches[5])
	if err != nil {
		return nil, fmt.Errorf("invalid qubit index: %s", matches[5])
	}

	return &GateStatement{
		Name:      "U",
		Targets:   []QubitRef{{Register: regName, Index: qubitIdx}},
		Controls:  []QubitRef{},
		Parameter: 0.0, // 参数将在执行时计算
		Params:    []string{theta, phi, lambda},
	}, nil
}

// parseCXGateStatement 解析CX门语句
func parseCXGateStatement(line string) (*GateStatement, error) {
	// CX q[0], q[1];
	re := regexp.MustCompile(`CX\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid CX gate statement: %s", line)
	}

	controlReg := matches[1]
	controlIdx, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid control qubit index: %s", matches[2])
	}

	targetReg := matches[3]
	targetIdx, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("invalid target qubit index: %s", matches[4])
	}

	return &GateStatement{
		Name:      "cx",
		Controls:  []QubitRef{{Register: controlReg, Index: controlIdx}},
		Targets:   []QubitRef{{Register: targetReg, Index: targetIdx}},
		Parameter: 0.0,
	}, nil
}

// parseCustomGateStatement 解析自定义门调用语句
func parseCustomGateStatement(line string) (Statement, error) {
	// 匹配带参数的门调用: name(param1,param2,...) q[0],q[1],...;
	reWithParams := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9]*)\s*\(([^)]+)\)\s+(.+)`)
	matches := reWithParams.FindStringSubmatch(line)
	if len(matches) == 4 {
		gateName := matches[1]
		paramsStr := matches[2]
		qubitsStr := matches[3]

		// 解析参数
		var paramStrs []string
		if paramsStr != "" {
			for _, p := range strings.Split(paramsStr, ",") {
				paramStrs = append(paramStrs, strings.TrimSpace(p))
			}
		}

		// 尝试将参数转换为浮点数
		params := make([]float64, len(paramStrs))
		for i, p := range paramStrs {
			val, err := strconv.ParseFloat(p, 64)
			if err != nil {
				// 如果不是数字，可能是表达式，暂时设为0
				params[i] = 0.0
			} else {
				params[i] = val
			}
		}

		// 解析量子比特引用
		var qubits []QubitRef
		for _, q := range strings.Split(qubitsStr, ",") {
			q = strings.TrimSpace(q)
			// 匹配q[i]格式
			reQubit := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
			qubitMatches := reQubit.FindStringSubmatch(q)
			if len(qubitMatches) != 3 {
				return nil, fmt.Errorf("invalid qubit reference: %s", q)
			}

			regName := qubitMatches[1]
			idx, err := strconv.Atoi(qubitMatches[2])
			if err != nil {
				return nil, fmt.Errorf("invalid qubit index: %s", qubitMatches[2])
			}

			qubits = append(qubits, QubitRef{Register: regName, Index: idx})
		}

		return &CustomGateStatement{
			Name:       gateName,
			Parameters: params,
			Qubits:     qubits,
		}, nil
	}

	// 匹配不带参数的门调用: name q[0],q[1],...;
	reNoParams := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9]*)\s+(.+)`)
	matches = reNoParams.FindStringSubmatch(line)
	if len(matches) == 3 {
		gateName := matches[1]
		qubitsStr := matches[2]

		// 解析量子比特引用
		var qubits []QubitRef
		for _, q := range strings.Split(qubitsStr, ",") {
			q = strings.TrimSpace(q)
			// 匹配q[i]格式
			reQubit := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
			qubitMatches := reQubit.FindStringSubmatch(q)
			if len(qubitMatches) != 3 {
				return nil, fmt.Errorf("invalid qubit reference: %s", q)
			}

			regName := qubitMatches[1]
			idx, err := strconv.Atoi(qubitMatches[2])
			if err != nil {
				return nil, fmt.Errorf("invalid qubit index: %s", qubitMatches[2])
			}

			qubits = append(qubits, QubitRef{Register: regName, Index: idx})
		}

		return &CustomGateStatement{
			Name:       gateName,
			Parameters: []float64{},
			Qubits:     qubits,
		}, nil
	}

	return nil, fmt.Errorf("unsupported gate operation: %s", line)
}
