package Usecases

import (
	"errors"
	"task_manager_clean/Domain"
	"task_manager_clean/Infrastructure"
	"task_manager_clean/Repositories"
	"time"
)

type UserUsecase interface {
	Register(username, password string) (*Domain.User, error)
	Login(username, password string) (string, error) // returns jwt token
	Promote(username string) error
}

type userUsecase struct {
	userRepo Repositories.UserRepository
	pw       Infrastructure.PasswordService
	jwt      Infrastructure.JWTService
}

func NewUserUsecase(r Repositories.UserRepository, pw Infrastructure.PasswordService, jwt Infrastructure.JWTService) UserUsecase {
	return &userUsecase{userRepo: r, pw: pw, jwt: jwt}
}

func (u *userUsecase) Register(username, password string) (*Domain.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password required")
	}

	if _, err := u.userRepo.FindByUsername(username); err == nil {
		return nil, errors.New("username already exists")
	}
	hash, err := u.pw.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &Domain.User{
		Username:     username,
		PasswordHash: hash,
		Role:         "user",
		CreatedAt:    time.Now().UTC(),
	}
	// if no users => admin
	if c, _ := u.userRepo.Count(); c == 0 {
		user.Role = "admin"
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (u *userUsecase) Login(username, password string) (string, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := u.pw.Compare(user.PasswordHash, password); err != nil {
		return "", errors.New("invalid credentials")
	}
	// generate token
	token, err := u.jwt.Generate(user.ID.Hex(), user.Username, user.Role, 24*time.Hour)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUsecase) Promote(username string) error {
	return u.userRepo.Promote(username)
}
