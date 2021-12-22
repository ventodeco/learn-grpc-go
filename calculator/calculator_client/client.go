package main

import (
	"context"
	"fmt"
	"grpc-course/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Calculator Client...")

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.Dial("localhost:50053", opts...)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	// close connection
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)

	// doServerStreaming(c)

	doClientStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Start to sum calculator with unary...")

	var firstNumber int32
	var secondNumber int32

	fmt.Print("Masukan angka pertama : ")
	fmt.Scanln(&firstNumber)
	fmt.Print("Masukan angka kedua : ")
	fmt.Scanln(&secondNumber)
	req := &calculatorpb.SumRequest{
		Calculating: &calculatorpb.Calculating{
			FirstNumber:  firstNumber,
			SecondNumber: secondNumber,
		},
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling calculator RPC: %v", err)
	}

	log.Printf("Response from Sum Calculator: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Prime Decomposition to do a Server Streaming RPC...")

	req := &calculatorpb.NumberRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling prime number decomposition RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from Prime Number Decomposition: %v", msg.GetResult())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC Compute Average...")

	requests := []*calculatorpb.NumberRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
		{
			Number: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("error while calling compute average RPC: %v", err)
	}

	// we iterate and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("Compute Average Response: %v\n", res)
}
