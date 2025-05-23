// 初始化代码编辑器
let editor;
document.addEventListener('DOMContentLoaded', function() {
    // 初始化CodeMirror编辑器
    editor = CodeMirror.fromTextArea(document.getElementById('qasmEditor'), {
        mode: 'javascript',
        theme: 'dracula',
        lineNumbers: true,
        autoCloseBrackets: true,
        matchBrackets: true,
        indentUnit: 4,
        tabSize: 4,
        lineWrapping: true
    });

    // 绑定按钮事件
    document.getElementById('runButton').addEventListener('click', runSimulation);
    document.getElementById('stepButton').addEventListener('click', runStepByStep);
    document.getElementById('resetButton').addEventListener('click', resetSimulation);
    document.getElementById('visualizeButton').addEventListener('click', visualizeCircuit);
    document.getElementById('formatButton').addEventListener('click', formatQASMCode);
    document.getElementById('loadExampleButton').addEventListener('click', loadExample);
});

// 运行模拟
function runSimulation() {
    const code = editor.getValue();
    const version = document.getElementById('qasmVersion').value;
    
    // 显示加载状态
    document.getElementById('simulationResults').innerHTML = '<div class="text-center py-5"><div class="spinner-border text-primary" role="status"></div><p class="mt-2">正在运行模拟...</p></div>';
    
    // 调用API运行QASM代码
    fetch('/api/v1/run', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ code, version })
    })
    .then(response => response.json())
    .then(data => {
        displayResults(data);
    })
    .catch(error => {
        document.getElementById('simulationResults').innerHTML = `<div class="alert alert-danger">${error.message}</div>`;
    });
}

// 步骤运行
function runStepByStep() {
    const code = editor.getValue();
    const version = document.getElementById('qasmVersion').value;
    
    // 显示加载状态
    document.getElementById('stateEvolution').innerHTML = '<div class="text-center py-5"><div class="spinner-border text-primary" role="status"></div><p class="mt-2">正在准备步骤运行...</p></div>';
    
    // 这里应该调用API的步骤运行功能，但目前后端可能还没有实现
    // 为了演示，我们将模拟步骤运行的效果
    simulateStepByStep(code);
}

// 模拟步骤运行效果
function simulateStepByStep(code) {
    // 解析代码，提取操作步骤
    const lines = code.split('\n');
    const steps = [];
    
    for (const line of lines) {
        const trimmedLine = line.trim();
        if (trimmedLine && !trimmedLine.startsWith('//') && !trimmedLine.startsWith('OPENQASM') && !trimmedLine.startsWith('include')) {
            if (trimmedLine.includes('qreg') || trimmedLine.includes('creg') || 
                trimmedLine.includes('h ') || trimmedLine.includes('cx ') || 
                trimmedLine.includes('measure')) {
                steps.push(trimmedLine);
            }
        }
    }
    
    // 显示步骤
    let stepsHtml = '<div class="list-group">';
    steps.forEach((step, index) => {
        stepsHtml += `<div class="list-group-item" id="step-${index}">${step}</div>`;
    });
    stepsHtml += '</div>';
    
    document.getElementById('stateEvolution').innerHTML = stepsHtml;
    
    // 模拟步骤执行
    let currentStep = 0;
    const stateInterval = setInterval(() => {
        // 移除之前的高亮
        if (currentStep > 0) {
            document.getElementById(`step-${currentStep - 1}`).classList.remove('step-highlight');
        }
        
        // 添加当前步骤高亮
        if (currentStep < steps.length) {
            document.getElementById(`step-${currentStep}`).classList.add('step-highlight');
            currentStep++;
        } else {
            // 所有步骤执行完毕
            clearInterval(stateInterval);
            // 运行完整模拟
            runSimulation();
        }
    }, 1000);
}

