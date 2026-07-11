package auth

import (
	"context"
	"fmt"
	"server/internal/application/storage"
	"server/internal/domain"
	"server/internal/security"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo       UserRepository
	storageRepo    storage.StorageRepository
	passwordHasher security.BcryptHasher
	jwtGenerate    security.JwtGenerate
}

func NewAuthService(userRepo UserRepository, storageRepo storage.StorageRepository, passwordHasher security.BcryptHasher, jwtGenerate security.JwtGenerate) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		storageRepo:    storageRepo,
		passwordHasher: passwordHasher,
		jwtGenerate:    jwtGenerate,
	}
}

func (as *AuthService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// check username
	existed, err := as.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	if existed != nil {
		return nil, fmt.Errorf("User aldready exist")
	}
	// hash password
	hashedpassword, err := as.passwordHasher.HashPassword(req.Password)
	// NewUser Domain
	user, err := domain.NewUser(req.Username, hashedpassword, req.Fullname)
	if err != nil {
		return nil, err
	}

	// CreateUser
	err = as.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("Create user fail")
	}

	storage := domain.NewStorage(user.ID)

	userStorage, err := as.storageRepo.Create(ctx, *storage)

	// generate refresh token + access token
	ttl, accessToken, err := as.jwtGenerate.GenerateAccessToken(user.ID, user.Role, req.Username)
	_, refreshToken, err := as.jwtGenerate.GenerateRefreshToken(user.ID, user.Role, req.Username)
	return DomainToDTo(accessToken, refreshToken, int(ttl), user.Fullname, user.Role, userStorage.CurrenSize, userStorage.LimitSize), nil
}

func (as *AuthService) SignIn(ctx context.Context, req SignInRequest) (*RegisterResponse, error) {
	user, err := as.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Wrong Username")
		}
	}
	checkPass := as.passwordHasher.HashedToPassword(req.Password, user.HashedPassword)
	if checkPass != nil {
		return nil, fmt.Errorf("Wrong Password")
	}
	ttl, accessToken, err := as.jwtGenerate.GenerateAccessToken(user.ID, user.Role, req.Username)
	_, refreshToken, err := as.jwtGenerate.GenerateRefreshToken(user.ID, user.Role, req.Username)

	userStorage, err := as.storageRepo.GetStorage(ctx, user.ID)
	return DomainToDTo(accessToken, refreshToken, int(ttl), user.Fullname, user.Role, userStorage.CurrenSize, userStorage.LimitSize), nil
}
