// 主要JavaScript文件，用于所有页面共享的功能

document.addEventListener('DOMContentLoaded', function() {
    // 检查是否从示例页面跳转过来
    if (window.location.href.includes('from=examples') && localStorage.getItem('qasmCode')) {
        // 获取存储的代码和版本
        const code = localStorage.getItem('qasmCode');
        const version = localStorage.getItem('qasmVersion');
        
        // 如果页面有编辑器，设置代码和版本
        const editor = document.querySelector('#qasmEditor');
        if (editor) {
            // 如果使用了CodeMirror，需要获取实例
            if (window.CodeMirror && CodeMirror.fromTextArea) {
                const cm = CodeMirror.fromTextArea(editor);
                if (cm) {
                    cm.setValue(code);
                }
            } else {
                editor.value = code;
            }
            
            // 设置版本
            const versionSelect = document.querySelector('#qasmVersion');
            if (versionSelect) {
                versionSelect.value = version;
            }
        }
        
        // 清除存储的代码和版本
        localStorage.removeItem('qasmCode');
        localStorage.removeItem('qasmVersion');
    }
    
    // 添加导航栏激活状态
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.navbar-nav .nav-link');
    
    navLinks.forEach(link => {
        if (link.getAttribute('href') === currentPath) {
            link.classList.add('active');
        } else {
            link.classList.remove('active');
        }
    });
});

// 格式化QASM代码
function formatQasmCode(code) {
    // 简单的格式化，增加缩进和空行
    const lines = code.split('\n');
    let formattedLines = [];
    let indentLevel = 0;
    
    for (let i = 0; i < lines.length; i++) {
        let line = lines[i].trim();
        
        // 跳过空行
        if (!line) {
            formattedLines.push('');
            continue;
        }
        
        // 处理注释
        if (line.startsWith('//')) {
            formattedLines.push(line);
            continue;
        }
        
        // 处理QASM版本和include语句
        if (line.startsWith('OPENQASM') || line.startsWith('include')) {
            formattedLines.push(line);
            continue;
        }
        
        // 处理寄存器声明
        if (line.startsWith('qreg') || line.startsWith('creg')) {
            formattedLines.push(line);
            continue;
        }
        
        // 处理门操作
        const indent = '  '.repeat(indentLevel);
        formattedLines.push(indent + line);
    }
    
    return formattedLines.join('\n');
}

// 显示通知
function showNotification(message, type = 'info') {
    // 创建通知元素
    const notification = document.createElement('div');
    notification.className = `alert alert-${type} alert-dismissible fade show position-fixed`;
    notification.style.top = '20px';
    notification.style.right = '20px';
    notification.style.zIndex = '9999';
    
    notification.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    `;
    
    // 添加到文档
    document.body.appendChild(notification);
    
    // 5秒后自动关闭
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => {
            notification.remove();
        }, 300);
    }, 5000);
}
