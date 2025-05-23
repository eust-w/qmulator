// Quantum Experience (QE) Standard Header
// File: qelib1.inc

// 单量子比特门
gate u3(theta,phi,lambda) q { U(theta,phi,lambda) q; }
gate u2(phi,lambda) q { U(pi/2,phi,lambda) q; }
gate u1(lambda) q { U(0,0,lambda) q; }
gate rx(theta) q { u3(theta,-pi/2,pi/2) q; }
gate ry(theta) q { u3(theta,0,0) q; }
gate rz(theta) q { u1(theta) q; }
gate id q { u1(0) q; }
gate x q { u3(pi,0,pi) q; }
gate y q { u3(pi,pi/2,pi/2) q; }
gate z q { u1(pi) q; }
gate h q { u2(0,pi) q; }
gate s q { u1(pi/2) q; }
gate sdg q { u1(-pi/2) q; }
gate t q { u1(pi/4) q; }
gate tdg q { u1(-pi/4) q; }

// 控制门
gate cx c,t { CX c,t; }
gate cz c,t { h t; cx c,t; h t; }
gate cy c,t { sdg t; cx c,t; s t; }
gate ch c,t { h t; sdg c; cx c,t; h t; s c; }
gate ccx a,b,c { h c; cx b,c; tdg c; cx a,c; t c; cx b,c; tdg c; cx a,c; t b; t c; h c; cx a,b; t a; tdg b; cx a,b; }
gate crz(lambda) c,t { u1(lambda/2) t; cx c,t; u1(-lambda/2) t; cx c,t; }
gate cu1(lambda) c,t { u1(lambda/2) c; cx c,t; u1(-lambda/2) t; cx c,t; u1(lambda/2) t; }
gate cu3(theta,phi,lambda) c,t { u1((lambda-phi)/2) t; cx c,t; u3(-theta/2,0,-(phi+lambda)/2) t; cx c,t; u3(theta/2,phi,0) t; }
