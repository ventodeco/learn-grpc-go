package main

import (
	"context"
	"fmt"
	"grpc-course/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Hello I'm a client")

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.Dial("localhost:50051", opts...)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	// close connection when all logic executed
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Vento",
			LastName:  "Deco",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	var firstName string
	var lastName string

	fmt.Print("Masukan first name : ")
	fmt.Scanln(&firstName)
	fmt.Print("Masukan last name : ")
	fmt.Scanln(&lastName)

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: firstName,
			LastName:  lastName,
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{FirstName: "Vento"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Yosep"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Deco"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Novento"},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
	}

	// we iterate and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{FirstName: "Vento"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Yosep"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Deco"},
		},
		{
			Greeting: &greetpb.Greeting{FirstName: "Novento"},
		},
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client(go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while receiving : %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
