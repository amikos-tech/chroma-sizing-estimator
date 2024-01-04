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
	resp, err := client.Calculate(ctx, &pb.CalculationRequest{NumberOfVectors: 2, DimensionOfVectors: 2})
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}
	if resp.Result != float32(4*2*2)/1024/1024/1024 {
		t.Fatalf("unexpected response: %v", resp.Result)
	}
}
