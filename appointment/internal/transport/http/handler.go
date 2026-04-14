package http

import (
	"errors"
	"log"
	"net/http"

	"appointment/internal/clients"
	"appointment/internal/model"
	"appointment/internal/repository"
	"appointment/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	uc *usecase.AppointmentUsecase
}

func NewAppointmentHandler(uc *usecase.AppointmentUsecase) *AppointmentHandler {
	return &AppointmentHandler{uc: uc}
}

func (h *AppointmentHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/appointments", h.Create)
	r.GET("/appointments/:id", h.GetByID)
	r.GET("/appointments", h.List)
	r.PATCH("/appointments/:id/status", h.UpdateStatus)
}

type createAppointmentRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DoctorID    string `json:"doctor_id"`
}

type UpdateStatusRequest struct {
	Status model.Status `json:"status"`
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	var req createAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appt, err := h.uc.Create(req.Title, req.Description, req.DoctorID)
	if err != nil {
		if errors.Is(err, clients.ErrDoctorServiceUnavailable) {
			log.Printf("[ERROR] Doctor service unavailable: %v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Doctor service is currently unavailable. Please try again later."})
			return
		}
		if errors.Is(err, clients.ErrDoctorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "referenced doctor does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appt)
}

func (h *AppointmentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	appt, err := h.uc.GetById(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandler) List(c *gin.Context) {
	appts, err := h.uc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appts)
}

func (h *AppointmentHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appt, err := h.uc.UpdateStatus(id, req.Status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appt)
}
