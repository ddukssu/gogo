package main

import (
	"github.com/ddukssu/gogo/appointment/internal/clients"
	"github.com/ddukssu/gogo/appointment/internal/repository"
	transport "github.com/ddukssu/gogo/appointment/internal/transport/grpc"
	"github.com/ddukssu/gogo/appointment/internal/usecase"
	pb "github.com/ddukssu/gogo/appointment/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	docClient, err := clients.NewDoctorClient("localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewInMemoryRepo()
	uc := usecase.NewAppointmentUseCase(repo, docClient)
	handler := transport.NewHandler(uc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAppointmentServiceServer(grpcServer, handler)

	log.Println("Appointment Service: 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
