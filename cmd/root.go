package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qmulator",
	Short: "QASM量子汇编模拟器",
	Long: `Qmulator是一个基于Go语言的QASM量子汇编模拟器，
支持QASM 1.0和QASM 2.0标准。该项目使用Cobra实现命令行功能，
使用Gin框架提供REST API接口。`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
