package repository

import (
	"github.com/ddukssu/gogo/appointment/internal/model"

	"errors"

	"github.com/google/uuid"

	"sync"
)

var ErrNotFound = errors.New("appointment not found")

type AppoiRepo interface {
	Create(a *model.Appointment) error
	GetByID(id string) (*model.Appointment, error)
	List() ([]*model.Appointment, error)
	Update(a *model.Appointment) error
}

type inMemoryRepo struct {
	mu           sync.RWMutex
	appointments map[string]*model.Appointment
}

func NewInMemoryRepo() AppoiRepo {
	return &inMemoryRepo{
		appointments: make(map[string]*model.Appointment),
	}
}

func (r *inMemoryRepo) Create(a *model.Appointment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a.ID = uuid.NewString()
	r.appointments[a.ID] = a
	return nil
}

func (r *inMemoryRepo) GetByID(id string) (*model.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.appointments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (r *inMemoryRepo) List() ([]*model.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*model.Appointment, 0, len(r.appointments))
	for _, a := range r.appointments {
		res = append(res, a)
	}
	return res, nil
}

func (r *inMemoryRepo) Update(a *model.Appointment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.appointments[a.ID]; !ok {
		return ErrNotFound
	}
	r.appointments[a.ID] = a
	return nil
}
