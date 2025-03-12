package users

import (
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository UserRepository
}
type UserService interface {
	CreateUser(username, email, password, role string) (*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
	Login(email, password string) (*User, error)
	GetAllUsers() ([]User, error)
}

func NewUserService(userRepository UserRepository) UserService {
	return &UserServiceImpl{UserRepository: userRepository}
}

func (s *UserServiceImpl) GetAllUsers() ([]User, error) {
	return s.UserRepository.GetAll()
}

func (s *UserServiceImpl) CreateUser(username, email, password, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	err = s.UserRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetUserByID(id string) (*User, error) {
	return s.UserRepository.GetByID(id)
}

func (s *UserServiceImpl) GetUserByUsername(username string) (*User, error) {
	return s.UserRepository.GetByUsername(username)
}

func (s *UserServiceImpl) UpdateUser(user *User) error {
	return s.UserRepository.Update(user)
}

func (s *UserServiceImpl) DeleteUser(id string) error {
	return s.UserRepository.Delete(id)
}

func (s *UserServiceImpl) Login(email, password string) (*User, error) {

	user, err := s.UserRepository.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
