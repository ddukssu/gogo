package app

import (
	"doc/internal/repository"
	transport "doc/internal/transport/http"
	"doc/internal/usecase"

	"github.com/gin-gonic/gin"
)

func Run(addr string) error {
	repo := repository.NewInMemoryRepo()
	uc := usecase.NewDocUsecase(repo)
	handler := transport.NewDoctorHandler(uc)

	r := gin.Default()
	handler.RegisterRoutes(r)
	return r.Run("0.0.0.0:8080")
}