// 重置模拟
function resetSimulation() {
    // 清空结果显示
    document.getElementById('simulationResults').innerHTML = '<div class="text-center py-5 text-muted"><p>点击"运行"按钮查看模拟结果</p></div>';
    document.getElementById('stateEvolution').innerHTML = '<div class="text-center py-5 text-muted"><p>点击"步骤运行"查看量子态演化</p></div>';
    document.getElementById('circuitVisualization').innerHTML = '<div class="text-center py-5 text-muted"><p>点击"可视化"按钮查看量子电路图</p></div>';
}

// 可视化电路
function visualizeCircuit() {
    const code = editor.getValue();
    const version = document.getElementById('qasmVersion').value;
    
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

// 显示模拟结果
function displayResults(data) {
    if (!data.success) {
        // 显示错误信息
        let errorHtml = '<div class="alert alert-danger">';
        if (data.errors && data.errors.length > 0) {
            errorHtml += `<p>模拟出错：</p><ul>`;
            data.errors.forEach(error => {
                errorHtml += `<li>${error}</li>`;
            });
            errorHtml += `</ul>`;
        } else {
            errorHtml += `<p>模拟出错，请检查QASM代码。</p>`;
        }
        errorHtml += '</div>';
        
        document.getElementById('simulationResults').innerHTML = errorHtml;
        return;
    }
    
    // 显示成功结果
    let resultsHtml = '<div class="card-text">';
    
    // 显示寄存器值
    if (data.registers && Object.keys(data.registers).length > 0) {
        resultsHtml += '<h6>寄存器状态：</h6>';
        
        for (const [name, bits] of Object.entries(data.registers)) {
            const binaryValue = bits.map(bit => bit ? '1' : '0').join('');
            const decimalValue = data.decimal[name] || 0;
            
            resultsHtml += `<div class="mb-2">`;
            resultsHtml += `<div><strong>${name}:</strong> ${binaryValue} (十进制: ${decimalValue})</div>`;
            
            // 添加可视化条形图
            if (decimalValue > 0) {
                const percentage = Math.min(100, decimalValue * 5); // 简单缩放以适应显示
                resultsHtml += `<div class="result-bar" style="width: ${percentage}%;">${decimalValue}</div>`;
            } else {
                resultsHtml += `<div class="result-bar" style="width: 5%;">0</div>`;
            }
            
            resultsHtml += `</div>`;
        }
    } else {
        resultsHtml += '<p>没有可用的寄存器结果。</p>';
    }
    
    resultsHtml += '</div>';
    document.getElementById('simulationResults').innerHTML = resultsHtml;
}

// 格式化QASM代码
function formatQASMCode() {
    const code = editor.getValue();
    
    // 显示格式化中的状态
    const originalButton = document.getElementById('formatButton').innerHTML;
    document.getElementById('formatButton').innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span> 格式化中...';
    document.getElementById('formatButton').disabled = true;
    
    // 调用API格式化QASM代码
    fetch('/api/v1/format', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ code })
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            // 如果有错误，显示错误消息
            alert('格式化错误: ' + data.error);
        } else {
            // 更新编辑器内容
            editor.setValue(data.formattedCode);
        }
        
        // 恢复按钮状态
        document.getElementById('formatButton').innerHTML = originalButton;
        document.getElementById('formatButton').disabled = false;
    })
    .catch(error => {
        // 处理错误
        alert('格式化请求失败: ' + error.message);
        document.getElementById('formatButton').innerHTML = originalButton;
        document.getElementById('formatButton').disabled = false;
    });
}

// 加载示例代码
function loadExample() {
    const exampleSelect = document.getElementById('exampleSelect');
    const selectedExample = exampleSelect.value;
    
    if (!selectedExample) {
        return;
    }
    
    // 根据选择的示例加载代码
    let exampleCode = '';
    let qasmVersion = '2.0';
    
    switch (selectedExample) {
        case 'bell':
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
    }
    
    // 设置编辑器内容和版本
    editor.setValue(exampleCode);
    document.getElementById('qasmVersion').value = qasmVersion;
    
    // 重置显示
    resetSimulation();
}
