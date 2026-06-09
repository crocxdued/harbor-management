package service

import (
	"errors"
	"harbor/models"
	"harbor/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      repository.Repository
	jwtSecret string
	jwtExpiry time.Duration
}

func New(repo repository.Repository, secret string, expiryHours int) *Service {
	return &Service{repo: repo, jwtSecret: secret, jwtExpiry: time.Duration(expiryHours) * time.Hour}
}

func (s *Service) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	u, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		return nil, errors.New("неверный email или пароль")
	}
	token, err := s.makeToken(u)
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{Token: token, User: *u}, nil
}

func (s *Service) makeToken(u *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"role":    u.Role,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
}

func (s *Service) ParseToken(tok string) (*models.Claims, error) {
	t, err := jwt.Parse(tok, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !t.Valid {
		return nil, errors.New("недействительный токен")
	}
	c := t.Claims.(jwt.MapClaims)
	return &models.Claims{
		UserID: int(c["user_id"].(float64)),
		Email:  c["email"].(string),
		Role:   c["role"].(string),
	}, nil
}

func (s *Service) GetAllUsers() ([]models.User, error)      { return s.repo.GetAllUsers() }
func (s *Service) GetUserByID(id int) (*models.User, error) { return s.repo.GetUserByID(id) }

func (s *Service) CreateUser(req models.UserCreateRequest) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("ошибка хеширования пароля")
	}
	return s.repo.CreateUser(&models.User{
		Name: req.Name, Email: req.Email, Age: req.Age, Role: req.Role, Password: string(hash),
	})
}

func (s *Service) UpdateUser(id int, req models.UserUpdateRequest) (*models.User, error) {
	return s.repo.UpdateUser(id, &models.User{Name: req.Name, Role: req.Role, Age: req.Age})
}

func (s *Service) DeleteUser(id int) error { return s.repo.DeleteUser(id) }

func (s *Service) GetAllShips() ([]models.Ship, error)      { return s.repo.GetAllShips() }
func (s *Service) GetShipByID(id int) (*models.Ship, error) { return s.repo.GetShipByID(id) }

func (s *Service) CreateShip(req models.ShipCreateRequest) (*models.Ship, error) {
	return s.repo.CreateShip(&models.Ship{
		Name: req.Name, IMONumber: req.IMONumber, ShipType: req.ShipType,
		FlagCountry: req.FlagCountry, GrossTonnage: req.GrossTonnage, YearBuilt: req.YearBuilt,
	})
}

func (s *Service) UpdateShip(id int, req models.ShipUpdateRequest) (*models.Ship, error) {
	return s.repo.UpdateShip(id, &models.Ship{
		Name: req.Name, ShipType: req.ShipType, FlagCountry: req.FlagCountry,
		GrossTonnage: req.GrossTonnage, YearBuilt: req.YearBuilt,
	})
}

func (s *Service) DeleteShip(id int) error { return s.repo.DeleteShip(id) }

func (s *Service) GetAllVisits() ([]models.Visit, error)      { return s.repo.GetAllVisits() }
func (s *Service) GetVisitByID(id int) (*models.Visit, error) { return s.repo.GetVisitByID(id) }

func (s *Service) CreateVisit(req models.VisitCreateRequest) (*models.Visit, error) {
	return s.repo.CreateVisit(&models.Visit{
		ShipID: req.ShipID, BerthID: req.BerthID,
		ArrivalTime: req.ArrivalTime, Purpose: req.Purpose,
	})
}

func (s *Service) UpdateVisit(id int, req models.VisitUpdateRequest) (*models.Visit, error) {
	return s.repo.UpdateVisit(id, &models.Visit{Status: req.Status, DepartureTime: req.DepartureTime})
}

func (s *Service) DeleteVisit(id int) error { return s.repo.DeleteVisit(id) }
