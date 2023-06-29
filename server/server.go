package server

import (
	"context"

	"platzi.com/go/grpc/models"
	"platzi.com/go/grpc/repository"
	grpc_studentpb "platzi.com/go/grpc/studentpb"
)

type Server struct {
	repo repository.Repository
	grpc_studentpb.UnimplementedStudentServiceServer
}

func NewStudentServer(repo repository.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetStudent(ctx context.Context, req *grpc_studentpb.GetStudentRequest) (*grpc_studentpb.Student, error) {
	student, err := s.repo.GetStudent(ctx, req.GetId())

	if err != nil {
		return nil, err
	}

	return &grpc_studentpb.Student{
		Id:   student.Id,
		Name: student.Name,
		Age:  student.Age,
	}, nil
}

func (s *Server) SetStudent(ctx context.Context, req *grpc_studentpb.Student) (*grpc_studentpb.SetStudentResponse, error) {
	student := &models.Student{
		Id:   req.GetId(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}

	err := s.repo.SetStudent(ctx, student)

	if err != nil {
		return nil, err
	}

	return &grpc_studentpb.SetStudentResponse{
		Id: student.Id,
	}, nil
}
