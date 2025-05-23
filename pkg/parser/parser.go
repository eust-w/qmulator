package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Parser 表示QASM解析器
type Parser struct {
	version QASMVersion
}

// GetVersion 获取解析器使用的QASM版本
func (p *Parser) GetVersion() QASMVersion {
	return p.version
}

// NewParser 创建一个新的QASM解析器
func NewParser(version QASMVersion) *Parser {
	return &Parser{
		version: version,
	}
}

// Parse 解析QASM代码
func (p *Parser) Parse(reader io.Reader) (*ParseResult, error) {
	scanner := bufio.NewScanner(reader)
	program := &Program{
		Version:    p.version,
		Includes:   []string{},
		Statements: []Statement{},
	}

	lineNum := 0
	var errors []error
	
	// 收集所有行以支持多行门定义
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	// 处理所有行
	for i := 0; i < len(lines); i++ {
		lineNum = i + 1
		line := lines[i]

		// 移除注释
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = line[:idx]
		}

		// 跳过空行
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检查QASM版本声明
		if strings.HasPrefix(line, "OPENQASM") {
			version, err := parseVersion(line)
			if err != nil {
				errors = append(errors, fmt.Errorf("line %d: %v", lineNum, err))
				continue
			}
			program.Version = version
			continue
		}

		// 解析include语句
		if strings.HasPrefix(line, "include") {
			include, err := parseInclude(line)
			if err != nil {
				errors = append(errors, fmt.Errorf("line %d: %v", lineNum, err))
				continue
			}
			program.Includes = append(program.Includes, include)
			continue
		}
		
		// 检查是否是门定义
		if strings.HasPrefix(line, "gate") && strings.Contains(line, "{") {
			// 收集门定义的所有行
			var gateDefLines []string
			gateDefLines = append(gateDefLines, line)
			
			// 查找门定义的结束
			braceCount := 1 // 已经遇到一个左大括号
			for j := i + 1; j < len(lines) && braceCount > 0; j++ {
				nextLine := lines[j]
				
				// 移除注释
				if idx := strings.Index(nextLine, "//"); idx >= 0 {
					nextLine = nextLine[:idx]
				}
				
				nextLine = strings.TrimSpace(nextLine)
				if nextLine == "" {
					i++
					continue
				}
				
				// 计算大括号
				for _, c := range nextLine {
					if c == '{' {
						braceCount++
					} else if c == '}' {
						braceCount--
					}
				}
				
				gateDefLines = append(gateDefLines, nextLine)
				i++
				
				if braceCount == 0 {
					break
				}
			}
			
			// 解析门定义
			gateDefStmt, _, err := ParseGateDefinition(gateDefLines)
			if err != nil {
				errors = append(errors, fmt.Errorf("line %d: %v", lineNum, err))
				continue
			}
			
			program.Statements = append(program.Statements, gateDefStmt)
			continue
		}

		// 解析普通语句
		stmt, err := parseStatement(line)
		if err != nil {
			errors = append(errors, fmt.Errorf("line %d: %v", lineNum, err))
			continue
		}

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &ParseResult{
		Program: program,
		Errors:  errors,
	}, nil
}

