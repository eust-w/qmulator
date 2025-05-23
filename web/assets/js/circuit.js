// 初始化代码编辑器
let editor;

// 全局变量，用于动画控制
let animationSpeed = 1000; // 动画速度（毫秒）
let animationPaused = false;
let currentStep = 0;
let circuitGates = [];
let circuitQubits = 0;
let animationInterval;
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
    document.getElementById('visualizeButton').addEventListener('click', generateCircuit);
    document.getElementById('exportButton').addEventListener('click', exportCircuitImage);
    document.getElementById('resetButton').addEventListener('click', resetCircuit);
    document.getElementById('loadExampleButton').addEventListener('click', loadExample);
    
    // 绑定动画控制按钮事件
    document.getElementById('playButton').addEventListener('click', playAnimation);
    document.getElementById('pauseButton').addEventListener('click', pauseAnimation);
    document.getElementById('stepForwardButton').addEventListener('click', stepForward);
    document.getElementById('stepBackButton').addEventListener('click', stepBack);
    document.getElementById('speedUpButton').addEventListener('click', speedUp);
    document.getElementById('slowDownButton').addEventListener('click', slowDown);
    document.getElementById('resetAnimationButton').addEventListener('click', resetAnimation);
});

// 生成电路图
function generateCircuit() {
    const code = editor.getValue();
    const version = document.getElementById('qasmVersion').value;
    
    // 显示加载状态
    document.getElementById('circuitVisualization').innerHTML = '<div class="text-center py-5"><div class="spinner-border text-primary" role="status"></div><p class="mt-2">正在生成电路图...</p></div>';
    
    // 重置动画状态
    resetAnimation();
    
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
        // 存储电路数据供动画使用
        circuitGates = data.gates || [];
        circuitQubits = data.qubits || 3;
        
        // 初始化电路图
        renderCircuit(data, true);
        generateCircuitDescription(data, code);
        
        // 显示动画控制面板
        document.getElementById('animationControls').classList.remove('d-none');
    })
    .catch(error => {
        document.getElementById('circuitVisualization').innerHTML = `<div class="alert alert-danger">${error.message}</div>`;
    });
}

