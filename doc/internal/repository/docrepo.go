package repository

import (
	"errors"
	"sync"

	"github.com/ddukssu/gogo/doc/internal/model"

	"github.com/google/uuid"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrEmailTaken = errors.New("email already taken")
)

type DoctorRepository interface {
	Create(d *model.Doctor) error
	GetByID(id string) (*model.Doctor, error)
	List() ([]*model.Doctor, error)
}

type inMemoryRepo struct {
	mu      sync.RWMutex
	doctors map[string]*model.Doctor
}

func NewInMemoryRepo() DoctorRepository {
	return &inMemoryRepo{
		doctors: make(map[string]*model.Doctor),
	}
}

func (r *inMemoryRepo) Create(d *model.Doctor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.doctors {
		if existing.Email == d.Email {
			return ErrEmailTaken
		}
	}

	d.ID = uuid.NewString()
	r.doctors[d.ID] = d
	return nil
}

func (r *inMemoryRepo) GetByID(id string) (*model.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.doctors[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (r *inMemoryRepo) List() ([]*model.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.Doctor, 0, len(r.doctors))
	for _, d := range r.doctors {
		result = append(result, d)
	}
	return result, nil
}
