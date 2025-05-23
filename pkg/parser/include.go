package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// IncludeManager 管理包含文件的处理
type IncludeManager struct {
	includePaths []string
	loadedFiles  map[string]bool
}

// NewIncludeManager 创建一个新的包含文件管理器
func NewIncludeManager(includePaths []string) *IncludeManager {
	return &IncludeManager{
		includePaths: includePaths,
		loadedFiles:  make(map[string]bool),
	}
}

// LoadIncludeFile 加载包含文件
func (im *IncludeManager) LoadIncludeFile(filename string) (string, error) {
	// 检查文件是否已加载，防止循环包含
	if im.loadedFiles[filename] {
		return "", nil
	}

	// 在包含路径中查找文件
	var filePath string
	var found bool

	// 首先检查是否是绝对路径
	if filepath.IsAbs(filename) {
		if _, err := os.Stat(filename); err == nil {
			filePath = filename
			found = true
		}
	} else {
		// 检查当前目录
		if _, err := os.Stat(filename); err == nil {
			absPath, err := filepath.Abs(filename)
			if err == nil {
				filePath = absPath
				found = true
			}
		}

		// 检查包含路径
		if !found {
			for _, path := range im.includePaths {
				fullPath := filepath.Join(path, filename)
				if _, err := os.Stat(fullPath); err == nil {
					filePath = fullPath
					found = true
					break
				}
			}
		}
	}

	if !found {
		return "", fmt.Errorf("include file not found: %s", filename)
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read include file: %v", err)
	}

	// 标记文件为已加载
	im.loadedFiles[filename] = true

	return string(content), nil
}

// ProcessIncludes 处理QASM代码中的包含语句
func (im *IncludeManager) ProcessIncludes(code string) (string, error) {
	lines := strings.Split(code, "\n")
	var processedLines []string

	// 添加空行以确保格式正确
	processedLines = append(processedLines, "")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// 检查是否是包含语句
		if strings.HasPrefix(trimmedLine, "include") {
			includeFile, err := parseInclude(trimmedLine)
			if err != nil {
				return "", err
			}

			// 加载包含文件
			includeContent, err := im.LoadIncludeFile(includeFile)
			if err != nil {
				return "", err
			}

			// 添加注释标记包含文件的开始
			processedLines = append(processedLines, fmt.Sprintf("// BEGIN INCLUDE: %s", includeFile))
			
			// 直接添加包含文件的内容，不进行递归处理
			// 这样可以确保门定义的完整性
			includeLines := strings.Split(includeContent, "\n")
			for _, includeLine := range includeLines {
				processedLines = append(processedLines, includeLine)
			}
			
			// 添加注释标记包含文件的结束
			processedLines = append(processedLines, fmt.Sprintf("// END INCLUDE: %s", includeFile))
			
			// 添加空行以确保格式正确
			processedLines = append(processedLines, "")
		} else {
			// 保留原始行
			processedLines = append(processedLines, line)
		}
	}

	return strings.Join(processedLines, "\n"), nil
}
