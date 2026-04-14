package grpc

import (
	"context"
	"github.com/ddukssu/gogo/doc/internal/model"
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

func (h *Handler) CreateDoctor(ctx context.Context, req *pb.CreateDoctorRequest) (*pb.DoctorResponse, error) {
	if req.Name == "" || req.Specialty == "" {
		return nil, status.Errorf(codes.InvalidArgument, "пустые поля")
	}
	doc := &model.Doctor{Name: req.Name, Specialty: req.Specialty}
	err := h.usecase.CreateDoctor(doc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DoctorResponse{Id: doc.ID, Name: doc.Name, Specialty: doc.Specialty}, nil
}

func (h *Handler) GetDoctor(ctx context.Context, req *pb.GetDoctorRequest) (*pb.DoctorResponse, error) {
	doc, err := h.usecase.GetDoctorByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "врач не найден")
	}
	return &pb.DoctorResponse{Id: doc.ID, Name: doc.Name, Specialty: doc.Specialty}, nil
}

func (h *Handler) ListDoctors(ctx context.Context, req *pb.Empty) (*pb.ListDoctorsResponse, error) {
	docs := h.usecase.GetAllDoctors()
	var resp []*pb.DoctorResponse
	for _, d := range docs {
		resp = append(resp, &pb.DoctorResponse{Id: d.ID, Name: d.Name, Specialty: d.Specialty})
	}
	return &pb.ListDoctorsResponse{Doctors: resp}, nil
}
