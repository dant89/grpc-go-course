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

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("GreetManyTimes function was invovked with %v", req)

	/*
		k = 2
		N = 210
		while N > 1:
			if N % k == 0:   // if k evenly divides into N
				print k      // this is a factor
				N = N / k    // divide N by k so that we have the rest of the number left.
			else:
				k = k + 1
	*/

	divisor := int64(2)
	primeNumber := req.GetPrimeNumberDecomposition().GetPrimeNumber()

	for primeNumber > 1 {
		if primeNumber%divisor == 0 {
			primeNumber = primeNumber / divisor
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: divisor,
			}
			stream.Send(res)
		} else {
			divisor++
		}
	}

	return nil
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
