package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {

	user := User{}

	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	// hasing password procs..
	passwodHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwodHash)
	user.Role = "user"

	// saving data to repository (Struct Input => Struct User)
	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil

}

func (s *service) Login(input LoginInput) (User, error) {

	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)

	// if error are exists
	if err != nil {
		return user, err
	}

	// if email is not exist
	if user.ID == 0 {
		return user, errors.New("user not found with that email")
	}

	// if no-error & email exist:
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, errors.New("password is not match")
	}

	return user, nil

}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {

	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err // false, email is taken
	}

	if user.ID == 0 {
		return true, nil // true, email is availabe to use
	}

	return false, nil // defaut bool is false
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	/*
	 dapatkan user bedasarkan ID
	 update attribute avatar file name
	 simpan perubahan avatar file name
	*/
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return user, err
	}

	return updatedUser, nil
}

func (s *service) GetUserByID(ID int) (User, error) {

	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	// if email is not exist
	if user.ID == 0 {
		return user, errors.New("no user found with that id")
	}

	return user, nil
}
