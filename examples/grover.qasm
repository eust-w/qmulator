// Grover算法示例 (QASM 2.0)
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
measure q -> c;
