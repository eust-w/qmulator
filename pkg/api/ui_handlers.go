package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/eust-w/qmulator/pkg/parser"
	"github.com/gin-gonic/gin"
)

// CircuitData 表示量子电路数据
type CircuitData struct {
	Gates     []GateData     `json:"gates"`
	Qubits    int            `json:"qubits"`
	Bits      int            `json:"bits"`
	Registers []RegisterData `json:"registers"`
}

// GateData 表示量子门数据
type GateData struct {
	Name      string   `json:"name"`
	Targets   []int    `json:"targets"`
	Controls  []int    `json:"controls"`
	Params    []string `json:"params,omitempty"`
	Position  int      `json:"position"`
	QubitsMap []int    `json:"qubitsMap,omitempty"`
}

// RegisterData 表示寄存器数据
type RegisterData struct {
	Name  string `json:"name"`
	Size  int    `json:"size"`
	IsQreg bool  `json:"isQreg"`
	Start int    `json:"start"`
}

// VersionsResponse 表示QASM版本响应
type VersionsResponse struct {
	Versions []string `json:"versions"`
}

// ExampleData 表示示例代码数据
type ExampleData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Version     string `json:"version"`
}

// RegisterUIHandlers 注册UI相关的路由处理器
func RegisterUIHandlers(router *gin.Engine, s *Server) {
	router.GET("/ui", s.handleUI)
	router.GET("/ui/simulator", s.handleSimulatorUI)
	router.GET("/ui/circuit", s.handleCircuitUI)
	router.GET("/ui/circuit_designer", s.handleCircuitDesignerUI)
	router.GET("/ui/examples", s.handleExamplesUI)
	router.GET("/ui/circuit/:id", s.handleCircuitDetail)
	router.GET("/api/v1/circuit", s.handleGetCircuit)
	router.POST("/api/v1/circuit", s.handleGetCircuit)
}

// handleUI 处理Web界面主页请求
func (s *Server) handleUI(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Qmulator - QASM量子汇编模拟器",
	})
}

// handleCircuitUI 处理电路可视化界面请求
func (s *Server) handleCircuitUI(c *gin.Context) {
	c.HTML(http.StatusOK, "circuit.html", gin.H{
		"title": "量子电路可视化",
	})
}

// handleSimulatorUI 处理模拟器界面请求
func (s *Server) handleSimulatorUI(c *gin.Context) {
	c.HTML(http.StatusOK, "simulator.html", gin.H{
		"title": "QASM模拟器",
	})
}

// handleExamplesUI 处理示例代码界面请求
func (s *Server) handleExamplesUI(c *gin.Context) {
	c.HTML(http.StatusOK, "examples.html", gin.H{
		"title": "QASM示例代码",
	})
}

// handleCircuitDesignerUI 处理电路设计器界面请求
func (s *Server) handleCircuitDesignerUI(c *gin.Context) {
	c.HTML(http.StatusOK, "circuit_designer.html", gin.H{
		"title": "量子电路设计器",
	})
}

// handleCircuitDetail 处理电路详情页面请求
func (s *Server) handleCircuitDetail(c *gin.Context) {
	circuitID := c.Param("id")
	// 这里可以根据电路 ID 加载特定电路的数据
	// 在实际应用中，可能需要从数据库或其他存储中获取电路数据
	c.HTML(http.StatusOK, "circuit_detail.html", gin.H{
		"title": "电路详情",
		"circuitID": circuitID,
	})
}

