package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// GateDefinition 表示自定义门定义
type GateDefinition struct {
	Name       string
	Parameters []string
	Qubits     []string
	Body       []Statement
}

// GateDefStatement 表示门定义语句
type GateDefStatement struct {
	Definition *GateDefinition
}

// Execute 执行门定义语句
func (s *GateDefStatement) Execute(ctx *ExecutionContext) error {
	// 将门定义添加到上下文中
	ctx.GateDefs[s.Definition.Name] = s.Definition
	return nil
}

// ParseGateDefinition 解析门定义语句
func ParseGateDefinition(lines []string) (*GateDefStatement, int, error) {
	if len(lines) == 0 {
		return nil, 0, fmt.Errorf("empty gate definition")
	}

	// 解析门定义头部
	headerLine := lines[0]
	def, err := parseGateHeader(headerLine)
	if err != nil {
		return nil, 0, err
	}

	// 查找门定义体
	bodyStart := 1
	bodyEnd := -1
	for i := bodyStart; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "}" {
			bodyEnd = i
			break
		}
	}

	if bodyEnd == -1 {
		return nil, 0, fmt.Errorf("unclosed gate definition")
	}

	// 解析门定义体
	body := lines[bodyStart:bodyEnd]
	statements, err := parseGateBody(body, def)
	if err != nil {
		return nil, 0, err
	}

	def.Body = statements

	return &GateDefStatement{Definition: def}, bodyEnd + 1, nil
}

// parseGateHeader 解析门定义头部
func parseGateHeader(line string) (*GateDefinition, error) {
	// 移除结尾的大括号和分号
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, "{")
	line = strings.TrimSpace(line)

	// 匹配带参数的门定义: gate name(param1,param2,...) qubit1,qubit2,...
	reWithParams := regexp.MustCompile(`gate\s+([a-zA-Z][a-zA-Z0-9]*)\s*\(([^)]*)\)\s+([^{]+)`)
	matches := reWithParams.FindStringSubmatch(line)
	if len(matches) == 4 {
		name := matches[1]
		paramsStr := matches[2]
		qubitsStr := matches[3]

		// 解析参数
		var params []string
		if paramsStr != "" {
			for _, p := range strings.Split(paramsStr, ",") {
				params = append(params, strings.TrimSpace(p))
			}
		}

		// 解析量子比特
		var qubits []string
		for _, q := range strings.Split(qubitsStr, ",") {
			qubits = append(qubits, strings.TrimSpace(q))
		}

		return &GateDefinition{
			Name:       name,
			Parameters: params,
			Qubits:     qubits,
			Body:       []Statement{},
		}, nil
	}

	// 匹配不带参数的门定义: gate name qubit1,qubit2,...
	reNoParams := regexp.MustCompile(`gate\s+([a-zA-Z][a-zA-Z0-9]*)\s+([^{]+)`)
	matches = reNoParams.FindStringSubmatch(line)
	if len(matches) == 3 {
		name := matches[1]
		qubitsStr := matches[2]

		// 解析量子比特
		var qubits []string
		for _, q := range strings.Split(qubitsStr, ",") {
			qubits = append(qubits, strings.TrimSpace(q))
		}

		return &GateDefinition{
			Name:       name,
			Parameters: []string{},
			Qubits:     qubits,
			Body:       []Statement{},
		}, nil
	}

	return nil, fmt.Errorf("invalid gate definition: %s", line)
}

// parseGateBody 解析门定义体
func parseGateBody(lines []string, def *GateDefinition) ([]Statement, error) {
	var statements []Statement

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// 确保语句以分号结尾
		if !strings.HasSuffix(line, ";") {
			return nil, fmt.Errorf("statement in gate body must end with semicolon: %s", line)
		}
		line = line[:len(line)-1] // 移除分号

		// 解析内置量子门（U门和CX门）
		if strings.HasPrefix(line, "U(") {
			stmt, err := parseUGate(line, def)
			if err != nil {
				return nil, err
			}
			statements = append(statements, stmt)
		} else if strings.HasPrefix(line, "CX") {
			stmt, err := parseCXGate(line, def)
			if err != nil {
				return nil, err
			}
			statements = append(statements, stmt)
		} else {
			// 解析自定义门调用
			stmt, err := parseCustomGateCall(line, def)
			if err != nil {
				return nil, err
			}
			statements = append(statements, stmt)
		}
	}

	return statements, nil
}

