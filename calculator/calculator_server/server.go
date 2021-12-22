package main

import (
	"context"
	"fmt"
	"grpc-course/calculator/calculatorpb"
	"io"
	"log"
	"net"
	"time"

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

func (*server) PrimeNumberDecomposition(req *calculatorpb.NumberRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	number := req.GetNumber()

	fmt.Printf("PrimeNumberManyTimes function was invoked with %v\n", number)

	var k, N int64
	k = 2
	N = number

	for N > 1 {
		if N%k == 0 {
			res := &calculatorpb.NumberResponse{
				Result: float32(k),
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			N = N / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with streaming request...\n")

	var result, jumlah int64

	result = 0
	jumlah = 0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			var _res float32
			_res = float32(result) / float32(jumlah)
			// we have finished reading the client streaming
			return stream.SendAndClose(&calculatorpb.NumberResponse{
				Result: _res,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client streaming: %v", err)
		}

		jumlah++
		result += req.GetNumber()
	}
}

func main() {
	fmt.Println("Calculator Sum Server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50053")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
