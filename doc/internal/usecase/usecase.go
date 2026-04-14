package usecase

import (
	"errors"

	"github.com/ddukssu/gogo/doc/internal/model"
	"github.com/ddukssu/gogo/doc/internal/repository"
)

var ErrValidation = errors.New("validation error")

type DocUsecase struct {
	repo repository.DoctorRepository
}

func NewDocUsecase(repo repository.DoctorRepository) *DocUsecase {
	return &DocUsecase{repo: repo}
}

func (uc *DocUsecase) Create(fullName, specialization, email string) (*model.Doctor, error) {
	if fullName == "" {
		return nil, errors.New("fullName is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}

	d := &model.Doctor{
		FullName:       fullName,
		Specialization: specialization,
		Email:          email,
	}

	if err := uc.repo.Create(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (uc *DocUsecase) GetByID(id string) (*model.Doctor, error) {
	return uc.repo.GetByID(id)
}

func (uc *DocUsecase) List() ([]*model.Doctor, error) {
	return uc.repo.List()
}
