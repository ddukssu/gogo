package app

import (
	"github.com/ddukssu/gogo/appointment/internal/clients"
	"github.com/ddukssu/gogo/appointment/internal/repository"
	transport "github.com/ddukssu/gogo/appointment/internal/transport/http"
	"github.com/ddukssu/gogo/appointment/internal/usecase"

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
