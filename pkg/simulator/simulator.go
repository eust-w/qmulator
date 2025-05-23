package simulator

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/d-robotics/qmulator/pkg/parser"
)

// SimulationResult 表示模拟结果
type SimulationResult struct {
	ClassicalRegisters map[string][]bool
	Errors             []error
}

// Simulator 表示量子模拟器
type Simulator struct {
	parser        *parser.Parser
	includePaths  []string
	includeManager *parser.IncludeManager
}

// NewSimulator 创建一个新的量子模拟器
func NewSimulator(version parser.QASMVersion) *Simulator {
	// 默认包含路径
	includePaths := []string{
		".",                    // 当前目录
		"./examples",           // examples目录
		"../examples",          // 相对于当前目录的examples
		"/usr/local/lib/qasm", // 系统库路径
	}
	
	includeManager := parser.NewIncludeManager(includePaths)
	
	return &Simulator{
		parser:        parser.NewParser(version),
		includePaths:  includePaths,
		includeManager: includeManager,
	}
}

// RunFile 运行QASM文件
func (s *Simulator) RunFile(filename string) (*SimulationResult, error) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	
	// 添加文件目录到包含路径
	fileDir := filepath.Dir(filename)
	if !containsPath(s.includePaths, fileDir) {
		s.includePaths = append(s.includePaths, fileDir)
		s.includeManager = parser.NewIncludeManager(s.includePaths)
	}
	
	// 处理包含文件
	processedCode, err := s.includeManager.ProcessIncludes(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to process includes: %v", err)
	}
	
	return s.RunString(processedCode)
}

// containsPath 检查路径数组是否包含指定路径
func containsPath(paths []string, path string) bool {
	for _, p := range paths {
		if p == path {
			return true
		}
	}
	return false
}

// RunString 运行QASM代码字符串
func (s *Simulator) RunString(code string) (*SimulationResult, error) {
	// 处理包含文件
	processedCode, err := s.includeManager.ProcessIncludes(code)
	if err != nil {
		return nil, fmt.Errorf("failed to process includes: %v", err)
	}
	
	return s.Run(strings.NewReader(processedCode))
}

// Run 运行QASM代码
func (s *Simulator) Run(reader io.Reader) (*SimulationResult, error) {
	// 解析QASM代码
	parseResult, err := s.parser.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	// 如果有解析错误，返回错误
	if len(parseResult.Errors) > 0 {
		return &SimulationResult{
			Errors: parseResult.Errors,
		}, nil
	}

	// 创建执行上下文
	ctx := parser.NewExecutionContext()

	// 执行程序
	for _, stmt := range parseResult.Program.Statements {
		if err := stmt.Execute(ctx); err != nil {
			return nil, fmt.Errorf("execution error: %v", err)
		}
	}

	// 收集结果
	result := &SimulationResult{
		ClassicalRegisters: make(map[string][]bool),
	}

	// 复制经典寄存器的值
	for name, reg := range ctx.ClassicalRegisters {
		result.ClassicalRegisters[name] = make([]bool, reg.Size)
		copy(result.ClassicalRegisters[name], reg.Bits)
	}

	return result, nil
}

// FormatResult 格式化模拟结果为字符串
func (s *Simulator) FormatResult(result *SimulationResult) string {
	if result == nil {
		return "No result available"
	}

	if len(result.Errors) > 0 {
		var errMsgs []string
		for _, err := range result.Errors {
			errMsgs = append(errMsgs, err.Error())
		}
		return fmt.Sprintf("Errors:\n%s", strings.Join(errMsgs, "\n"))
	}

	var output strings.Builder
	output.WriteString("Simulation Result:\n")

	for name, bits := range result.ClassicalRegisters {
		output.WriteString(fmt.Sprintf("Register %s: ", name))
		
		// 将布尔值转换为0和1
		var bitValues []string
		for _, bit := range bits {
			if bit {
				bitValues = append(bitValues, "1")
			} else {
				bitValues = append(bitValues, "0")
			}
		}
		
		// 打印二进制值
		output.WriteString(strings.Join(bitValues, ""))
		
		// 计算十进制值
		decimalValue := 0
		for i, bit := range bits {
			if bit {
				decimalValue += 1 << (len(bits) - i - 1)
			}
		}
		
		output.WriteString(fmt.Sprintf(" (decimal: %d)\n", decimalValue))
	}

	return output.String()
}

// GetVersion 获取模拟器使用的QASM版本
func (s *Simulator) GetVersion() parser.QASMVersion {
	return s.parser.GetVersion()
}

// SetVersion 设置模拟器使用的QASM版本
func (s *Simulator) SetVersion(version parser.QASMVersion) {
	s.parser = parser.NewParser(version)
}