// 渲染电路图
function renderCircuit(circuitData, initialRender = false) {
    // 获取电路数据
    const qubits = circuitData.qubits || 3;
    const gates = circuitData.gates || [
        { name: 'h', targets: [0], controls: [], position: 0 },
        { name: 'cx', targets: [1], controls: [0], position: 1 },
        { name: 'measure', targets: [0, 1], controls: [], position: 2 }
    ];
    
    // 创建电路容器
    let circuitHtml = '';
    
    // 添加动画状态指示器
    circuitHtml += `<div class="mb-3 animation-status">
        <div class="d-flex justify-content-between align-items-center">
            <span class="badge bg-primary">当前步骤: <span id="currentStepDisplay">${currentStep}</span> / ${gates.length}</span>
            <span class="badge bg-secondary">动画速度: <span id="speedDisplay">${animationSpeed}ms</span></span>
        </div>
    </div>`;
    
    // 创建一个表格来表示电路
    circuitHtml += `<div class="circuit-table">`;
    
    // 为每个量子比特创建一条线
    for (let i = 0; i < qubits; i++) {
        circuitHtml += `<div class="d-flex align-items-center mb-3 qubit-row" data-qubit="${i}">`;
        circuitHtml += `<div class="me-3" style="width: 50px;"><strong>q[${i}]</strong></div>`;
        
        // 创建量子线
        circuitHtml += `<div class="flex-grow-1 position-relative">`;
        circuitHtml += `<div class="quantum-wire"></div>`;
        
        // 用于存储当前行的门
        const rowGates = [];
        
        // 在线上放置量子门
        gates.forEach((gate, gateIndex) => {
            const isTarget = gate.targets.includes(i);
            const isControl = gate.controls.includes(i);
            
            if (isTarget || isControl) {
                const position = gate.position * 50 + 10; // 根据门的位置计算偏移量
                
                // 根据当前动画步骤决定是否显示门
                const isVisible = initialRender || gateIndex <= currentStep;
                const activeClass = gateIndex === currentStep ? 'active-gate' : '';
                const visibilityStyle = isVisible ? '' : 'visibility: hidden;';
                
                // 为控制点和目标点使用不同的样式
                if (isControl) {
                    // 控制点使用黑点
                    circuitHtml += `<div class="quantum-gate control ${activeClass}" 
                                   data-gate-index="${gateIndex}"
                                   data-gate-position="${gate.position}"
                                   data-gate-name="${gate.name}"
                                   style="position: absolute; top: -15px; left: ${position}px; ${visibilityStyle}">
                                   •
                               </div>`;
                } else if (gate.name === 'measure') {
                    // 测量门使用特殊样式
                    circuitHtml += `<div class="quantum-gate measure ${activeClass}" 
                                   data-gate-index="${gateIndex}"
                                   data-gate-position="${gate.position}"
                                   data-gate-name="${gate.name}"
                                   style="position: absolute; top: -15px; left: ${position}px; ${visibilityStyle}">
                                   M
                               </div>`;
                } else {
                    // 目标门显示门名称和参数
                    let gateText = gate.name;
                    
                    // 如果有参数，显示第一个参数
                    if (gate.params && gate.params.length > 0) {
                        gateText = `${gate.name}(${gate.params[0]})`;
                    }
                    
                    circuitHtml += `<div class="quantum-gate target ${activeClass}" 
                                   data-gate-index="${gateIndex}"
                                   data-gate-position="${gate.position}"
                                   data-gate-name="${gate.name}"
                                   style="position: absolute; top: -15px; left: ${position}px; ${visibilityStyle}">
                                   ${gateText}
                               </div>`;
                }
                
                // 存储门信息用于后续绘制连线
                rowGates.push({
                    position: position,
                    isControl: isControl,
                    isTarget: isTarget,
                    gateIndex: gateIndex,
                    gateName: gate.name,
                    qubit: i
                });
            }
        });
        
        circuitHtml += `</div>`;
        circuitHtml += `</div>`;
    }
    
    circuitHtml += `</div>`;
    
    // 添加当前门操作的说明
    if (gates.length > 0 && currentStep < gates.length) {
        const currentGate = gates[currentStep];
        let paramsDisplay = '';
        
        // 格式化参数显示
        if (currentGate.params && currentGate.params.length > 0) {
            paramsDisplay = `<p><strong>参数:</strong> ${currentGate.params.join(', ')}</p>`;
        }
        
        circuitHtml += `<div class="mt-3 p-2 border rounded bg-light">
            <h5>当前操作: ${currentGate.name.toUpperCase()}</h5>
            <p><strong>目标量子比特:</strong> ${currentGate.targets.map(t => `q[${t}]`).join(', ')}</p>
            ${currentGate.controls.length > 0 ? `<p><strong>控制量子比特:</strong> ${currentGate.controls.map(c => `q[${c}]`).join(', ')}</p>` : ''}
            ${paramsDisplay}
        </div>`;
    }
    
    document.getElementById('circuitVisualization').innerHTML = circuitHtml;
    
    // 在DOM更新后添加控制线
    if (initialRender || currentStep >= 0) {
        drawControlLines(gates, currentStep, initialRender);
    }
}

