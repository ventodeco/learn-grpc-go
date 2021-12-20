package main

import (
	"context"
	"fmt"
	"grpc-course/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	_numberOne := req.GetCalculating().GetFirstNumber()
	_numberTwo := req.GetCalculating().GetSecondNumber()
	result := _numberOne + _numberTwo

	fmt.Printf("Hasil dari penjumlahan %d + %d = %d\n", _numberOne, _numberTwo, result)

	res := &calculatorpb.SumResponse{
		SumResult: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Calculator Sum Server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
