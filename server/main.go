package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/amikos-tech/chroma-sizing-estimator/pb"
)

type server struct {
	pb.CalculatorServiceServer
}

func (s *server) Calculate(
	_ context.Context, in *pb.CalculationRequest,
) (*pb.CalculationResponse, error) {

	if in.NumberOfVectors <= 0 {
		return nil, status.Errorf(
			codes.InvalidArgument, "At least one vector is required",
		)
	}
	if in.VectorDimensions <= 0 {
		return nil, status.Errorf(
			codes.InvalidArgument, "Dimension of vectors must be positive",
		)
	}
	var systemMemoryOverhead float32 = 0.2
	if in.SystemMemoryOverhead != nil {
		systemMemoryOverhead = *in.SystemMemoryOverhead
	}
	var binaryIndexEstimate = float32(4*in.NumberOfVectors*in.VectorDimensions) / 1024 / 1024 / 1024
	var memorySizeEstimate = binaryIndexEstimate * (1 + systemMemoryOverhead)
	return &pb.CalculationResponse{
		MemorySizeEstimate: memorySizeEstimate,
		EstimateUnit:       pb.EstimateUnit_GB,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("failed to create listener:", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalln("failed to serve:", err)
	}
}
