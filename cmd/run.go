package cmd

import (
	"fmt"
	"os"

	"github.com/d-robotics/qmulator/pkg/parser"
	"github.com/d-robotics/qmulator/pkg/simulator"
	"github.com/spf13/cobra"
)

var (
	runFile    string
	qasmVersion string
)

// runCmd 表示运行命令
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行QASM量子汇编代码",
	Long:  `运行QASM量子汇编代码文件，并显示模拟结果。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查文件参数
		if runFile == "" {
			fmt.Println("错误: 必须指定QASM文件路径")
			cmd.Help()
			os.Exit(1)
		}

		// 检查文件是否存在
		if _, err := os.Stat(runFile); os.IsNotExist(err) {
			fmt.Printf("错误: 文件 %s 不存在\n", runFile)
			os.Exit(1)
		}

		// 确定QASM版本
		version := parser.QASM2 // 默认使用QASM 2.0
		if qasmVersion != "" {
			switch qasmVersion {
			case "1.0":
				version = parser.QASM1
			case "2.0":
				version = parser.QASM2
			default:
				fmt.Printf("错误: 不支持的QASM版本 %s，使用默认版本 2.0\n", qasmVersion)
			}
		}

		// 创建模拟器
		sim := simulator.NewSimulator(version)

		// 运行QASM文件
		fmt.Printf("运行QASM文件: %s (版本: %s)\n", runFile, version)
		result, err := sim.RunFile(runFile)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			os.Exit(1)
		}

		// 打印结果
		fmt.Println(sim.FormatResult(result))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// 添加命令行参数
	runCmd.Flags().StringVarP(&runFile, "file", "f", "", "QASM文件路径")
	runCmd.Flags().StringVarP(&qasmVersion, "version", "v", "", "QASM版本 (1.0 或 2.0)")
}
