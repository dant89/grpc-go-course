package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dant89/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Calculate function was invovked with %v", req)
	sum := req.GetCalculator().GetFirstDigit() + req.GetCalculator().GetSecondDigit()
	res := &calculatorpb.CalculatorResponse{
		Result: sum,
	}
	return res, nil
}

func main() {
	fmt.Println("Calculator server listening...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