// handleGetCircuit 处理获取电路数据请求
func (s *Server) handleGetCircuit(c *gin.Context) {
	var req RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析QASM代码，提取电路数据
	circuit, err := s.extractCircuitData(req.Code, req.Version)
	if err != nil {
		// 检查错误类型，分类处理
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()
		
		// 如果是语法错误、输入错误或不支持的门，返回400
		if strings.Contains(errorMessage, "syntax error") || 
		   strings.Contains(errorMessage, "empty QASM code") ||
		   strings.Contains(errorMessage, "failed to run QASM code") ||
		   strings.Contains(errorMessage, "unsupported gates") ||
		   strings.Contains(errorMessage, "unclosed gate definition") {
			statusCode = http.StatusBadRequest
			
			// 如果是不支持的门，提供更友好的错误信息
			if strings.Contains(errorMessage, "unsupported gates") {
				// 提取不支持的门列表
				unsupportedGates := strings.TrimPrefix(errorMessage, "unsupported gates detected: ")
				errorMessage = fmt.Sprintf(
					"您的代码使用了当前模拟器不支持的量子门: %s。请使用基本门（如x, y, z, h, cx等）替代。",
					unsupportedGates,
				)
			} else if strings.Contains(errorMessage, "unclosed gate definition") {
				// 提供更友好的错误信息给用户
				errorMessage = "自定义门定义未关闭。请确保每个门定义都有匹配的大括号。例如:\n\ngate mygate a,b {\n  // 门操作\n}  // 确保有这个结束的大括号"
			}
		}
		
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, circuit)
}

// handleGetVersions 处理获取支持的QASM版本请求
func (s *Server) handleGetVersions(c *gin.Context) {
	versions := VersionsResponse{
		Versions: []string{"1.0", "2.0"},
	}
	c.JSON(http.StatusOK, versions)
}

// isGateSupported 检查量子门是否受支持
func isGateSupported(gateName string) bool {
	// 支持的基本量子门列表
	// 注意：这个列表与模拟器实际支持的门一致
	supportedGates := map[string]bool{
		// 基本门
		"x": true, "y": true, "z": true, "h": true, "s": true, "t": true,
		"sdg": true, "tdg": true, "id": true,
		// 旋转门
		"rx": true, "ry": true, "rz": true,
		// 双量子门
		"cx": true, "cz": true, "cy": true, "ch": true, "swap": true,
		// 三量子门
		"ccx": true, "toffoli": true, "cswap": true, "fredkin": true,
		// 测量
		"measure": true, "reset": true,
		// 参数化门
		"u": true, "u1": true, "u2": true, "u3": true, "p": true,
		// 受控参数化门
		"cu1": true, "cu2": true, "cu3": true, "crz": true,
	}

	_, supported := supportedGates[strings.ToLower(gateName)]
	return supported
}

