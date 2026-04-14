package app

import (
	"github.com/ddukssu/gogo/doc/internal/repository"
	transport "github.com/ddukssu/gogo/doc/internal/transport/http"
	"github.com/ddukssu/gogo/doc/internal/usecase"

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
