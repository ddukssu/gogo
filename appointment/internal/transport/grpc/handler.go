package grpc

import (
	"context"

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

func (h *Handler) CreateAppointment(_ context.Context, req *pb.CreateAppointmentRequest) (*pb.AppointmentResponse, error) {
	appt, err := h.usecase.Create(req.PatientId, req.Date, req.DoctorId)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unavailable {
			return nil, status.Errorf(codes.Unavailable, st.Message())
		}
		if err.Error() == "doctor does not exist" {
			return nil, status.Errorf(codes.NotFound, "врач не найден")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.AppointmentResponse{
		Id:        appt.ID,
		PatientId: appt.Title,
		DoctorId:  appt.DoctorID,
		Date:      appt.Description,
	}, nil
}