// 绘制控制点和目标点之间的连线
function drawControlLines(gates, currentStep, initialRender) {
    // 对于每个门，如果它有控制点和目标点，绘制它们之间的连线
    gates.forEach((gate, gateIndex) => {
        // 只处理当前步骤之前的门
        if (!initialRender && gateIndex > currentStep) {
            return;
        }
        
        // 只处理有控制点的门
        if (gate.controls.length === 0 || gate.targets.length === 0) {
            return;
        }
        
        // 找到控制点和目标点的DOM元素
        const position = gate.position;
        const controlElements = [];
        const targetElements = [];
        
        // 查找所有控制点元素
        gate.controls.forEach(controlQubit => {
            const controlSelector = `.qubit-row[data-qubit="${controlQubit}"] .quantum-gate[data-gate-position="${position}"]`;
            const controlElement = document.querySelector(controlSelector);
            if (controlElement) {
                controlElements.push(controlElement);
            }
        });
        
        // 查找所有目标点元素
        gate.targets.forEach(targetQubit => {
            const targetSelector = `.qubit-row[data-qubit="${targetQubit}"] .quantum-gate[data-gate-position="${position}"]`;
            const targetElement = document.querySelector(targetSelector);
            if (targetElement) {
                targetElements.push(targetElement);
            }
        });
        
        // 如果找到了控制点和目标点，绘制连线
        if (controlElements.length > 0 && targetElements.length > 0) {
            // 找出最小和最大的量子比特索引
            const allQubits = [...gate.controls, ...gate.targets];
            const minQubit = Math.min(...allQubits);
            const maxQubit = Math.max(...allQubits);
            
            // 创建一个连线元素
            const lineElement = document.createElement('div');
            lineElement.className = 'control-line';
            
            // 设置连线的位置和大小
            const topQubitRow = document.querySelector(`.qubit-row[data-qubit="${minQubit}"]`);
            const bottomQubitRow = document.querySelector(`.qubit-row[data-qubit="${maxQubit}"]`);
            
            if (topQubitRow && bottomQubitRow) {
                const topRect = topQubitRow.getBoundingClientRect();
                const bottomRect = bottomQubitRow.getBoundingClientRect();
                const circuitRect = document.querySelector('.circuit-table').getBoundingClientRect();
                
                // 计算连线的位置和高度
                const left = gate.position * 50 + 10 + 12; // 门的位置 + 偏移量 + 门的一半宽度
                const top = topRect.top - circuitRect.top + 15; // 顶部量子比特的位置 + 偏移量
                const height = bottomRect.top - topRect.top; // 连线的高度
                
                // 设置连线的样式
                lineElement.style.position = 'absolute';
                lineElement.style.left = `${left}px`;
                lineElement.style.top = `${top}px`;
                lineElement.style.width = '2px';
                lineElement.style.height = `${height}px`;
                lineElement.style.backgroundColor = '#007bff';
                lineElement.style.zIndex = '1';
                
                // 将连线添加到电路图中
                document.querySelector('.circuit-table').appendChild(lineElement);
            }
        }
    });
}

// 生成电路说明
function generateCircuitDescription(circuitData, code) {
    // 解析代码，提取注释和关键信息
    const lines = code.split('\n');
    let description = '';
    
    // 提取注释作为描述
    for (const line of lines) {
        if (line.trim().startsWith('//')) {
            description += line.trim().substring(2).trim() + ' ';
        }
    }
    
    // 添加电路统计信息
    const qubits = circuitData.qubits || 3;
    const gates = circuitData.gates || [];
    
    let descriptionHtml = '<div class="card-text">';
    
    if (description) {
        descriptionHtml += `<p><strong>电路描述：</strong> ${description}</p>`;
    }
    
    descriptionHtml += `<p><strong>电路统计：</strong></p>`;
    descriptionHtml += `<ul>`;
    descriptionHtml += `<li>量子比特数量：${qubits}</li>`;
    
    // 统计不同类型的门
    const gateCounts = {};
    gates.forEach(gate => {
        gateCounts[gate.name] = (gateCounts[gate.name] || 0) + 1;
    });
    
    descriptionHtml += `<li>门操作总数：${gates.length}</li>`;
    
    for (const [gateName, count] of Object.entries(gateCounts)) {
        descriptionHtml += `<li>${gateName} 门数量：${count}</li>`;
    }
    
    descriptionHtml += `</ul>`;
    descriptionHtml += '</div>';
    
    document.getElementById('circuitDescription').innerHTML = descriptionHtml;
}

