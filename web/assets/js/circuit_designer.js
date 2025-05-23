// 电路设计器JavaScript

// 全局变量
let editor;
let qubitCount = 3;
let columnCount = 10;
let circuit = []; // 存储电路信息
let selectedGate = null;
let selectedSlot = null;
let draggedGate = null;

// 初始化
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
        lineWrapping: true,
        readOnly: true
    });

    // 初始化电路
    initCircuit();
    
    // 绑定事件
    bindEvents();
    
    // 生成初始QASM代码
    generateQASM();
});

// 初始化电路
function initCircuit() {
    // 初始化电路数据结构
    circuit = [];
    for (let i = 0; i < qubitCount; i++) {
        circuit.push(new Array(columnCount).fill(null));
    }
    
    // 渲染电路
    renderCircuit();
}

// 渲染电路
function renderCircuit() {
    const circuitGrid = document.getElementById('circuitGrid');
    circuitGrid.innerHTML = '';
    
    // 创建量子比特线
    for (let i = 0; i < qubitCount; i++) {
        const qubitLine = document.createElement('div');
        qubitLine.className = 'qubit-line';
        qubitLine.dataset.qubit = i;
        
        // 添加量子比特标签
        const qubitLabel = document.createElement('div');
        qubitLabel.className = 'qubit-label';
        qubitLabel.textContent = `q[${i}]`;
        qubitLine.appendChild(qubitLabel);
        
        // 添加量子线
        const qubitWire = document.createElement('div');
        qubitWire.className = 'qubit-wire';
        qubitLine.appendChild(qubitWire);
        
        // 添加门插槽
        for (let j = 0; j < columnCount; j++) {
            const slotLeft = 60 + j * 50;
            const gateSlot = document.createElement('div');
            gateSlot.className = 'gate-slot';
            gateSlot.style.left = `${slotLeft}px`;
            gateSlot.dataset.qubit = i;
            gateSlot.dataset.column = j;
            
            // 如果该位置有门，则渲染门
            if (circuit[i][j]) {
                const gate = circuit[i][j];
                gateSlot.classList.add('occupied');
                
                const gateElement = document.createElement('div');
                gateElement.className = `gate-on-wire ${getGateClass(gate.name)}`;
                gateElement.textContent = getGateDisplay(gate.name);
                gateElement.dataset.gate = gate.name;
                gateElement.dataset.qubit = i;
                gateElement.dataset.column = j;
                
                if (gate.control) {
                    gateElement.dataset.control = true;
                    gateElement.dataset.controlQubit = gate.controlQubit;
                }
                
                if (gate.param) {
                    gateElement.dataset.param = true;
                    gateElement.dataset.paramValue = gate.paramValue || 'pi/2';
                }
                
                if (gate.secondTarget) {
                    gateElement.dataset.secondTarget = true;
                    gateElement.dataset.secondTargetQubit = gate.secondTargetQubit;
                }
                
                gateSlot.appendChild(gateElement);
                
                // 如果是控制门，添加控制线
                if (gate.control && gate.controlQubit !== undefined) {
                    addControlLine(i, j, gate.controlQubit);
                }
                
                // 如果是双目标门，添加连接线
                if (gate.secondTarget && gate.secondTargetQubit !== undefined) {
                    addConnectionLine(i, j, gate.secondTargetQubit);
                }
            }
            
            qubitLine.appendChild(gateSlot);
            
            // 添加点击事件
            gateSlot.addEventListener('click', handleSlotClick);
        }
        
        circuitGrid.appendChild(qubitLine);
    }
    
    // 更新QASM代码
    generateQASM();
}

// 添加控制线
function addControlLine(targetQubit, column, controlQubit) {
    const circuitGrid = document.getElementById('circuitGrid');
    const slotLeft = 60 + column * 50 + 20; // 中心位置
    
    const controlLine = document.createElement('div');
    controlLine.className = 'control-line';
    controlLine.style.left = `${slotLeft}px`;
    
    // 确定控制线的起点和终点
    const startQubit = Math.min(targetQubit, controlQubit);
    const endQubit = Math.max(targetQubit, controlQubit);
    const startY = startQubit * 50 + 25;
    const height = (endQubit - startQubit) * 50;
    
    controlLine.style.top = `${startY}px`;
    controlLine.style.height = `${height}px`;
    
    circuitGrid.appendChild(controlLine);
}

