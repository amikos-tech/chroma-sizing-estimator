# Chroma Sizing Calculator

This is a tiny experiment in calculating ChromaDB memory requirement based on input number of vectors and dimensionality
of vectors over a gRPC interface.

## Usage

Standalone server:

```bash
make build
```

Run server:

```bash
make run
```

Build docker image:

```bash
make docker-build
```

Run docker image:

```bash
make docker-run
```

### Testing

With grpcurl:

```bash
grpcurl -plaintext -d '{"number_of_vectors": 1000000, "dimension_of_vectors": 128}' localhost:8080 calculator.CalculatorService/Calculate
```

Should return:

```json
{
  "result": 0.47683716
}
```
