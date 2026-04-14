package app

import (
	"appointment/internal/clients"
	"appointment/internal/repository"
	transport "appointment/internal/transport/http"
	"appointment/internal/usecase"

	"github.com/gin-gonic/gin"
)

func Run(addr, doctorServiceURL string) error {
	repo := repository.NewInMemoryRepo()
	doctorClient := clients.NewHTTPDoctorClient(doctorServiceURL)
	uc := usecase.NewAppointmentUseCase(repo, doctorClient)
	handler := transport.NewAppointmentHandler(uc)

	r := gin.Default()
	handler.RegisterRoutes(r)
	return r.Run(addr)
}
