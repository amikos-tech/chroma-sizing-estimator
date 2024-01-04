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

Testing with System Buffer Memory Size (20% of total memory):

```bash
grpcurl -plaintext -d '{"number_of_vectors": 1000000, "vector_dimensions": 128}' localhost:8080 calculator.CalculatorService/Calculate
```

Should return:

```json
{
  "memorySizeEstimate": 0.47683716,
  "estimateUnit": "GB"
}
```

Testing with System Buffer Memory Size (30% of total memory):

```bash
grpcurl -plaintext -d '{"number_of_vectors": 1000000, "vector_dimensions": 128, "system_memory_overhead": 0.3}' localhost:8080 calculator.CalculatorService/Calculate
```

Should return:

```json
{
  "memorySizeEstimate": 0.6198883,
  "estimateUnit": "GB"
}
```

