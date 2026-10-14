# Release validation

Date: 2026-10-14. Checks are executed in an isolated working checkout. Screenshots use local synthetic/demo data. A source implementation, a passing automated test and external-service acceptance are separate claims.

## Commands

```bash
go test ./...
go build .
go run . run -f examples/bell_state.qasm
```

PASS: Go tests/build and bundled Bell parser/execution. Inline gate bodies and nested parameter expressions now parse; CLI returns nonzero on parse errors. Independent-qubit results are not evidence of correct entanglement simulation. Hardware, external credentials, signing and production deployment are separate setup/acceptance steps.
