package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"platzi.com/go/grpc/database"
	"platzi.com/go/grpc/server"
	grpc_studentpb "platzi.com/go/grpc/studentpb"
)

func main() {
	list, err := net.Listen("tcp", ":5060")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository("", "", "")

	if err != nil {
		log.Fatal(err)
	}

	server := server.NewStudentServer(repo)
	s := grpc.NewServer()
	grpc_studentpb.RegisterStudentServiceServer(s, server)

	// Provee metadata a los clientes
	reflection.Register(s)
	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