// parseVersion 解析QASM版本声明
func parseVersion(line string) (QASMVersion, error) {
	re := regexp.MustCompile(`OPENQASM\s+(\d+\.\d+)\s*;`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 2 {
		return "", fmt.Errorf("invalid QASM version declaration: %s", line)
	}

	version := QASMVersion(matches[1])
	if version != QASM1 && version != QASM2 {
		return "", fmt.Errorf("unsupported QASM version: %s", version)
	}

	return version, nil
}

// parseInclude 解析include语句
func parseInclude(line string) (string, error) {
	re := regexp.MustCompile(`include\s+"([^"]+)"\s*;`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 2 {
		return "", fmt.Errorf("invalid include statement: %s", line)
	}

	return matches[1], nil
}

// parseStatement 解析QASM语句
func parseStatement(line string) (Statement, error) {
	// 移除末尾的分号
	if !strings.HasSuffix(line, ";") {
		return nil, fmt.Errorf("statement must end with semicolon: %s", line)
	}
	line = line[:len(line)-1]
	line = strings.TrimSpace(line)

	// 解析量子寄存器声明
	if strings.HasPrefix(line, "qreg") {
		return parseQregStatement(line)
	}

	// 解析经典寄存器声明
	if strings.HasPrefix(line, "creg") {
		return parseCregStatement(line)
	}

	// 解析测量语句
	if strings.HasPrefix(line, "measure") {
		return parseMeasureStatement(line)
	}
	
	// 解析U门
	if strings.HasPrefix(line, "U(") {
		return parseUGateStatement(line)
	}
	
	// 解析CX门
	if strings.HasPrefix(line, "CX ") {
		return parseCXGateStatement(line)
	}
	
	// 解析CU1门
	if strings.HasPrefix(line, "cu1(") || strings.HasPrefix(line, "cu1 (") {
		return parseCU1GateStatement(line)
	}
	
	// 解析CU2门
	if strings.HasPrefix(line, "cu2(") || strings.HasPrefix(line, "cu2 (") {
		return parseCU2GateStatement(line)
	}
	
	// 解析CU3门
	if strings.HasPrefix(line, "cu3(") || strings.HasPrefix(line, "cu3 (") {
		return parseCU3GateStatement(line)
	}
	
	// 解析CRZ门
	if strings.HasPrefix(line, "crz(") || strings.HasPrefix(line, "crz (") {
		return parseCRZGateStatement(line)
	}
	
	// 解析SWAP门
	if strings.HasPrefix(line, "swap ") {
		return parseSWAPGateStatement(line)
	}
	
	// 解析CSWAP门
	if strings.HasPrefix(line, "cswap ") {
		return parseCSWAPGateStatement(line)
	}
	
	// 解析Fredkin门
	if strings.HasPrefix(line, "fredkin ") {
		return parseFredkinGateStatement(line)
	}

	// 解析其他门操作
	return parseGateStatement(line)
}

// parseQregStatement 解析量子寄存器声明
func parseQregStatement(line string) (*QregStatement, error) {
	re := regexp.MustCompile(`qreg\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid quantum register declaration: %s", line)
	}

	name := matches[1]
	size, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid register size: %s", matches[2])
	}

	return &QregStatement{
		Name: name,
		Size: size,
	}, nil
}

// parseCregStatement 解析经典寄存器声明
func parseCregStatement(line string) (*CregStatement, error) {
	re := regexp.MustCompile(`creg\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid classical register declaration: %s", line)
	}

	name := matches[1]
	size, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid register size: %s", matches[2])
	}

	return &CregStatement{
		Name: name,
		Size: size,
	}, nil
}

// parseMeasureStatement 解析测量语句
func parseMeasureStatement(line string) (*MeasureStatement, error) {
	// 支持两种格式：
	// 1. measure q -> c;
	// 2. measure q[0] -> c[0];
	
	// 检查第一种格式
	re1 := regexp.MustCompile(`measure\s+([a-zA-Z][a-zA-Z0-9]*)\s+->\s+([a-zA-Z][a-zA-Z0-9]*)`)
	matches := re1.FindStringSubmatch(line)
	if len(matches) == 3 {
		return &MeasureStatement{
			QregName:  matches[1],
			CregName:  matches[2],
			QbitIndex: -1, // -1表示整个寄存器
			CbitIndex: -1,
		}, nil
	}

	// 检查第二种格式
	re2 := regexp.MustCompile(`measure\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]\s+->\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	matches = re2.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid measure statement: %s", line)
	}

	qregName := matches[1]
	qbitIdx, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid qubit index: %s", matches[2])
	}

	cregName := matches[3]
	cbitIdx, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("invalid classical bit index: %s", matches[4])
	}

	return &MeasureStatement{
		QregName:  qregName,
		CregName:  cregName,
		QbitIndex: qbitIdx,
		CbitIndex: cbitIdx,
	}, nil
}

// parseGateStatement 解析门操作语句
func parseGateStatement(line string) (Statement, error) {
	// 解析单量子比特门
	singleQubitGates := []string{"h", "x", "y", "z", "s", "t", "sdg", "tdg", "id"}
	for _, gate := range singleQubitGates {
		if strings.HasPrefix(line, gate+" ") {
			return parseSingleQubitGate(line, gate)
		}
	}

	// 解析双量子比特门
	twoQubitGates := []string{"cx", "cz", "cy", "ch", "crz", "cu1", "cu3"}
	for _, gate := range twoQubitGates {
		if strings.HasPrefix(line, gate+" ") {
			return parseTwoQubitGate(line, gate)
		}
	}

	// 解析三量子比特门
	if strings.HasPrefix(line, "ccx ") {
		return parseThreeQubitGate(line, "ccx")
	}

	// 解析旋转门
	rotationGates := []string{"rx", "ry", "rz", "u1", "u2", "u3"}
	for _, gate := range rotationGates {
		if strings.HasPrefix(line, gate+"(") {
			return parseRotationGate(line, gate)
		}
	}
	
	// 如果不是内置门，尝试解析为自定义门调用
	return parseCustomGateStatement(line)
}

// parseSingleQubitGate 解析单量子比特门
func parseSingleQubitGate(line, gateName string) (*GateStatement, error) {
	re := regexp.MustCompile(gateName + `\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid %s gate statement: %s", gateName, line)
	}

	regName := matches[1]
	qubitIdx, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid qubit index: %s", matches[2])
	}

	return &GateStatement{
		Name:      gateName,
		Targets:   []QubitRef{{Register: regName, Index: qubitIdx}},
		Controls:  []QubitRef{},
		Parameter: 0.0,
	}, nil
}

