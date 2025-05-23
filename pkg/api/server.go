package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/d-robotics/qmulator/pkg/parser"
	"github.com/d-robotics/qmulator/pkg/simulator"
	"github.com/gin-gonic/gin"
)

// Server 表示API服务器
type Server struct {
	engine    *gin.Engine
	simulator *simulator.Simulator
	port      int
	enableUI  bool
}

// RunRequest 表示运行QASM代码的请求
type RunRequest struct {
	Code    string `json:"code" binding:"required"`
	Version string `json:"version,omitempty"`
}

// RunResponse 表示运行QASM代码的响应
type RunResponse struct {
	Success   bool                `json:"success"`
	Registers map[string][]bool   `json:"registers,omitempty"`
	Decimal   map[string]int      `json:"decimal,omitempty"`
	Errors    []string            `json:"errors,omitempty"`
}

// FormatRequest 表示格式化QASM代码的请求
type FormatRequest struct {
	Code string `json:"code" binding:"required"`
}

// FormatResponse 表示格式化QASM代码的响应
type FormatResponse struct {
	FormattedCode string `json:"formattedCode,omitempty"`
	Error         string `json:"error,omitempty"`
}

// NewServer 创建一个新的API服务器
func NewServer(port int, enableUI bool) *Server {
	engine := gin.Default()
	
	// 如果启用UI，设置HTML模板渲染器
	if enableUI {
		engine.LoadHTMLGlob("web/templates/*html")
	}
	
	server := &Server{
		engine:    engine,
		simulator: simulator.NewSimulator(parser.QASM2), // 默认使用QASM 2.0
		port:      port,
		enableUI:  enableUI,
	}
	
	server.setupRoutes()
	
	return server
}

// setupRoutes 设置API路由
func (s *Server) setupRoutes() {
	// API路由
	api := s.engine.Group("/api/v1")
	{
		api.GET("/status", s.handleStatus)
		api.POST("/run", s.handleRun)
		api.POST("/format", s.handleFormat)
		api.GET("/versions", s.handleGetVersions)
	}
	
	// 如果启用UI，设置Web界面路由
	if s.enableUI {
		// 静态文件服务
		s.engine.Static("/assets", "./web/assets")
		
		// 注册UI路由
		RegisterUIHandlers(s.engine, s)
	}
}

// handleStatus 处理状态请求
func (s *Server) handleStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
		"version": string(s.simulator.GetVersion()),
	})
}

// handleRun 处理运行QASM代码的请求
func (s *Server) handleRun(c *gin.Context) {
	var req RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 设置QASM版本（如果指定）
	if req.Version != "" {
		version := parser.QASMVersion(req.Version)
		if version != parser.QASM1 && version != parser.QASM2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported QASM version"})
			return
		}
		s.simulator.SetVersion(version)
	}
	
	// 验证输入
	if req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty QASM code"})
		return
	}

	// 运行QASM代码
	result, err := s.simulator.RunString(req.Code)
	if err != nil {
		// 打印详细错误信息以便调试
		fmt.Printf("Error running QASM code: %v\n", err)
		
		// 将语法错误返回为400，其他错误返回为500
		if strings.Contains(err.Error(), "syntax") || strings.Contains(err.Error(), "parse") {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("QASM syntax error: %v", err)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to run QASM code: %v", err)})
		}
		return
	}
	
	// 构建响应
	response := RunResponse{
		Success:   len(result.Errors) == 0,
		Registers: make(map[string][]bool),
		Decimal:   make(map[string]int),
	}
	
	if len(result.Errors) > 0 {
		response.Errors = make([]string, len(result.Errors))
		for i, err := range result.Errors {
			response.Errors[i] = err.Error()
		}
	} else {
		for name, bits := range result.ClassicalRegisters {
			response.Registers[name] = bits
			
			// 计算十进制值
			decimalValue := 0
			for i, bit := range bits {
				if bit {
					decimalValue += 1 << (len(bits) - i - 1)
				}
			}
			response.Decimal[name] = decimalValue
		}
	}
	
	c.JSON(http.StatusOK, response)
}

// handleFormat 处理格式化QASM代码的请求
func (s *Server) handleFormat(c *gin.Context) {
	var req FormatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse{Error: err.Error()})
		return
	}
	
	// 格式化QASM代码
	formattedCode, err := s.formatQASMCode(req.Code)
	if err != nil {
		c.JSON(http.StatusOK, FormatResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, FormatResponse{FormattedCode: formattedCode})
}

// formatQASMCode 格式化QASM代码
func (s *Server) formatQASMCode(code string) (string, error) {
	// 分行处理
	lines := strings.Split(code, "\n")
	formattedLines := make([]string, 0, len(lines))
	
	// 跟踪空行
	var lastLineWasEmpty bool
	
	// 处理每一行
	for _, line := range lines {
		// 去除前后空白
		trimmedLine := strings.TrimSpace(line)
		
		// 处理空行，但避免多个连续空行
		if trimmedLine == "" {
			// 只有当前一行不是空行时才添加空行
			if !lastLineWasEmpty {
				formattedLines = append(formattedLines, "")
				lastLineWasEmpty = true
			}
			continue
		}
		
		// 非空行
		lastLineWasEmpty = false
		
		// 处理注释，保留原样
		if strings.HasPrefix(trimmedLine, "//") {
			formattedLines = append(formattedLines, trimmedLine)
			continue
		}
		
		// 处理OPENQASM版本声明
		if strings.HasPrefix(trimmedLine, "OPENQASM") {
			// 删除可能存在的分号
			if strings.HasSuffix(trimmedLine, ";") {
				trimmedLine = trimmedLine[:len(trimmedLine)-1]
			}
			formattedLines = append(formattedLines, trimmedLine + ";")
			continue
		}
		
		// 处理include语句
		if strings.HasPrefix(trimmedLine, "include") {
			// 删除可能存在的分号
			if strings.HasSuffix(trimmedLine, ";") {
				trimmedLine = trimmedLine[:len(trimmedLine)-1]
			}
			formattedLines = append(formattedLines, trimmedLine + ";")
			continue
		}
		
		// 处理量子寄存器和经典寄存器声明
		if strings.HasPrefix(trimmedLine, "qreg") || strings.HasPrefix(trimmedLine, "creg") {
			// 删除可能存在的分号
			if strings.HasSuffix(trimmedLine, ";") {
				trimmedLine = trimmedLine[:len(trimmedLine)-1]
			}
			formattedLines = append(formattedLines, trimmedLine + ";")
			continue
		}
		
		// 处理量子门操作
		// 删除可能存在的分号
		if strings.HasSuffix(trimmedLine, ";") {
			trimmedLine = trimmedLine[:len(trimmedLine)-1]
		}
		
		// 处理测量操作
		if strings.HasPrefix(trimmedLine, "measure") {
			formattedLines = append(formattedLines, trimmedLine + ";")
			continue
		}
		
		// 处理常规量子门
		// 分割操作符和参数，使用Fields处理多个空格
		parts := strings.Fields(trimmedLine)
		if len(parts) > 0 {
			// 处理带参数的门，如rx(pi/2)
			// 保持原有格式，只添加分号
			formattedLines = append(formattedLines, trimmedLine + ";")
		} else {
			// 如果没有内容，直接跳过
			continue
		}
	}
	
	// 返回格式化后的代码
	return strings.Join(formattedLines, "\n"), nil
}

// Start 启动API服务器
func (s *Server) Start() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.port))
}
