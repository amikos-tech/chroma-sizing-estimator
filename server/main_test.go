package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"

	"github.com/amikos-tech/chroma-sizing-estimator/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCalculate(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)
	client := pb.NewCalculatorServiceClient(conn)

	t.Run("Test Calculate With Default System Buffer Memory", func(t *testing.T) {
		resp, err := client.Calculate(ctx, &pb.CalculationRequest{NumberOfVectors: 2, VectorDimensions: 2})
		if err != nil {
			t.Fatalf("Calculate failed: %v", err)
		}
		if resp.GetMemorySizeEstimate() != (float32(4*2*2)/1024/1024/1024)*(1+0.2) {
			t.Fatalf("unexpected response: %v", resp.MemorySizeEstimate)
		}
		if resp.GetEstimateUnit() != pb.EstimateUnit_GB {
			t.Fatalf("unexpected response: %v", resp.EstimateUnit)
		}
	})

	t.Run("Test Calculate With Custom System Buffer Memory", func(t *testing.T) {
		systemMemoryOverhead := float32(0.3)
		resp, err := client.Calculate(ctx, &pb.CalculationRequest{NumberOfVectors: 2, VectorDimensions: 2, SystemMemoryOverhead: &systemMemoryOverhead})
		if err != nil {
			t.Fatalf("Calculate failed: %v", err)
		}
		if resp.GetMemorySizeEstimate() != (float32(4*2*2)/1024/1024/1024)*(1+systemMemoryOverhead) {
			t.Fatalf("unexpected response: %v", resp.MemorySizeEstimate)
		}
		if resp.GetEstimateUnit() != pb.EstimateUnit_GB {
			t.Fatalf("unexpected response: %v", resp.EstimateUnit)
		}
	})
}