// extractCircuitData 从QASM代码中提取电路数据
func (s *Server) extractCircuitData(code string, version string) (*CircuitData, error) {
	// 验证输入
	if code == "" {
		return nil, fmt.Errorf("empty QASM code")
	}
	
	// 检查代码中是否使用了不支持的门
	unsupportedGates := []string{}
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		// 去除注释和空白
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 跳过OPENQASM、include、qreg、creg等非门操作行
		if strings.HasPrefix(line, "OPENQASM") || 
		   strings.HasPrefix(line, "include") || 
		   strings.HasPrefix(line, "qreg") || 
		   strings.HasPrefix(line, "creg") {
			continue
		}
		
		// 提取门名
		parts := strings.Split(line, " ")
		if len(parts) > 0 {
			// 获取门名，处理可能的参数
			gateName := strings.TrimSpace(parts[0])
			
			// 如果是自定义门定义，跳过
			if strings.HasPrefix(gateName, "gate") {
				continue
			}
			
			// 处理带参数的门调用格式，如cu1(pi/2)
			if parenthesisIndex := strings.Index(gateName, "("); parenthesisIndex > 0 {
				gateName = gateName[:parenthesisIndex]
			}
			
			// 检查门是否受支持
			if gateName != "" && !isGateSupported(gateName) {
				// 将原始的门调用格式添加到不支持列表
				unsupportedGates = append(unsupportedGates, strings.TrimSpace(parts[0]))
			}
		}
	}
	
	// 如果有不支持的门，返回错误
	if len(unsupportedGates) > 0 {
		return nil, fmt.Errorf("unsupported gates detected: %s", strings.Join(unsupportedGates, ", "))
	}

	// 设置QASM版本
	qasmVersion := parser.QASM2
	if version == "1.0" {
		qasmVersion = parser.QASM1
	}

	// 创建模拟器实例
	sim := s.simulator
	sim.SetVersion(qasmVersion)

	// 尝试解析QASM代码
	result, err := sim.RunString(code)
	if err != nil {
		// 打印详细错误信息以便调试
		fmt.Printf("Error running QASM code: %v\n", err)
		return nil, fmt.Errorf("failed to run QASM code: %v", err)
	}

	// 检查解析错误
	if len(result.Errors) > 0 {
		// 打印详细错误信息以便调试
		fmt.Printf("QASM parse errors: %v\n", result.Errors)
		return nil, fmt.Errorf("QASM syntax error: %v", result.Errors[0])
	}

	// 创建电路数据
	circuit := &CircuitData{
		Gates:     []GateData{},
		Qubits:    0,
		Bits:      0,
		Registers: []RegisterData{},
	}

	// 从代码中提取量子寄存器和经典寄存器信息
	qregPattern := regexp.MustCompile(`qreg\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	cregPattern := regexp.MustCompile(`creg\s+([a-zA-Z][a-zA-Z0-9]*)\[([0-9]+)\]`)
	
	// 存储所有寄存器
	qregs := make(map[string]int)
	cregs := make(map[string]int)
	
	// 扫描代码行寻找寄存器声明
	codeLines := strings.Split(code, "\n")
	for _, line := range codeLines {
		// 去除注释和空白
		var idx int
		if idx = strings.Index(line, "//"); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 匹配量子寄存器
		if qregMatches := qregPattern.FindStringSubmatch(line); len(qregMatches) == 3 {
			regName := qregMatches[1]
			regSize, _ := strconv.Atoi(qregMatches[2])
			qregs[regName] = regSize
		}
		
		// 匹配经典寄存器
		if cregMatches := cregPattern.FindStringSubmatch(line); len(cregMatches) == 3 {
			regName := cregMatches[1]
			regSize, _ := strconv.Atoi(cregMatches[2])
			cregs[regName] = regSize
		}
	}
	
	// 计算总量子比特数和经典比特数
	totalQubits := 0
	totalBits := 0
	qregStart := 0
	cregStart := 0
	
	// 添加量子寄存器
	for name, size := range qregs {
		circuit.Registers = append(circuit.Registers, RegisterData{
			Name:   name,
			Size:   size,
			IsQreg: true,
			Start:  qregStart,
		})
		totalQubits += size
		qregStart += size
	}
	
	// 添加经典寄存器
	for name, size := range cregs {
		circuit.Registers = append(circuit.Registers, RegisterData{
			Name:   name,
			Size:   size,
			IsQreg: false,
			Start:  cregStart,
		})
		totalBits += size
		cregStart += size
	}
	
	// 设置总量子比特数和经典比特数
	circuit.Qubits = totalQubits
	circuit.Bits = totalBits
	
	// 如果没有检测到寄存器，使用默认值
	if totalQubits == 0 {
		circuit.Qubits = 3
		circuit.Registers = append(circuit.Registers, RegisterData{
			Name:   "q",
			Size:   3,
			IsQreg: true,
			Start:  0,
		})
	}
	
	if totalBits == 0 && totalQubits > 0 {
		circuit.Bits = circuit.Qubits
		circuit.Registers = append(circuit.Registers, RegisterData{
			Name:   "c",
			Size:   circuit.Qubits,
			IsQreg: false,
			Start:  0,
		})
	}

	// 直接从代码中提取门操作
	position := 0
	
	// 定义各种门的正则表达式
	
	// 1. 单量子比特门（h, x, y, z, s, t等）
	singleQubitGatePattern := regexp.MustCompile(`([hxyzstid]+)\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	
	// 2. 双量子比特门（cx, cz, swap等）
	twoQubitGatePattern := regexp.MustCompile(`(cx|cz|cy|ch|swap)\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	
	// 3. 带参数的受控门（cu1, cu2, cu3, crz等）
	controlledParamGatePattern := regexp.MustCompile(`(cu1|cu2|cu3|crz)\s*\(([^)]*)\)\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	
	// 4. 旋转门（rx, ry, rz, u1, u2, u3等）
	rotationGatePattern := regexp.MustCompile(`(rx|ry|rz|u1|u2|u3)\s*\(([^)]*)\)\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	
	// 5. 测量门
	measurePattern := regexp.MustCompile(`measure\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]\s+->\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	measureAllPattern := regexp.MustCompile(`measure\s+([a-zA-Z][a-zA-Z0-9]*)\s+->\s+([a-zA-Z][a-zA-Z0-9]*)`)
	
	// 6. 多量子比特门（如ccx, cswap等）
	multiQubitGatePattern := regexp.MustCompile(`(ccx|cswap|fredkin|toffoli)\s+([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]\s*,\s*([a-zA-Z][a-zA-Z0-9]*)\[(\d+)\]`)
	
	// 遍历代码行
	for _, line := range codeLines {
		// 去除注释和空白
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "OPENQASM") || strings.HasPrefix(line, "include") || 
		   strings.HasPrefix(line, "qreg") || strings.HasPrefix(line, "creg") {
			continue
		}
		
		// 1. 匹配单量子比特门
		if matches := singleQubitGatePattern.FindStringSubmatch(line); len(matches) == 4 {
			gateName := matches[1]
			regName := matches[2]
			qubitIdx, _ := strconv.Atoi(matches[3])
			
			// 计算实际的量子比特索引
			globalIdx := -1
			for _, reg := range circuit.Registers {
				if reg.Name == regName && reg.IsQreg {
					globalIdx = reg.Start + qubitIdx
					break
				}
			}
			
			if globalIdx >= 0 {
				circuit.Gates = append(circuit.Gates, GateData{
					Name:     gateName,
					Targets:  []int{globalIdx},
					Controls: []int{},
					Position: position,
				})
				position++
			}
			continue
		}
		
		// 2. 匹配双量子比特门
		if matches := twoQubitGatePattern.FindStringSubmatch(line); len(matches) == 6 {
			gateName := matches[1]
			controlRegName := matches[2]
			controlIdx, _ := strconv.Atoi(matches[3])
			targetRegName := matches[4]
			targetIdx, _ := strconv.Atoi(matches[5])
			
			// 计算实际的量子比特索引
			controlGlobalIdx := -1
			targetGlobalIdx := -1
			
			for _, reg := range circuit.Registers {
				if reg.Name == controlRegName && reg.IsQreg {
					controlGlobalIdx = reg.Start + controlIdx
				}
				if reg.Name == targetRegName && reg.IsQreg {
					targetGlobalIdx = reg.Start + targetIdx
				}
			}
			
			if controlGlobalIdx >= 0 && targetGlobalIdx >= 0 {
				// 对于swap门，没有控制比特
				if gateName == "swap" {
					circuit.Gates = append(circuit.Gates, GateData{
						Name:     gateName,
						Targets:  []int{controlGlobalIdx, targetGlobalIdx},
						Controls: []int{},
						Position: position,
					})
				} else {
					circuit.Gates = append(circuit.Gates, GateData{
						Name:     gateName,
						Targets:  []int{targetGlobalIdx},
						Controls: []int{controlGlobalIdx},
						Position: position,
					})
				}
				position++
			}
			continue
		}
		
		// 3. 匹配带参数的受控门
		if matches := controlledParamGatePattern.FindStringSubmatch(line); len(matches) == 7 {
			gateName := matches[1]
			params := strings.Split(matches[2], ",")
			for i := range params {
				params[i] = strings.TrimSpace(params[i])
			}
			controlRegName := matches[3]
			controlIdx, _ := strconv.Atoi(matches[4])
			targetRegName := matches[5]
			targetIdx, _ := strconv.Atoi(matches[6])
			
			// 计算实际的量子比特索引
			controlGlobalIdx := -1
			targetGlobalIdx := -1
			
			for _, reg := range circuit.Registers {
				if reg.Name == controlRegName && reg.IsQreg {
					controlGlobalIdx = reg.Start + controlIdx
				}
				if reg.Name == targetRegName && reg.IsQreg {
					targetGlobalIdx = reg.Start + targetIdx
				}
			}
			
			if controlGlobalIdx >= 0 && targetGlobalIdx >= 0 {
				circuit.Gates = append(circuit.Gates, GateData{
					Name:     gateName,
					Targets:  []int{targetGlobalIdx},
					Controls: []int{controlGlobalIdx},
					Params:   params,
					Position: position,
				})
				position++
			}
			continue
		}
		
		// 4. 匹配旋转门
		if matches := rotationGatePattern.FindStringSubmatch(line); len(matches) == 5 {
			gateName := matches[1]
			params := strings.Split(matches[2], ",")
			for i := range params {
				params[i] = strings.TrimSpace(params[i])
			}
			regName := matches[3]
			qubitIdx, _ := strconv.Atoi(matches[4])
			
			// 计算实际的量子比特索引
			globalIdx := -1
			for _, reg := range circuit.Registers {
				if reg.Name == regName && reg.IsQreg {
					globalIdx = reg.Start + qubitIdx
					break
				}
			}
			
			if globalIdx >= 0 {
				circuit.Gates = append(circuit.Gates, GateData{
					Name:     gateName,
					Targets:  []int{globalIdx},
					Controls: []int{},
					Params:   params,
					Position: position,
				})
				position++
			}
			continue
		}
		
		// 5. 匹配测量门 (单比特测量)
		if matches := measurePattern.FindStringSubmatch(line); len(matches) == 5 {
			qregName := matches[1]
			qubitIdx, _ := strconv.Atoi(matches[2])
			
			// 计算实际的量子比特索引
			globalIdx := -1
			for _, reg := range circuit.Registers {
				if reg.Name == qregName && reg.IsQreg {
					globalIdx = reg.Start + qubitIdx
					break
				}
			}
			
			if globalIdx >= 0 {
				circuit.Gates = append(circuit.Gates, GateData{
					Name:     "measure",
					Targets:  []int{globalIdx},
					Controls: []int{},
					Position: position,
				})
				position++
			}
			continue
		}
		
		// 6. 匹配测量门 (整个寄存器测量)
		if matches := measureAllPattern.FindStringSubmatch(line); len(matches) == 3 {
			qregName := matches[1]
			
			// 收集该寄存器的所有量子比特
			var qubits []int
			for _, reg := range circuit.Registers {
				if reg.Name == qregName && reg.IsQreg {
					for i := 0; i < reg.Size; i++ {
						qubits = append(qubits, reg.Start + i)
					}
					break
				}
			}
			
			if len(qubits) > 0 {
				circuit.Gates = append(circuit.Gates, GateData{
					Name:     "measure",
					Targets:  qubits,
					Controls: []int{},
					Position: position,
				})
				position++
			}
			continue
		}
		
		// 7. 匹配多量子比特门（如ccx, cswap等）
		if matches := multiQubitGatePattern.FindStringSubmatch(line); len(matches) == 8 {
			gateName := matches[1]
			control1RegName := matches[2]
			control1Idx, _ := strconv.Atoi(matches[3])
			control2RegName := matches[4]
			control2Idx, _ := strconv.Atoi(matches[5])
			targetRegName := matches[6]
			targetIdx, _ := strconv.Atoi(matches[7])
			
			// 计算实际的量子比特索引
			control1GlobalIdx := -1
			control2GlobalIdx := -1
			targetGlobalIdx := -1
			
			for _, reg := range circuit.Registers {
				if reg.Name == control1RegName && reg.IsQreg {
					control1GlobalIdx = reg.Start + control1Idx
				}
				if reg.Name == control2RegName && reg.IsQreg {
					control2GlobalIdx = reg.Start + control2Idx
				}
				if reg.Name == targetRegName && reg.IsQreg {
					targetGlobalIdx = reg.Start + targetIdx
				}
			}
			
			if control1GlobalIdx >= 0 && control2GlobalIdx >= 0 && targetGlobalIdx >= 0 {
				circuit.Gates = append(circuit.Gates, GateData{
					Name:     gateName,
					Targets:  []int{targetGlobalIdx},
					Controls: []int{control1GlobalIdx, control2GlobalIdx},
					Position: position,
				})
				position++
			}
			continue
		}
	}

	return circuit, nil
}
