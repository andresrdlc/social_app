package application

import (
	"errors"
	"regexp"

	"github.com/unmsmfisi-socialapplication/social_app/internal/register/domain"
)

var (
	ErrEmailInUse = errors.New("EMAIL_IN_USE")
	ErrFormat = errors.New("INVALID_PASSWORD")
	ErrPhone = errors.New("INVALID_PHONE")
	nombreErr=errors.New("PLANTILLA_ERROR")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
)

type UserRepository interface {
	GetUserByEmail(email string) (*domain.User, error)
	InsertUser(newUser *domain.User) (*domain.User, error)
}

type RegistrationUseCase struct {
	repo UserRepository
}

func NewRegistrationUseCase(r UserRepository) *RegistrationUseCase {
	return &RegistrationUseCase{repo: r}
}
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	
	if !regexp.MustCompile(`[!@#$%^&*()_+{}\[\]:;<>,.?~\\-]`).MatchString(password) {
		return false
	}

	return true
}

func (r *RegistrationUseCase) RegisterUser(phone, email, username,password string) (*domain.User, error) {
	
	existingUser, err := r.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailInUse
	}
	if !isValidPassword(password) {
		return nil, ErrFormat
	}

	if len(phone)!=9 {
		return nil, ErrPhone

	}
	newUser, err := domain.NewUser(phone, email, username,password) // Utilizamos el correo electrónico como identificador
	if err != nil {
		return nil, err
	}
    newUser, err = r.repo.InsertUser(newUser)
	return newUser, nil
}
