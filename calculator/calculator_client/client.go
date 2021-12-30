package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dant89/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Calcuate client running...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)

	doBiDirectionalStreaming(c)

	doUErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Unary RPC...")
	req := &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{
			FirstDigit:  3,
			SecondDigit: 10,
		},
	}

	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculate RPC: %v", err)
	}

	log.Printf("Response from Calculate: %v\n", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		PrimeNumberDecomposition: &calculatorpb.PrimeNumberDecomposition{
			PrimeNumber: 120,
		},
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// Reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeNumberDecomposition: %v\n", msg.GetResult())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}

	numbers := []int64{3, 5, 9, 54, 23}

	for _, number := range numbers {
		stream.Send(&calculatorpb.ComputeAverageRequest{
			ComputeAverage: &calculatorpb.ComputeAverage{
				Digit: number,
			},
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while recieving response from LongGreet: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}

func doBiDirectionalStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Bi Directional Streaming RPC...")

	// Create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*calculatorpb.FindMaximumRequest{}
	digits := []int64{4, 7, 2, 19, 4, 6, 32}
	for _, digit := range digits {
		requests = append(requests, &calculatorpb.FindMaximumRequest{
			FindMaximum: &calculatorpb.FindMaximum{
				Digit: digit,
			},
		})
	}

	waitc := make(chan struct{})
	// Send messages to the client in a go routine
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending mesage: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 + time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Recieve messages from the client in a go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Recieved: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// Block until channel is empty
	<-waitc
}

func doUErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do SquareRoot Unary RPC...")

	doErrorCall(c, int32(9))
	doErrorCall(c, int32(-2))
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: n,
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// Actual error form gRPC
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Error while calling SquareRoot RPC: %v", err)
		}
	}

	fmt.Printf("Result of square root of %v:%v\n", n, res.GetNumberRoot())
}