// parseUGate 解析U门
func parseUGate(line string, def *GateDefinition) (*GateStatement, error) {
	// U(theta,phi,lambda) q;
	re := regexp.MustCompile(`U\(([^,]+),([^,]+),([^)]+)\)\s+([a-zA-Z][a-zA-Z0-9]*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid U gate: %s", line)
	}

	theta := matches[1]
	phi := matches[2]
	lambda := matches[3]
	qubit := matches[4]

	// 检查量子比特是否在门定义中声明
	found := false
	for _, q := range def.Qubits {
		if q == qubit {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("qubit %s not declared in gate definition", qubit)
	}

	// 创建U门语句
	return &GateStatement{
		Name:      "U",
		Targets:   []QubitRef{{Register: qubit, Index: -1}}, // 在门定义中，我们使用Index=-1表示形参
		Controls:  []QubitRef{},
		Parameter: 0.0, // 参数将在实例化时填充
		Params:    []string{theta, phi, lambda},
	}, nil
}

// parseCXGate 解析CX门
func parseCXGate(line string, def *GateDefinition) (*GateStatement, error) {
	// CX q,r;
	re := regexp.MustCompile(`CX\s+([a-zA-Z][a-zA-Z0-9]*)\s*,\s*([a-zA-Z][a-zA-Z0-9]*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid CX gate: %s", line)
	}

	control := matches[1]
	target := matches[2]

	// 检查量子比特是否在门定义中声明
	foundControl := false
	foundTarget := false
	for _, q := range def.Qubits {
		if q == control {
			foundControl = true
		}
		if q == target {
			foundTarget = true
		}
	}
	if !foundControl {
		return nil, fmt.Errorf("control qubit %s not declared in gate definition", control)
	}
	if !foundTarget {
		return nil, fmt.Errorf("target qubit %s not declared in gate definition", target)
	}

	// 创建CX门语句
	return &GateStatement{
		Name:      "cx",
		Controls:  []QubitRef{{Register: control, Index: -1}}, // 在门定义中，我们使用Index=-1表示形参
		Targets:   []QubitRef{{Register: target, Index: -1}},
		Parameter: 0.0,
	}, nil
}

// parseCustomGateCall 解析自定义门调用
func parseCustomGateCall(line string, def *GateDefinition) (*GateStatement, error) {
	// 匹配带参数的门调用: name(param1,param2,...) qubit1,qubit2,...
	reWithParams := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9]*)\s*\(([^)]*)\)\s+(.+)`)
	matches := reWithParams.FindStringSubmatch(line)
	if len(matches) == 4 {
		gateName := matches[1]
		paramsStr := matches[2]
		qubitsStr := matches[3]

		// 解析参数
		var params []string
		if paramsStr != "" {
			for _, p := range strings.Split(paramsStr, ",") {
				params = append(params, strings.TrimSpace(p))
			}
		}

		// 解析量子比特
		var qubits []string
		for _, q := range strings.Split(qubitsStr, ",") {
			qubits = append(qubits, strings.TrimSpace(q))
		}

		// 检查量子比特是否在门定义中声明
		for _, q := range qubits {
			found := false
			for _, defQ := range def.Qubits {
				if q == defQ {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("qubit %s not declared in gate definition", q)
			}
		}

		// 创建自定义门语句
		return &GateStatement{
			Name:      gateName,
			Targets:   []QubitRef{{Register: qubits[0], Index: -1}}, // 简化处理，假设第一个是目标
			Controls:  []QubitRef{},                                 // 控制比特将在实例化时填充
			Parameter: 0.0,
			Params:    params,
		}, nil
	}

	// 匹配不带参数的门调用: name qubit1,qubit2,...
	reNoParams := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9]*)\s+(.+)`)
	matches = reNoParams.FindStringSubmatch(line)
	if len(matches) == 3 {
		gateName := matches[1]
		qubitsStr := matches[2]

		// 解析量子比特
		var qubits []string
		for _, q := range strings.Split(qubitsStr, ",") {
			qubits = append(qubits, strings.TrimSpace(q))
		}

		// 检查量子比特是否在门定义中声明
		for _, q := range qubits {
			found := false
			for _, defQ := range def.Qubits {
				if q == defQ {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("qubit %s not declared in gate definition", q)
			}
		}

		// 创建自定义门语句
		targets := []QubitRef{}
		controls := []QubitRef{}

		// 简化处理，假设第一个是目标，其余是控制
		if len(qubits) > 0 {
			targets = append(targets, QubitRef{Register: qubits[0], Index: -1})
		}
		for i := 1; i < len(qubits); i++ {
			controls = append(controls, QubitRef{Register: qubits[i], Index: -1})
		}

		return &GateStatement{
			Name:      gateName,
			Targets:   targets,
			Controls:  controls,
			Parameter: 0.0,
		}, nil
	}

	return nil, fmt.Errorf("invalid gate call: %s", line)
}