// 添加连接线（用于SWAP等双目标门）
function addConnectionLine(targetQubit1, column, targetQubit2) {
    const circuitGrid = document.getElementById('circuitGrid');
    const slotLeft = 60 + column * 50 + 20; // 中心位置
    
    const connectionLine = document.createElement('div');
    connectionLine.className = 'control-line';
    connectionLine.style.left = `${slotLeft}px`;
    
    // 确定连接线的起点和终点
    const startQubit = Math.min(targetQubit1, targetQubit2);
    const endQubit = Math.max(targetQubit1, targetQubit2);
    const startY = startQubit * 50 + 25;
    const height = (endQubit - startQubit) * 50;
    
    connectionLine.style.top = `${startY}px`;
    connectionLine.style.height = `${height}px`;
    
    circuitGrid.appendChild(connectionLine);
}

// 获取门的显示文本
function getGateDisplay(gateName) {
    const gateDisplayMap = {
        'h': 'H',
        'x': 'X',
        'y': 'Y',
        'z': 'Z',
        's': 'S',
        't': 'T',
        'cx': 'X',
        'cz': 'Z',
        'swap': '⇄',
        'measure': 'M',
        'rx': 'RX',
        'ry': 'RY',
        'rz': 'RZ',
        'cu1': 'U1'
    };
    
    return gateDisplayMap[gateName] || gateName.toUpperCase();
}

// 获取门的CSS类
function getGateClass(gateName) {
    if (['cx', 'cz', 'cu1'].includes(gateName)) {
        return 'target';
    } else if (gateName === 'measure') {
        return 'measure';
    } else {
        return '';
    }
}

// 绑定事件
function bindEvents() {
    // 添加量子比特按钮
    document.getElementById('addQubitBtn').addEventListener('click', function() {
        if (qubitCount < 100) {
            qubitCount++;
            circuit.push(new Array(columnCount).fill(null));
            renderCircuit();
        }
    });
    
    // 移除量子比特按钮
    document.getElementById('removeQubitBtn').addEventListener('click', function() {
        if (qubitCount > 1) {
            qubitCount--;
            circuit.pop();
            renderCircuit();
        }
    });
    
    // 清空电路按钮
    document.getElementById('clearCircuitBtn').addEventListener('click', function() {
        initCircuit();
    });
    
    // 复制代码按钮
    document.getElementById('copyCodeBtn').addEventListener('click', function() {
        const code = editor.getValue();
        navigator.clipboard.writeText(code)
            .then(() => {
                alert('代码已复制到剪贴板！');
            })
            .catch(err => {
                console.error('复制失败: ', err);
                alert('复制失败，请手动选择代码并复制。');
            });
    });
    
    // 运行代码按钮
    document.getElementById('runCodeBtn').addEventListener('click', function() {
        const code = editor.getValue();
        const version = document.getElementById('qasmVersion').value;
        
        // 将代码存储到localStorage
        localStorage.setItem('qasmCode', code);
        localStorage.setItem('qasmVersion', version);
        
        // 跳转到模拟器页面
        window.location.href = '/ui/simulator?from=designer';
    });
    
    // 量子门拖拽
    const gateItems = document.querySelectorAll('.gate-item');
    gateItems.forEach(gate => {
        gate.addEventListener('dragstart', handleDragStart);
        gate.addEventListener('dragend', handleDragEnd);
    });
    
    // 门插槽拖放
    const circuitGrid = document.getElementById('circuitGrid');
    circuitGrid.addEventListener('dragover', handleDragOver);
    circuitGrid.addEventListener('drop', handleDrop);
    
    // QASM版本选择
    document.getElementById('qasmVersion').addEventListener('change', generateQASM);
}

// 处理门拖拽开始
function handleDragStart(e) {
    draggedGate = {
        name: this.dataset.gate,
        control: this.dataset.control === 'true',
        param: this.dataset.param === 'true',
        twoTarget: this.dataset.twoTarget === 'true'
    };
    
    e.dataTransfer.setData('text/plain', this.dataset.gate);
    e.dataTransfer.effectAllowed = 'copy';
}

