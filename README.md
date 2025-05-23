# Qmulator - QASM量子汇编模拟器

Qmulator是一个基于Go语言的QASM量子汇编模拟器，支持QASM 1.0和QASM 2.0标准。该项目使用Cobra实现命令行功能，使用Gin框架提供REST API接口。

## 功能特点

- 支持QASM 1.0和QASM 2.0语法解析
- 模拟执行量子电路
- 命令行界面支持直接执行QASM文件
- REST API接口支持远程调用
- 可视化量子电路执行结果

## 安装

确保您已安装Go 1.16或更高版本，然后运行：

```bash
go get github.com/d-robotics/qmulator
```

## 使用方法

### 命令行使用

```bash
# 运行QASM文件
qmulator run -f path/to/file.qasm

# 指定QASM版本
qmulator run -f path/to/file.qasm --version 2.0

# 启动API服务器
qmulator serve --port 8080
```

### API使用

启动服务器后，可以通过以下API端点访问：

- `POST /api/v1/run` - 运行QASM代码
- `GET /api/v1/status` - 检查服务器状态

## 项目结构

```
.
├── cmd/           # 命令行应用程序
├── pkg/           # 核心包
│   ├── simulator/ # 量子模拟器实现
│   ├── parser/    # QASM解析器
│   └── api/       # REST API实现
└── examples/      # QASM示例文件
```

## 许可证

MIT