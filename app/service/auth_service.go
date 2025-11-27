package service

import (
	"errors"
	"time"
	"uas_be/app/model"
	"uas_be/app/repository"
	"uas_be/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	// Login melakukan autentikasi dan mengembalikan token
	Login(username, password string) (*model.LoginResponse, error)
	// GetUserProfile mengambil profile user dengan permissions
	GetUserProfile(userID string) (*model.UserProfile, error)
}

type authServiceImpl struct {
	userRepo       repository.UserRepository
	permissionRepo repository.PermissionRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
	permissionRepo repository.PermissionRepository,
) AuthService {
	return &authServiceImpl{
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
	}
}

// Login melakukan autentikasi user dengan username dan password, mengembalikan token JWT
func (s *authServiceImpl) Login(username, password string) (*model.LoginResponse, error) {
	// Validasi input wajib
	if username == "" || password == "" {
		return nil, errors.New("username dan password tidak boleh kosong")
	}

	// Cari user berdasarkan username
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Business rule: user harus aktif
	if !user.IsActive {
		return nil, errors.New("user tidak aktif")
	}

	// Verify password hash dengan bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("password salah")
	}

	// Ambil permissions dari role user
	permissions, err := s.permissionRepo.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.RoleID, permissions, 24*time.Hour)
	if err != nil {
		return nil, errors.New("gagal generate token: " + err.Error())
	}

	return &model.LoginResponse{
		Token: token,
		User: model.UserWithRole{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			RoleName:    user.RoleID, // TODO: ambil nama role dari database
			IsActive:    user.IsActive,
			Permissions: permissions,
			CreatedAt:   user.CreatedAt.String(),
		},
	}, nil
}

// GetUserProfile mengambil profile user dengan permissions berdasarkan user ID
func (s *authServiceImpl) GetUserProfile(userID string) (*model.UserProfile, error) {
	// Cari user berdasarkan ID
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Ambil permissions dari role user
	permissions, err := s.permissionRepo.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return nil, err
	}

	return &model.UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		RoleName:    user.RoleID, // TODO: ambil nama role dari database
		Permissions: permissions,
		CreatedAt:   user.CreatedAt.String(),
	}, nil
}
