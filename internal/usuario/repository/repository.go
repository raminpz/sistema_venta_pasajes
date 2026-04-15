package repository

import (
	"context"
	"errors"
	"sistema_venta_pasajes/internal/usuario/domain"

	"gorm.io/gorm"
)

type UsuarioRepository interface {
	Create(ctx context.Context, usuario *domain.Usuario) error
	Update(ctx context.Context, id int, usuario *domain.Usuario) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*domain.Usuario, error)
	GetByEmail(ctx context.Context, email string) (*domain.Usuario, error)
	GetByDNI(ctx context.Context, dni string) (*domain.Usuario, error)
	List(ctx context.Context, filtro map[string]interface{}, offset, limit int) ([]*domain.Usuario, int, error)
}

type usuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) UsuarioRepository {
	return &usuarioRepository{db: db}
}

func (r *usuarioRepository) Create(ctx context.Context, usuario *domain.Usuario) error {
	return r.db.WithContext(ctx).Create(usuario).Error
}

func (r *usuarioRepository) Update(ctx context.Context, id int, usuario *domain.Usuario) error {
	return r.db.WithContext(ctx).Model(&domain.Usuario{}).Where("id_usuario = ?", id).Updates(usuario).Error
}

func (r *usuarioRepository) Delete(ctx context.Context, id int) error {
	res := r.db.WithContext(ctx).Delete(&domain.Usuario{}, "id_usuario = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *usuarioRepository) GetByID(ctx context.Context, id int) (*domain.Usuario, error) {
	var usuario domain.Usuario
	err := r.db.WithContext(ctx).First(&usuario, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &usuario, err
}

func (r *usuarioRepository) GetByEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	var usuario domain.Usuario
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&usuario).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &usuario, err
}

func (r *usuarioRepository) GetByDNI(ctx context.Context, dni string) (*domain.Usuario, error) {
	var usuario domain.Usuario
	err := r.db.WithContext(ctx).Where("dni = ?", dni).First(&usuario).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &usuario, err
}

func (r *usuarioRepository) List(ctx context.Context, filtro map[string]interface{}, offset, limit int) ([]*domain.Usuario, int, error) {
	var usuarios []*domain.Usuario
	var total int64
	tx := r.db.WithContext(ctx).Model(&domain.Usuario{})
	for k, v := range filtro {
		tx = tx.Where(k+" = ?", v)
	}
	err := tx.Offset(offset).Limit(limit).Find(&usuarios).Error
	if err != nil {
		return nil, 0, err
	}
	err = tx.Count(&total).Error
	return usuarios, int(total), err
}


