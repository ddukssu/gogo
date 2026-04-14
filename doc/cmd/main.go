package main

import (
	"github.com/ddukssu/gogo/doc/internal/repository"
	transport "github.com/ddukssu/gogo/doc/internal/transport/grpc"
	"github.com/ddukssu/gogo/doc/internal/usecase"
	pb "github.com/ddukssu/gogo/doc/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	repo := repository.NewDocRepo()
	uc := usecase.NewDocUsecase(repo)
	handler := transport.NewHandler(uc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDoctorServiceServer(grpcServer, handler)

	log.Println("Doctor Service: 50051")
	grpcServer.Serve(lis)
}
