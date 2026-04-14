package http

import (
	"errors"
	"net/http"

	"doc/internal/repository"
	"doc/internal/usecase"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	uc *usecase.DocUsecase
}

func NewDoctorHandler(uc *usecase.DocUsecase) *DoctorHandler {
	return &DoctorHandler{uc: uc}
}

func (h *DoctorHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/doctors", h.Create)
	r.GET("/doctors", h.List)
	r.GET("/doctors/:id", h.GetByID)
}

type createDoctorRequest struct {
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization"`
	Email          string `json:"email"`
}

func (h *DoctorHandler) Create(c *gin.Context) {
	var req createDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor, err := h.uc.Create(req.FullName, req.Specialization, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, doctor)
}

func (h *DoctorHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	doctor, err := h.uc.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) List(c *gin.Context) {
	doctors, err := h.uc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}