// 导出电路图为图片
function exportCircuitImage() {
    // 这里应该实现将电路图导出为图片的功能
    // 为了演示，我们将显示一个提示消息
    alert('导出功能正在开发中，敬请期待！');
}

// 重置电路
function resetCircuit() {
    // 清空电路图和描述
    document.getElementById('circuitVisualization').innerHTML = '<div class="text-center py-5 text-muted"><p>点击"生成电路图"按钮查看量子电路图</p></div>';
    document.getElementById('circuitDescription').innerHTML = '<div class="text-center py-5 text-muted"><p>生成电路图后将显示电路说明</p></div>';
    
    // 隐藏动画控制面板
    document.getElementById('animationControls').classList.add('d-none');
    
    // 重置动画状态
    resetAnimation();
}

// 重置动画状态
function resetAnimation() {
    // 清除现有的动画定时器
    if (animationInterval) {
        clearInterval(animationInterval);
        animationInterval = null;
    }
    
    // 重置动画状态
    currentStep = 0;
    animationPaused = true;
}

// 播放动画
function playAnimation() {
    // 如果已经到达最后一步，重置到开始
    if (currentStep >= circuitGates.length) {
        currentStep = 0;
    }
    
    // 如果已经有定时器在运行，先清除
    if (animationInterval) {
        clearInterval(animationInterval);
    }
    
    // 设置动画状态
    animationPaused = false;
    
    // 创建新的定时器
    animationInterval = setInterval(() => {
        // 如果已经到达最后一步，停止动画
        if (currentStep >= circuitGates.length) {
            pauseAnimation();
            return;
        }
        
        // 更新当前步骤
        renderCircuit({ qubits: circuitQubits, gates: circuitGates });
        currentStep++;
        
        // 更新步骤显示
        document.getElementById('currentStepDisplay').textContent = currentStep;
    }, animationSpeed);
}

// 暂停动画
function pauseAnimation() {
    if (animationInterval) {
        clearInterval(animationInterval);
        animationInterval = null;
    }
    animationPaused = true;
}

// 向前一步
function stepForward() {
    // 如果已经到达最后一步，不做任何操作
    if (currentStep >= circuitGates.length) {
        return;
    }
    
    // 暂停当前动画
    pauseAnimation();
    
    // 增加当前步骤
    currentStep++;
    
    // 重新渲染电路图
    renderCircuit({ qubits: circuitQubits, gates: circuitGates });
}

// 向后一步
function stepBack() {
    // 如果已经到达第一步，不做任何操作
    if (currentStep <= 0) {
        return;
    }
    
    // 暂停当前动画
    pauseAnimation();
    
    // 减少当前步骤
    currentStep--;
    
    // 重新渲染电路图
    renderCircuit({ qubits: circuitQubits, gates: circuitGates });
}

// 加快动画速度
function speedUp() {
    // 最小速度为100毫秒
    animationSpeed = Math.max(100, animationSpeed / 1.5);
    
    // 更新速度显示
    document.getElementById('speedDisplay').textContent = Math.round(animationSpeed) + 'ms';
    
    // 如果动画正在运行，重新启动以应用新速度
    if (!animationPaused && animationInterval) {
        pauseAnimation();
        playAnimation();
    }
}

// 减慢动画速度
function slowDown() {
    // 最大速度为5000毫秒
    animationSpeed = Math.min(5000, animationSpeed * 1.5);
    
    // 更新速度显示
    document.getElementById('speedDisplay').textContent = Math.round(animationSpeed) + 'ms';
    
    // 如果动画正在运行，重新启动以应用新速度
    if (!animationPaused && animationInterval) {
        pauseAnimation();
        playAnimation();
    }
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
    resetCircuit();
}
