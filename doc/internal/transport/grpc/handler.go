package grpc

import (
	"context"

	"github.com/ddukssu/gogo/doc/internal/usecase"
	pb "github.com/ddukssu/gogo/doc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedDoctorServiceServer
	usecase *usecase.DocUsecase
}

func NewHandler(uc *usecase.DocUsecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) CreateDoctor(_ context.Context, req *pb.CreateDoctorRequest) (*pb.DoctorResponse, error) {
	if req.Name == "" || req.Specialty == "" {
		return nil, status.Errorf(codes.InvalidArgument, "пустые поля")
	}

	doc, err := h.usecase.Create(req.Name, req.Specialty, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.DoctorResponse{
		Id:        doc.ID,
		Name:      doc.FullName,
		Specialty: doc.Specialization,
	}, nil
}

func (h *Handler) GetDoctor(_ context.Context, req *pb.GetDoctorRequest) (*pb.DoctorResponse, error) {
	doc, err := h.usecase.GetByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "врач не найден")
	}

	return &pb.DoctorResponse{
		Id:        doc.ID,
		Name:      doc.FullName,
		Specialty: doc.Specialization,
	}, nil
}

func (h *Handler) ListDoctors(_ context.Context, _ *pb.Empty) (*pb.ListDoctorsResponse, error) {
	docs, err := h.usecase.List()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var resp []*pb.DoctorResponse
	for _, d := range docs {
		resp = append(resp, &pb.DoctorResponse{
			Id:        d.ID,
			Name:      d.FullName,
			Specialty: d.Specialization,
		})
	}

	return &pb.ListDoctorsResponse{Doctors: resp}, nil
}