// CustomGateStatement 表示自定义门调用语句
type CustomGateStatement struct {
	Name       string
	Parameters []float64
	Qubits     []QubitRef
	Definition *GateDefinition
}

// Execute 执行自定义门调用
func (s *CustomGateStatement) Execute(ctx *ExecutionContext) error {
	// 获取门定义
	def, exists := ctx.GateDefs[s.Name]
	if !exists {
		return fmt.Errorf("gate %s not defined", s.Name)
	}

	// 检查参数和量子比特数量是否匹配
	if len(s.Parameters) != len(def.Parameters) {
		return fmt.Errorf("gate %s expects %d parameters, got %d", s.Name, len(def.Parameters), len(s.Parameters))
	}
	if len(s.Qubits) != len(def.Qubits) {
		return fmt.Errorf("gate %s expects %d qubits, got %d", s.Name, len(def.Qubits), len(s.Qubits))
	}

	// 创建参数映射
	paramMap := make(map[string]float64)
	for i, name := range def.Parameters {
		paramMap[name] = s.Parameters[i]
	}

	// 创建量子比特映射
	qubitMap := make(map[string]QubitRef)
	for i, name := range def.Qubits {
		qubitMap[name] = s.Qubits[i]
	}

	// 执行门定义体中的每个语句
	for _, stmt := range def.Body {
		// 根据门定义中的语句类型进行处理
		switch s := stmt.(type) {
		case *GateStatement:
			// 创建一个新的门语句，替换形参为实参
			newStmt := &GateStatement{
				Name:      s.Name,
				Controls:  make([]QubitRef, len(s.Controls)),
				Targets:   make([]QubitRef, len(s.Targets)),
				Parameter: s.Parameter,
				Params:    make([]string, len(s.Params)),
			}

			// 替换控制比特
			for i, control := range s.Controls {
				if control.Index == -1 {
					// 这是一个形参，需要替换
					actualQubit, exists := qubitMap[control.Register]
					if !exists {
						return fmt.Errorf("qubit %s not found in gate definition", control.Register)
					}
					newStmt.Controls[i] = actualQubit
				} else {
					// 这是一个实际的量子比特索引
					newStmt.Controls[i] = control
				}
			}

			// 替换目标比特
			for i, target := range s.Targets {
				if target.Index == -1 {
					// 这是一个形参，需要替换
					actualQubit, exists := qubitMap[target.Register]
					if !exists {
						return fmt.Errorf("qubit %s not found in gate definition", target.Register)
					}
					newStmt.Targets[i] = actualQubit
				} else {
					// 这是一个实际的量子比特索引
					newStmt.Targets[i] = target
				}
			}

			// 替换参数
			for i, param := range s.Params {
				// 检查参数是否是一个数值
				if _, err := strconv.ParseFloat(param, 64); err == nil {
					// 这是一个数值，直接使用
					newStmt.Params[i] = param
				} else {
					// 这是一个形参，需要替换
					if paramValue, exists := paramMap[param]; exists {
						newStmt.Params[i] = fmt.Sprintf("%f", paramValue)
					} else {
						return fmt.Errorf("parameter %s not found in gate definition", param)
					}
				}
			}

			// 执行替换后的门语句
			if err := newStmt.Execute(ctx); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported statement type in gate definition")
		}
	}

	return nil
}