// 处理门拖拽结束
function handleDragEnd(e) {
    draggedGate = null;
}

// 处理拖拽经过
function handleDragOver(e) {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'copy';
    
    // 高亮当前悬停的插槽
    const slot = e.target.closest('.gate-slot');
    if (slot && !slot.classList.contains('occupied')) {
        slot.classList.add('active');
    }
    
    // 移除其他插槽的高亮
    document.querySelectorAll('.gate-slot.active').forEach(activeSlot => {
        if (activeSlot !== slot) {
            activeSlot.classList.remove('active');
        }
    });
}

// 处理拖拽放置
function handleDrop(e) {
    e.preventDefault();
    
    // 移除所有插槽的高亮
    document.querySelectorAll('.gate-slot.active').forEach(activeSlot => {
        activeSlot.classList.remove('active');
    });
    
    // 获取放置的插槽
    const slot = e.target.closest('.gate-slot');
    if (!slot || slot.classList.contains('occupied')) {
        return;
    }
    
    const qubit = parseInt(slot.dataset.qubit);
    const column = parseInt(slot.dataset.column);
    
    // 获取拖拽的门
    const gateName = e.dataTransfer.getData('text/plain');
    
    // 如果是控制门，需要选择控制比特
    if (draggedGate.control) {
        selectedGate = {
            name: gateName,
            qubit: qubit,
            column: column,
            control: true
        };
        
        // 提示用户选择控制比特
        alert('请点击另一个量子比特作为控制比特');
        return;
    }
    
    // 如果是双目标门，需要选择第二个目标比特
    if (draggedGate.twoTarget) {
        selectedGate = {
            name: gateName,
            qubit: qubit,
            column: column,
            secondTarget: true
        };
        
        // 提示用户选择第二个目标比特
        alert('请点击另一个量子比特作为第二个目标比特');
        return;
    }
    
    // 如果是带参数的门，需要设置参数
    if (draggedGate.param) {
        // 默认参数
        const paramValue = 'pi/2';
        
        // 添加门到电路
        circuit[qubit][column] = {
            name: gateName,
            param: true,
            paramValue: paramValue
        };
        
        // 显示参数设置面板
        showGateParams(qubit, column, gateName, paramValue);
    } else {
        // 添加普通门到电路
        circuit[qubit][column] = {
            name: gateName
        };
    }
    
    // 重新渲染电路
    renderCircuit();
}

// 处理插槽点击
function handleSlotClick(e) {
    const slot = e.target.closest('.gate-slot');
    if (!slot) return;
    
    const qubit = parseInt(slot.dataset.qubit);
    const column = parseInt(slot.dataset.column);
    
    // 如果当前有选中的门，且需要选择控制比特
    if (selectedGate && selectedGate.control) {
        // 不能选择自己作为控制比特
        if (qubit === selectedGate.qubit) {
            alert('控制比特不能与目标比特相同');
            return;
        }
        
        // 添加控制门到电路
        circuit[selectedGate.qubit][selectedGate.column] = {
            name: selectedGate.name,
            control: true,
            controlQubit: qubit
        };
        
        // 重置选中状态
        selectedGate = null;
        
        // 重新渲染电路
        renderCircuit();
        return;
    }
    
    // 如果当前有选中的门，且需要选择第二个目标比特
    if (selectedGate && selectedGate.secondTarget) {
        // 不能选择自己作为第二个目标比特
        if (qubit === selectedGate.qubit) {
            alert('第二个目标比特不能与第一个目标比特相同');
            return;
        }
        
        // 添加双目标门到电路
        circuit[selectedGate.qubit][selectedGate.column] = {
            name: selectedGate.name,
            secondTarget: true,
            secondTargetQubit: qubit
        };
        
        // 重置选中状态
        selectedGate = null;
        
        // 重新渲染电路
        renderCircuit();
        return;
    }
    
    // 如果点击的是已有门，则删除
    if (circuit[qubit][column]) {
        circuit[qubit][column] = null;
        renderCircuit();
        return;
    }
    
    // 否则，记录当前选中的插槽
    selectedSlot = {
        qubit: qubit,
        column: column
    };
}

