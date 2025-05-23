// 初始化代码编辑器
let editor;
document.addEventListener('DOMContentLoaded', function() {
    // 初始化CodeMirror编辑器
    editor = CodeMirror.fromTextArea(document.getElementById('exampleCode'), {
        mode: 'javascript',
        theme: 'dracula',
        lineNumbers: true,
        autoCloseBrackets: true,
        matchBrackets: true,
        indentUnit: 4,
        tabSize: 4,
        lineWrapping: true,
        readOnly: true // 示例代码默认为只读
    });

    // 绑定按钮事件
    document.getElementById('runButton').addEventListener('click', runExample);
    document.getElementById('visualizeButton').addEventListener('click', visualizeExample);
    document.getElementById('copyButton').addEventListener('click', copyExampleCode);
    
    // 绑定示例列表点击事件
    const exampleItems = document.querySelectorAll('#examplesList a');
    exampleItems.forEach(item => {
        item.addEventListener('click', function(e) {
            e.preventDefault();
            
            // 移除其他项目的激活状态
            exampleItems.forEach(i => i.classList.remove('active'));
            
            // 激活当前项目
            this.classList.add('active');
            
            // 加载示例
            loadExample(this.getAttribute('data-example'));
        });
    });
});

// 运行示例代码
function runExample() {
    const code = editor.getValue();
    const version = document.getElementById('exampleVersion').textContent.includes('2.0') ? '2.0' : '1.0';
    
    // 打开模拟器页面并传递代码
    localStorage.setItem('qasmCode', code);
    localStorage.setItem('qasmVersion', version);
    
    // 跳转到模拟器页面
    window.location.href = '/ui/simulator?from=examples';
}

// 可视化示例电路
function visualizeExample() {
    const code = editor.getValue();
    const version = document.getElementById('exampleVersion').textContent.includes('2.0') ? '2.0' : '1.0';
    
    // 显示加载状态
    document.getElementById('circuitVisualization').innerHTML = '<div class="text-center py-5"><div class="spinner-border text-primary" role="status"></div><p class="mt-2">正在生成电路图...</p></div>';
    
    // 调用API获取电路数据
    fetch('/api/v1/circuit', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ code, version })
    })
    .then(response => response.json())
    .then(data => {
        renderCircuit(data);
    })
    .catch(error => {
        document.getElementById('circuitVisualization').innerHTML = `<div class="alert alert-danger">${error.message}</div>`;
    });
}

// 复制示例代码
function copyExampleCode() {
    const code = editor.getValue();
    
    // 使用Clipboard API复制代码
    navigator.clipboard.writeText(code)
        .then(() => {
            // 显示成功消息
            const copyButton = document.getElementById('copyButton');
            const originalText = copyButton.textContent;
            copyButton.textContent = '复制成功！';
            copyButton.classList.remove('btn-secondary');
            copyButton.classList.add('btn-success');
            
            // 2秒后恢复按钮状态
            setTimeout(() => {
                copyButton.textContent = originalText;
                copyButton.classList.remove('btn-success');
                copyButton.classList.add('btn-secondary');
            }, 2000);
        })
        .catch(err => {
            console.error('复制失败: ', err);
            alert('复制失败，请手动选择代码并复制。');
        });
}

