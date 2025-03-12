// internal/users/repository.go
package users

import (
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	GetAll() ([]User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) Create(user *User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepositoryImpl) GetByID(id string) (*User, error) {
	var user User
	err := r.DB.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserRepositoryImpl) GetByUsername(username string) (*User, error) {
	var user User
	err := r.DB.First(&user, "username = ?", username).Error
	return &user, err
}

func (r *UserRepositoryImpl) GetByEmail(email string) (*User, error) {
	var user User
	err := r.DB.First(&user, "email = ?", email).Error
	return &user, err
}

func (r *UserRepositoryImpl) Update(user *User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepositoryImpl) Delete(id string) error {
	return r.DB.Delete(&User{}, "id = ?", id).Error
}

func (r *UserRepositoryImpl) GetAll() ([]User, error) {
	var users []User
	err := r.DB.Find(&users).Error
	return users, err
}
