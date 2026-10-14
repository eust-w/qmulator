<div align="center">
  <img src="docs/readme/hero.svg" width="100%" alt="Qmulator workflow overview">
  <h1>Qmulator</h1>
  <p><strong>用轻量教学原型探索 QASM 程序与电路。</strong></p>
  <p><a href="LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-d99536?style=flat-square"></a> <a href="docs/readme/validation.md"><img alt="Validation notes" src="https://img.shields.io/badge/Validation-repository%20checks-547ac4?style=flat-square"></a></p>
  <p><a href="#overview">项目亮点</a> · <a href="#quick-start">快速开始</a> · <a href="#workflow">工作流</a> · <a href="#validation">验证</a> · <a href="#documentation">文档</a></p>
  <p><a href="README.en.md">English</a> · <a href="README.md">简体中文</a></p>
</div>

---

<a id="overview"></a>
## 项目亮点

Go QASM 解析器、CLI 与 Web 教学界面。当前模拟器采用独立量子位表示，并非完整联合态量子引擎。

| 能力 | 内容 |
| :--- | :--- |
| **QASM 探索** | 解析器与示例程序，用于查看已实现语法。 |
| **CLI / API** | 运行文件或查看本地 HTTP 接口。 |
| **电路界面** | 浏览器电路示例与可视化。 |

### 本地页面

![qmulator local UI](docs/readme/demo.jpg)

<a id="quick-start"></a>
## 快速开始

```bash
git clone https://github.com/eust-w/qmulator.git
cd qmulator
go mod download
go build -o qmulator .
./qmulator run -f examples/bell_state.qasm
./qmulator serve --port 8080 --ui
# Open http://127.0.0.1:8080/ui
```

<a id="workflow"></a>
## 工作流

```text
PARSE → EXECUTE → INSPECT
```

```bash
go test ./...
go build .
```

<a id="validation"></a>
## 验证

```bash
go test ./...
go build .
go run . run -f examples/bell_state.qasm
```

运行环境与已执行结果见 [validation notes](docs/readme/validation.md).

| 项目 | 当前边界 |
| :--- | :--- |
| 状态模型 | 独立 Alpha/Beta 量子位无法忠实表达一般纠缠联合态。 |
| 示例 | Bell 示例验证解析与执行链路，不证明纠缠模拟正确。 |
| 兼容性 | 部分解析构造不代表完整 QASM 标准兼容。 |

<a id="documentation"></a>
## 文档

- [Parser / 解析器](pkg/parser/)
- [Simulator / 模拟器](pkg/simulator/simulator.go)
- [Examples / 示例](examples/)

## 许可证与来源

原创代码和新增演示资源采用 [MIT](LICENSE). 第三方代码、模型、数据和外部资源保留各自许可，详见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