// 加载示例代码
function loadExample(exampleName) {
    // 根据选择的示例加载代码
    let exampleCode = '';
    let exampleTitle = '';
    let exampleDescription = '';
    let qasmVersion = 'QASM 2.0';
    
    switch (exampleName) {
        case 'bell':
            exampleTitle = 'Bell态';
            exampleDescription = 'Bell态是量子纠缠的最简单例子，也称为EPR对。这个示例创建了两个量子比特的最大纠缠态。';
            exampleCode = `// Bell状态示例 (QASM 2.0)
OPENQASM 2.0;
include "qelib1.inc";

qreg q[2];
creg c[2];

h q[0];
cx q[0], q[1];

measure q -> c;`;
            break;
            
        case 'teleportation':
            exampleTitle = '量子隐形传态';
            exampleDescription = '量子隐形传态是一种利用量子纠缠和经典通信将量子态从一个地方传送到另一个地方的协议，而无需实际传输量子比特本身。';
            exampleCode = `// 量子隐形传态示例 (QASM 2.0)
OPENQASM 2.0;
include "qelib1.inc";

// 声明量子寄存器和经典寄存器
qreg q[3];  // q[0]是Alice的源量子比特，q[1]和q[2]是纠缠对
creg c[2];  // 用于存储测量结果

// 初始化源量子比特为任意状态（这里使用RY门创建一个非平凡状态）
ry(pi/4) q[0];

// 创建Bell对（q[1]和q[2]之间的纠缠）
h q[1];
cx q[1], q[2];

// 开始隐形传态协议
// 步骤1：Alice将她的量子比特与她的纠缠量子比特进行纠缠
cx q[0], q[1];
h q[0];

// 步骤2：Alice测量她的两个量子比特
measure q[0] -> c[0];
measure q[1] -> c[1];

// 步骤3：根据测量结果，Bob对他的量子比特应用相应的门操作
// 注意：在真实的量子计算机上，这需要经典条件操作，这里我们使用两个可能的操作
// 如果c[1]==1，应用X门
// 如果c[0]==1，应用Z门

// 最终，q[2]将包含原始q[0]的量子态`;
            break;
            
        case 'grover':
            exampleTitle = 'Grover算法';
            exampleDescription = 'Grover算法是一种量子搜索算法，可以在O(√N)时间内在无序数据库中找到目标元素，相比经典算法的O(N)时间复杂度有显著提升。';
            exampleCode = `// Grover算法示例 (QASM 2.0)
OPENQASM 2.0;
include "qelib1.inc";

qreg q[3];
creg c[3];

// 初始化为叠加态
h q[0];
h q[1];
h q[2];

// Oracle (标记状态|101⟩)
x q[0];
x q[2];
h q[2];
ccx q[0], q[1], q[2];
h q[2];
x q[0];
x q[2];

// 振幅放大
h q[0];
h q[1];
h q[2];
x q[0];
x q[1];
x q[2];
h q[2];
ccx q[0], q[1], q[2];
h q[2];
x q[0];
x q[1];
x q[2];
h q[0];
h q[1];
h q[2];

// 测量
measure q -> c;`;
            break;
            
        case 'deutsch':
            exampleTitle = 'Deutsch算法';
            exampleDescription = 'Deutsch算法是最早的量子算法之一，它可以用一次查询确定一个黑盒函数是常数函数还是平衡函数，而经典算法需要两次查询。';
            exampleCode = `// Deutsch算法示例 (QASM 2.0)
OPENQASM 2.0;
include "qelib1.inc";

qreg q[2];
creg c[1];

// 初始化
x q[1];
h q[0];
h q[1];

// 应用Oracle (这里实现的是平衡函数f(x) = x)
cx q[0], q[1];

// 测量
h q[0];
measure q[0] -> c[0];`;
            break;
            
        case 'qft':
            exampleTitle = '量子傅里叶变换';
            exampleDescription = '量子傅里叶变换(QFT)是经典傅里叶变换的量子版本，是许多量子算法的核心组件，包括Shor算法。';
            exampleCode = `// 量子傅里叶变换示例 (QASM 2.0)
OPENQASM 2.0;
include "qelib1.inc";

qreg q[3];
creg c[3];

// 初始化为非平凡状态
x q[0];

// 3比特QFT
// 第一个量子比特
h q[0];
cu1(pi/2) q[0], q[1];
cu1(pi/4) q[0], q[2];

// 第二个量子比特
h q[1];
cu1(pi/2) q[1], q[2];

// 第三个量子比特
h q[2];

// 交换量子比特顺序（可选）
swap q[0], q[2];

// 测量
measure q -> c;`;
            break;
            
        case 'shor':
            exampleTitle = 'Shor算法';
            exampleDescription = 'Shor算法是一种量子算法，用于对大整数进行质因数分解，它比已知的任何经典算法都要快得多。这里展示的是简化版本。';
            exampleCode = `// Shor算法简化示例 (QASM 2.0)
OPENQASM 2.0;
include "qelib1.inc";

// 这是Shor算法的简化版本，用于演示目的
// 实际的Shor算法需要更多的量子比特和更复杂的电路

qreg q[7];  // 4个工作量子比特 + 3个辅助量子比特
creg c[4];  // 用于存储测量结果

// 初始化工作寄存器为叠加态
h q[0];
h q[1];
h q[2];
h q[3];

// 模拟模幂运算（这里是简化版本）
// 在实际的Shor算法中，这部分会根据要分解的数字和选择的a值来构建

// 应用QFT到工作寄存器
// QFT开始
h q[0];
cu1(pi/2) q[0], q[1];
cu1(pi/4) q[0], q[2];
cu1(pi/8) q[0], q[3];

h q[1];
cu1(pi/2) q[1], q[2];
cu1(pi/4) q[1], q[3];

h q[2];
cu1(pi/2) q[2], q[3];

h q[3];
// QFT结束

// 测量工作寄存器
measure q[0] -> c[0];
measure q[1] -> c[1];
measure q[2] -> c[2];
measure q[3] -> c[3];`;
            break;
    }
    
    // 更新页面内容
    document.getElementById('exampleTitle').textContent = exampleTitle;
    document.getElementById('exampleDescription').innerHTML = `<p>${exampleDescription}</p>`;
    document.getElementById('exampleVersion').textContent = qasmVersion;
    
    // 更新编辑器内容
    editor.setValue(exampleCode);
    
    // 清空电路可视化
    document.getElementById('circuitVisualization').innerHTML = '<div class="text-center py-5 text-muted"><p>点击"可视化电路"按钮查看量子电路图</p></div>';
}

