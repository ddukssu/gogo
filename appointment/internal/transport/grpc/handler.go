package grpc

import (
	"context"
	"github.com/ddukssu/gogo/appointment/internal/model"
	"github.com/ddukssu/gogo/appointment/internal/usecase"
	pb "github.com/ddukssu/gogo/appointment/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedAppointmentServiceServer
	usecase *usecase.AppointmentUsecase
}

func NewHandler(uc *usecase.AppointmentUsecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) CreateAppointment(ctx context.Context, req *pb.CreateAppointmentRequest) (*pb.AppointmentResponse, error) {
	app := &model.Appointment{PatientID: req.PatientId, DoctorID: req.DoctorId, Date: req.Date}
	err := h.usecase.CreateAppointment(app)
	if err != nil {
		// Трансляция ошибки из usecase в gRPC код
		if err.Error() == "doctor does not exist" {
			return nil, status.Errorf(codes.NotFound, "врач не найден")
		}
		if err.Error() == "doctor service недоступен" {
			return nil, status.Errorf(codes.Unavailable, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.AppointmentResponse{Id: app.ID, PatientId: app.PatientID, DoctorId: app.DoctorID, Date: app.Date}, nil
}
