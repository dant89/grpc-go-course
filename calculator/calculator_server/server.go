package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/dant89/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	fmt.Printf("PrimeNumberDecomposition function was invovked with %v", req)

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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage function was invovked")

	var digits []int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// End of client stream
			var total int64
			for _, digit := range digits {
				total += digit
			}
			result := float64(total) / float64(len(digits))

			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		digits = append(digits, req.GetComputeAverage().GetDigit())
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a bi directional streaming request")
	maximum := int64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil

		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		digit := req.FindMaximum.GetDigit()
		if digit > maximum {
			maximum = digit
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Result: maximum,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot function was invovked with %v", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