// parseTwoQubitGate 解析双量子比特门
func parseTwoQubitGate(line, gateName string) (*GateStatement, error) {
	re := regexp.MustCompile(gateName + `\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\],\s*([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid %s gate statement: %s", gateName, line)
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
		Name:      gateName,
		Controls:  []QubitRef{{Register: controlReg, Index: controlIdx}},
		Targets:   []QubitRef{{Register: targetReg, Index: targetIdx}},
		Parameter: 0.0,
	}, nil
}

// parseThreeQubitGate 解析三量子比特门
func parseThreeQubitGate(line, gateName string) (*GateStatement, error) {
	re := regexp.MustCompile(gateName + `\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\],\s*([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\],\s*([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 7 {
		return nil, fmt.Errorf("invalid %s gate statement: %s", gateName, line)
	}

	control1Reg := matches[1]
	control1Idx, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid first control qubit index: %s", matches[2])
	}

	control2Reg := matches[3]
	control2Idx, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("invalid second control qubit index: %s", matches[4])
	}

	targetReg := matches[5]
	targetIdx, err := strconv.Atoi(matches[6])
	if err != nil {
		return nil, fmt.Errorf("invalid target qubit index: %s", matches[6])
	}

	return &GateStatement{
		Name: gateName,
		Controls: []QubitRef{
			{Register: control1Reg, Index: control1Idx},
			{Register: control2Reg, Index: control2Idx},
		},
		Targets:   []QubitRef{{Register: targetReg, Index: targetIdx}},
		Parameter: 0.0,
	}, nil
}

// parseRotationGate 解析旋转门
func parseRotationGate(line, gateName string) (*GateStatement, error) {
	re := regexp.MustCompile(gateName + `\(([^)]+)\)\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid %s gate statement: %s", gateName, line)
	}

	paramStr := matches[1]
	param, err := strconv.ParseFloat(paramStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid rotation parameter: %s", paramStr)
	}

	regName := matches[2]
	qubitIdx, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid qubit index: %s", matches[3])
	}

	return &GateStatement{
		Name:      gateName,
		Targets:   []QubitRef{{Register: regName, Index: qubitIdx}},
		Controls:  []QubitRef{},
		Parameter: param,
	}, nil
}
