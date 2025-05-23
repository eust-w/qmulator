package parser

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// parseSymbolicExpression 解析符号表达式，如pi/2, pi*2等
func parseSymbolicExpression(expr string) (float64, error) {
	// 去除空格
	expr = strings.TrimSpace(expr)
	
	// 尝试直接解析为浮点数
	if val, err := strconv.ParseFloat(expr, 64); err == nil {
		return val, nil
	}
	
	// 替换常量
	expr = strings.ReplaceAll(expr, "pi", fmt.Sprintf("%f", math.Pi))
	expr = strings.ReplaceAll(expr, "PI", fmt.Sprintf("%f", math.Pi))
	expr = strings.ReplaceAll(expr, "Pi", fmt.Sprintf("%f", math.Pi))
	
	// 处理简单的除法表达式，如pi/2
	divRegex := regexp.MustCompile(`([0-9.]+)\s*/\s*([0-9.]+)`)
	if matches := divRegex.FindStringSubmatch(expr); len(matches) == 3 {
		numerator, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid numerator: %s", matches[1])
		}
		
		denominator, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid denominator: %s", matches[2])
		}
		
		if denominator == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		
		return numerator / denominator, nil
	}
	
	// 处理简单的乘法表达式，如pi*2
	mulRegex := regexp.MustCompile(`([0-9.]+)\s*\*\s*([0-9.]+)`)
	if matches := mulRegex.FindStringSubmatch(expr); len(matches) == 3 {
		factor1, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid factor: %s", matches[1])
		}
		
		factor2, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid factor: %s", matches[2])
		}
		
		return factor1 * factor2, nil
	}
	
	// 处理简单的加法表达式，如pi+2
	addRegex := regexp.MustCompile(`([0-9.]+)\s*\+\s*([0-9.]+)`)
	if matches := addRegex.FindStringSubmatch(expr); len(matches) == 3 {
		term1, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid term: %s", matches[1])
		}
		
		term2, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid term: %s", matches[2])
		}
		
		return term1 + term2, nil
	}
	
	// 处理简单的减法表达式，如pi-2
	subRegex := regexp.MustCompile(`([0-9.]+)\s*-\s*([0-9.]+)`)
	if matches := subRegex.FindStringSubmatch(expr); len(matches) == 3 {
		term1, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid term: %s", matches[1])
		}
		
		term2, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid term: %s", matches[2])
		}
		
		return term1 - term2, nil
	}
	
	// 如果以上都不匹配，尝试直接解析
	return strconv.ParseFloat(expr, 64)
}