// 显示门参数设置面板
function showGateParams(qubit, column, gateName, paramValue) {
    const paramsContainer = document.getElementById('gateParams');
    paramsContainer.innerHTML = '';
    
    const formGroup = document.createElement('div');
    formGroup.className = 'mb-3';
    
    const label = document.createElement('label');
    label.className = 'form-label';
    label.textContent = `${gateName.toUpperCase()} 门参数`;
    formGroup.appendChild(label);
    
    const input = document.createElement('input');
    input.type = 'text';
    input.className = 'form-control';
    input.value = paramValue;
    input.dataset.qubit = qubit;
    input.dataset.column = column;
    formGroup.appendChild(input);
    
    const helpText = document.createElement('div');
    helpText.className = 'form-text';
    helpText.textContent = '可以使用 pi/2, pi/4 等表达式';
    formGroup.appendChild(helpText);
    
    const button = document.createElement('button');
    button.className = 'btn btn-primary mt-2';
    button.textContent = '设置参数';
    button.addEventListener('click', function() {
        const paramValue = input.value;
        const qubit = parseInt(input.dataset.qubit);
        const column = parseInt(input.dataset.column);
        
        // 更新门参数
        if (circuit[qubit][column]) {
            circuit[qubit][column].paramValue = paramValue;
            renderCircuit();
        }
    });
    formGroup.appendChild(button);
    
    paramsContainer.appendChild(formGroup);
}

// 生成QASM代码
function generateQASM() {
    const version = document.getElementById('qasmVersion').value;
    let code = '';
    
    // 添加版本声明
    code += `OPENQASM ${version};\n`;
    
    // 添加包含语句（仅QASM 2.0）
    if (version === '2.0') {
        code += 'include "qelib1.inc";\n';
    }
    
    code += '\n';
    
    // 添加量子寄存器声明
    code += `qreg q[${qubitCount}];\n`;
    
    // 添加经典寄存器声明（用于测量）
    let hasMeasure = false;
    for (let i = 0; i < qubitCount; i++) {
        for (let j = 0; j < columnCount; j++) {
            if (circuit[i][j] && circuit[i][j].name === 'measure') {
                hasMeasure = true;
                break;
            }
        }
        if (hasMeasure) break;
    }
    
    if (hasMeasure) {
        code += `creg c[${qubitCount}];\n`;
    }
    
    code += '\n';
    
    // 添加门操作
    // 按列遍历，确保门操作按顺序
    for (let j = 0; j < columnCount; j++) {
        let hasGateInColumn = false;
        
        // 首先处理单量子比特门
        for (let i = 0; i < qubitCount; i++) {
            const gate = circuit[i][j];
            if (!gate) continue;
            
            // 跳过控制门和双目标门，它们将在后面处理
            if ((gate.control && gate.controlQubit !== undefined) || 
                (gate.secondTarget && gate.secondTargetQubit !== undefined)) {
                continue;
            }
            
            hasGateInColumn = true;
            
            // 生成门操作代码
            if (gate.name === 'measure') {
                code += `measure q[${i}] -> c[${i}];\n`;
            } else if (gate.param) {
                code += `${gate.name}(${gate.paramValue || 'pi/2'}) q[${i}];\n`;
            } else {
                code += `${gate.name} q[${i}];\n`;
            }
        }
        
        // 然后处理双量子比特门
        for (let i = 0; i < qubitCount; i++) {
            const gate = circuit[i][j];
            if (!gate) continue;
            
            // 只处理控制门
            if (gate.control && gate.controlQubit !== undefined) {
                hasGateInColumn = true;
                
                if (gate.param) {
                    code += `${gate.name}(${gate.paramValue || 'pi/2'}) q[${gate.controlQubit}], q[${i}];\n`;
                } else {
                    code += `${gate.name} q[${gate.controlQubit}], q[${i}];\n`;
                }
            }
            
            // 处理双目标门
            if (gate.secondTarget && gate.secondTargetQubit !== undefined) {
                hasGateInColumn = true;
                
                code += `${gate.name} q[${i}], q[${gate.secondTargetQubit}];\n`;
            }
        }
        
        // 如果这一列有门操作，添加空行
        if (hasGateInColumn) {
            code += '\n';
        }
    }
    
    // 更新编辑器内容
    editor.setValue(code);
}
