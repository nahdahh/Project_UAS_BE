package service

import (
	"errors"
	"uas_be/app/model"
	"uas_be/app/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req *model.CreateUserRequest, roleID string) (*model.User, error)
	GetUserByID(id string) (*model.UserWithRole, error)
	GetAllUsers(page, pageSize int) ([]*model.UserWithRole, int, error)
	UpdateUser(id string, username, email, fullName string) error
	DeleteUser(id string) error
}

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// CreateUser membuat user baru (admin only)
func (s *userServiceImpl) CreateUser(req *model.CreateUserRequest, roleID string) (*model.User, error) {
	// Validasi input wajib
	if req.Username == "" || req.Email == "" || req.Password == "" || req.FullName == "" {
		return nil, errors.New("semua field harus diisi")
	}

	// Business rule: username harus unik
	existingUser, _ := s.userRepo.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username sudah terdaftar")
	}

	// Business rule: email harus unik
	existingUser, _ = s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Hash password dengan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal hash password")
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		RoleID:       roleID,
		IsActive:     true,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID mengambil detail user dengan permissions
func (s *userServiceImpl) GetUserByID(id string) (*model.UserWithRole, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Ambil permissions dari role
	permissions, err := s.userRepo.GetUserPermissions(id)
	if err != nil {
		return nil, err
	}

	return &model.UserWithRole{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		IsActive:    user.IsActive,
		Permissions: permissions,
		CreatedAt:   user.CreatedAt.String(),
	}, nil
}

// GetAllUsers mengambil semua user dengan pagination (admin only)
func (s *userServiceImpl) GetAllUsers(page, pageSize int) ([]*model.UserWithRole, int, error) {
	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.userRepo.GetAllUsers(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser mengupdate data user
func (s *userServiceImpl) UpdateUser(id string, username, email, fullName string) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user tidak ditemukan")
	}

	// Hanya update field yang diberikan (optional fields)
	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if fullName != "" {
		user.FullName = fullName
	}

	return s.userRepo.UpdateUser(user)
}

// DeleteUser menghapus user
func (s *userServiceImpl) DeleteUser(id string) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user tidak ditemukan")
	}

	return s.userRepo.DeleteUser(id)
}
