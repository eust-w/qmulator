package cmd

import (
	"fmt"

	"github.com/d-robotics/qmulator/pkg/api"
	"github.com/spf13/cobra"
)

var (
	serverPort int
	enableUI   bool
)

// serveCmd 表示服务命令
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动QASM模拟器API服务器",
	Long:  `启动QASM模拟器API服务器，提供REST API接口。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("启动API服务器，监听端口: %d\n", serverPort)
		
		// 创建API服务器
		server := api.NewServer(serverPort, enableUI)
		
		// 如果启用UI，显示相关信息
		if enableUI {
			fmt.Printf("启用Web界面，请访问: http://localhost:%d/ui\n", serverPort)
			fmt.Println("支持功能：")
			fmt.Println("  - 在线运行QASM代码")
			fmt.Println("  - 量子电路可视化")
			fmt.Println("  - 步骤运行图表")
			fmt.Println("  - 支持QASM 1.0和2.0版本")
		}
		
		// 启动服务器
		if err := server.Start(); err != nil {
			fmt.Printf("服务器启动失败: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// 添加命令行参数
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "API服务器端口")
	serveCmd.Flags().BoolVar(&enableUI, "ui", false, "启用Web界面，支持在线运行QASM代码和可视化")
}
