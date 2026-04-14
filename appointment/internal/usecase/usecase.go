package usecase

import (
	"errors"
	"time"

	"appointment/internal/clients"
	"appointment/internal/model"
	"appointment/internal/repository"
)

type AppointmentUsecase struct {
	repo         repository.AppoiRepo
	doctorClient clients.DoctorClient
}

func NewAppointmentUseCase(repo repository.AppoiRepo, dc clients.DoctorClient) *AppointmentUsecase {
	return &AppointmentUsecase{repo: repo, doctorClient: dc}
}

func (uc *AppointmentUsecase) Create(title, description, doctorID string) (*model.Appointment, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if doctorID == "" {
		return nil, errors.New("doctorID is required")
	}

	if err := uc.doctorClient.DoctorExists(doctorID); err != nil {
		return nil, err
	}

	a := &model.Appointment{
		Title:       title,
		Description: description,
		DoctorID:    doctorID,
		Status:      model.StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (uc *AppointmentUsecase) GetById(id string) (*model.Appointment, error) {
	return uc.repo.GetByID(id)
}

func (uc *AppointmentUsecase) List() ([]*model.Appointment, error) {
	return uc.repo.List()
}

func (uc *AppointmentUsecase) UpdateStatus(id string, newStatus model.Status) (*model.Appointment, error) {
	if !newStatus.IsValid() {
		return nil, errors.New("invalid status")
	}

	a, err := uc.GetById(id)
	if err != nil {
		return nil, err
	}

	if a.Status == model.StatusDone && newStatus == model.StatusNew {
		return nil, errors.New("cannot update done")
	}

	a.Status = newStatus
	a.UpdatedAt = time.Now()

	if err := uc.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}
