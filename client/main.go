package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"platzi.com/go/grpc/testpb"
)

func main() {
	// Conexión insegura pues no cuenta con SSL
	cc, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := testpb.NewTestServiceClient(cc)

	// DoUnary(c)
	// DoClientStreaming(c)
	DoBidirectionalStreaming(c)

}

func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "tst1",
	}

	res, err := c.GetTest(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling GetTest: %v", err)
	}

	log.Printf("response from server: %v", res)
}

func DoClientStreaming(c testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "q8t1",
			Answer:   "Azul",
			Question: "Color asociado a golang",
			TestId:   "tst1",
		},
		{
			Id:       "q9t1",
			Answer:   "Google",
			Question: "Empresa que desarrolló en lenguaje golang",
			TestId:   "tst1",
		},
		{
			Id:       "q10t1",
			Answer:   "Backend",
			Question: "Especialidad de golang",
			TestId:   "tst1",
		},
	}

	stream, err := c.SetQuestions(context.Background())

	if err != nil {
		log.Fatalf("error while calling SetQuestions: %v", err)
	}

	for _, question := range questions {
		log.Println("sending question: ", question.Id)
		stream.Send(question)
		time.Sleep(2 * time.Second)
	}

	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while reciving response: %v", err)
	}
	log.Printf("response from server: %v", msg)
}

func DoServerStreaming(c testpb.TestServiceClient) {
	req := &testpb.GetStudentsPerTestRequest{
		TestId: "tst1",
	}

	stream, err := c.GetStudentsPerTest(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling GetStudentsPerTest: %v", err)
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("response from server: %v", msg)
	}
}

func DoBidirectionalStreaming(c testpb.TestServiceClient) {
	answer := testpb.TakeTestRequest{
		Answer: "42",
	}

	numberOfQuestions := 4

	waitChannel := make(chan struct{})

	stream, err := c.TakeTest(context.Background())
	if err != nil {
		log.Fatalf("error while calling TakeTest: %v", err)
	}

	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			stream.Send(&answer)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
				break
			}

			log.Printf("response from server: %v", res)
		}
		close(waitChannel)
	}()

	<-waitChannel
}
