package repository

import "sistema_venta_pasajes/internal/programacion/domain"

type MockProgramacionRepository struct {
	CreateFunc  func(programacion *domain.Programacion) error
	UpdateFunc  func(programacion *domain.Programacion) error
	DeleteFunc  func(id int64) error
	GetByIDFunc func(id int64) (*domain.Programacion, error)
	ListFunc    func(offset, limit int) ([]domain.Programacion, int, error)
}

func (m *MockProgramacionRepository) Create(programacion *domain.Programacion) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(programacion)
	}
	return nil
}

func (m *MockProgramacionRepository) Update(programacion *domain.Programacion) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(programacion)
	}
	return nil
}

func (m *MockProgramacionRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockProgramacionRepository) GetByID(id int64) (*domain.Programacion, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockProgramacionRepository) List(offset, limit int) ([]domain.Programacion, int, error) {
	if m.ListFunc != nil {
		return m.ListFunc(offset, limit)
	}
	return nil, 0, nil
}
