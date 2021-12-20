package main

import (
	"context"
	"fmt"
	"grpc-course/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Calculator Client...")

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.Dial("localhost:50052", opts...)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	// close connection
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)
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