// 渲染电路图
function renderCircuit(circuitData) {
    // 这里应该根据circuitData渲染电路图
    // 为了演示，我们将生成一个简单的电路图
    
    const qubits = circuitData.qubits || 3;
    const gates = circuitData.gates || [
        { name: 'h', targets: [0], controls: [], position: 0 },
        { name: 'cx', targets: [1], controls: [0], position: 1 },
        { name: 'measure', targets: [0, 1], controls: [], position: 2 }
    ];
    
    let circuitHtml = '';
    
    // 为每个量子比特创建一条线
    for (let i = 0; i < qubits; i++) {
        circuitHtml += `<div class="d-flex align-items-center mb-3">`;
        circuitHtml += `<div class="me-3" style="width: 50px;"><strong>q[${i}]</strong></div>`;
        
        // 创建量子线
        circuitHtml += `<div class="flex-grow-1 position-relative">`;
        circuitHtml += `<div class="quantum-wire"></div>`;
        
        // 在线上放置量子门
        gates.forEach(gate => {
            const isTarget = gate.targets.includes(i);
            const isControl = gate.controls.includes(i);
            
            if (isTarget || isControl) {
                const position = gate.position * 50 + 10; // 根据门的位置计算偏移量
                const gateClass = isControl ? 'control' : (gate.name === 'measure' ? 'measure' : 'target');
                const gateText = isControl ? '•' : (gate.name === 'measure' ? 'M' : gate.name);
                
                circuitHtml += `<div class="quantum-gate ${gateClass}" style="position: absolute; top: -15px; left: ${position}px;">${gateText}</div>`;
            }
        });
        
        circuitHtml += `</div>`;
        circuitHtml += `</div>`;
    }
    
    document.getElementById('circuitVisualization').innerHTML = circuitHtml;
}
