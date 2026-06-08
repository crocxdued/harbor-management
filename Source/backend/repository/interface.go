package repository

import "harbor/models"

type Repository interface {
	// Users
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(u *models.User) (*models.User, error)
	UpdateUser(id int, u *models.User) (*models.User, error)
	DeleteUser(id int) error

	// Ships
	GetAllShips() ([]models.Ship, error)
	GetShipByID(id int) (*models.Ship, error)
	CreateShip(s *models.Ship) (*models.Ship, error)
	UpdateShip(id int, s *models.Ship) (*models.Ship, error)
	DeleteShip(id int) error

	// Visits
	GetAllVisits() ([]models.Visit, error)
	GetVisitByID(id int) (*models.Visit, error)
	CreateVisit(v *models.Visit) (*models.Visit, error)
	UpdateVisit(id int, v *models.Visit) (*models.Visit, error)
	DeleteVisit(id int) error
}
