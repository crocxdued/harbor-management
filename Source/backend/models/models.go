package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       *int      `json:"age,omitempty"`
	Role      string    `json:"role"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UserCreateRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Age      *int   `json:"age"`
	Role     string `json:"role"     binding:"required,oneof=admin dispatcher operator"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserUpdateRequest struct {
	Name string `json:"name"`
	Role string `json:"role" binding:"omitempty,oneof=admin dispatcher operator"`
	Age  *int   `json:"age"`
}

type Ship struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	IMONumber    string    `json:"imo_number"`
	ShipType     string    `json:"ship_type"`
	FlagCountry  string    `json:"flag_country"`
	GrossTonnage int       `json:"gross_tonnage"`
	YearBuilt    int       `json:"year_built"`
	CreatedAt    time.Time `json:"created_at"`
}

type ShipCreateRequest struct {
	Name         string `json:"name"          binding:"required"`
	IMONumber    string `json:"imo_number"    binding:"required"`
	ShipType     string `json:"ship_type"     binding:"required,oneof=cargo tanker passenger warship"`
	FlagCountry  string `json:"flag_country"  binding:"required"`
	GrossTonnage int    `json:"gross_tonnage" binding:"required,min=1"`
	YearBuilt    int    `json:"year_built"    binding:"required"`
}

type ShipUpdateRequest struct {
	Name         string `json:"name"`
	ShipType     string `json:"ship_type"     binding:"omitempty,oneof=cargo tanker passenger warship"`
	FlagCountry  string `json:"flag_country"`
	GrossTonnage int    `json:"gross_tonnage"`
	YearBuilt    int    `json:"year_built"`
}

type Visit struct {
	ID            int        `json:"id"`
	ShipID        int        `json:"ship_id"`
	BerthID       int        `json:"berth_id"`
	ArrivalTime   time.Time  `json:"arrival_time"`
	DepartureTime *time.Time `json:"departure_time,omitempty"`
	Status        string     `json:"status"`
	Purpose       string     `json:"purpose"`
	CreatedAt     time.Time  `json:"created_at"`

	ShipName    string `json:"ship_name,omitempty"`
	BerthNumber string `json:"berth_number,omitempty"`
}

type VisitCreateRequest struct {
	ShipID      int       `json:"ship_id"      binding:"required"`
	BerthID     int       `json:"berth_id"     binding:"required"`
	ArrivalTime time.Time `json:"arrival_time" binding:"required"`
	Purpose     string    `json:"purpose"      binding:"required"`
}

type VisitUpdateRequest struct {
	Status        string     `json:"status"         binding:"omitempty,oneof=planned active completed cancelled"`
	DepartureTime *time.Time `json:"departure_time"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Claims struct {
	UserID int
	Email  string
	Role   string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MsgResponse struct {
	Message string `json:"message"`
}
