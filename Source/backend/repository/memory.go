package repository

import (
	"errors"
	"harbor/models"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type MemoryRepository struct {
	mu       sync.RWMutex
	users    map[int]*models.User
	ships    map[int]*models.Ship
	visits   map[int]*models.Visit
	userSeq  int
	shipSeq  int
	visitSeq int
}

func NewMemoryRepository() *MemoryRepository {
	r := &MemoryRepository{
		users:  make(map[int]*models.User),
		ships:  make(map[int]*models.Ship),
		visits: make(map[int]*models.Visit),
	}
	r.seed()
	return r
}

func (r *MemoryRepository) seed() {
	hash := func(p string) string {
		b, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		return string(b)
	}
	a1, a2, a3 := 35, 28, 42
	r.users[1] = &models.User{ID: 1, Name: "Иван Петров", Email: "admin@harbor.ru", Age: &a1, Role: "admin", Password: hash("admin123"), CreatedAt: time.Now()}
	r.users[2] = &models.User{ID: 2, Name: "Мария Сидорова", Email: "dispatcher@harbor.ru", Age: &a2, Role: "dispatcher", Password: hash("disp123"), CreatedAt: time.Now()}
	r.users[3] = &models.User{ID: 3, Name: "Алексей Козлов", Email: "operator@harbor.ru", Age: &a3, Role: "operator", Password: hash("oper123"), CreatedAt: time.Now()}
	r.userSeq = 3

	r.ships[1] = &models.Ship{ID: 1, Name: "Северный Ветер", IMONumber: "IMO9876543", ShipType: "cargo", FlagCountry: "Россия", GrossTonnage: 15000, YearBuilt: 2010, CreatedAt: time.Now()}
	r.ships[2] = &models.Ship{ID: 2, Name: "Arctic Explorer", IMONumber: "IMO1234567", ShipType: "tanker", FlagCountry: "Норвегия", GrossTonnage: 50000, YearBuilt: 2015, CreatedAt: time.Now()}
	r.ships[3] = &models.Ship{ID: 3, Name: "Балтийская Звезда", IMONumber: "IMO5555555", ShipType: "passenger", FlagCountry: "Финляндия", GrossTonnage: 30000, YearBuilt: 2018, CreatedAt: time.Now()}
	r.shipSeq = 3

	arr1 := time.Now().Add(-48 * time.Hour)
	arr2 := time.Now().Add(-120 * time.Hour)
	dep2 := time.Now().Add(-96 * time.Hour)
	r.visits[1] = &models.Visit{ID: 1, ShipID: 1, BerthID: 1, ArrivalTime: arr1, Status: "active", Purpose: "Разгрузка контейнеров", ShipName: "Северный Ветер", BerthNumber: "A-01", CreatedAt: time.Now()}
	r.visits[2] = &models.Visit{ID: 2, ShipID: 2, BerthID: 3, ArrivalTime: arr2, DepartureTime: &dep2, Status: "completed", Purpose: "Погрузка нефти", ShipName: "Arctic Explorer", BerthNumber: "B-01", CreatedAt: time.Now()}
	r.visitSeq = 2
}

func (r *MemoryRepository) GetAllUsers() ([]models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.User, 0, len(r.users))
	for _, u := range r.users {
		out = append(out, *u)
	}
	return out, nil
}

func (r *MemoryRepository) GetUserByID(id int) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("пользователь не найден")
	}
	return u, nil
}

func (r *MemoryRepository) GetUserByEmail(email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("пользователь не найден")
}

func (r *MemoryRepository) CreateUser(u *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, x := range r.users {
		if x.Email == u.Email {
			return nil, errors.New("email уже занят")
		}
	}
	r.userSeq++
	u.ID = r.userSeq
	u.CreatedAt = time.Now()
	r.users[u.ID] = u
	return u, nil
}

func (r *MemoryRepository) UpdateUser(id int, upd *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("пользователь не найден")
	}
	if upd.Name != "" {
		u.Name = upd.Name
	}
	if upd.Role != "" {
		u.Role = upd.Role
	}
	if upd.Age != nil {
		u.Age = upd.Age
	}
	return u, nil
}

func (r *MemoryRepository) DeleteUser(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return errors.New("пользователь не найден")
	}
	delete(r.users, id)
	return nil
}

func (r *MemoryRepository) GetAllShips() ([]models.Ship, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.Ship, 0, len(r.ships))
	for _, s := range r.ships {
		out = append(out, *s)
	}
	return out, nil
}

func (r *MemoryRepository) GetShipByID(id int) (*models.Ship, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.ships[id]
	if !ok {
		return nil, errors.New("судно не найдено")
	}
	return s, nil
}

func (r *MemoryRepository) CreateShip(s *models.Ship) (*models.Ship, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, x := range r.ships {
		if x.IMONumber == s.IMONumber {
			return nil, errors.New("номер ИМО уже существует")
		}
	}
	r.shipSeq++
	s.ID = r.shipSeq
	s.CreatedAt = time.Now()
	r.ships[s.ID] = s
	return s, nil
}

func (r *MemoryRepository) UpdateShip(id int, upd *models.Ship) (*models.Ship, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.ships[id]
	if !ok {
		return nil, errors.New("судно не найдено")
	}
	if upd.Name != "" {
		s.Name = upd.Name
	}
	if upd.ShipType != "" {
		s.ShipType = upd.ShipType
	}
	if upd.FlagCountry != "" {
		s.FlagCountry = upd.FlagCountry
	}
	if upd.GrossTonnage > 0 {
		s.GrossTonnage = upd.GrossTonnage
	}
	if upd.YearBuilt > 0 {
		s.YearBuilt = upd.YearBuilt
	}
	return s, nil
}

func (r *MemoryRepository) DeleteShip(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.ships[id]; !ok {
		return errors.New("судно не найдено")
	}
	delete(r.ships, id)
	return nil
}

func (r *MemoryRepository) GetAllVisits() ([]models.Visit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.Visit, 0, len(r.visits))
	for _, v := range r.visits {
		vv := *v
		if s, ok := r.ships[v.ShipID]; ok {
			vv.ShipName = s.Name
		}
		out = append(out, vv)
	}
	return out, nil
}

func (r *MemoryRepository) GetVisitByID(id int) (*models.Visit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.visits[id]
	if !ok {
		return nil, errors.New("визит не найден")
	}
	vv := *v
	if s, ok := r.ships[v.ShipID]; ok {
		vv.ShipName = s.Name
	}
	return &vv, nil
}

func (r *MemoryRepository) CreateVisit(v *models.Visit) (*models.Visit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.ships[v.ShipID]; !ok {
		return nil, errors.New("судно не найдено")
	}
	r.visitSeq++
	v.ID = r.visitSeq
	v.Status = "planned"
	v.CreatedAt = time.Now()
	if s, ok := r.ships[v.ShipID]; ok {
		v.ShipName = s.Name
	}
	r.visits[v.ID] = v
	return v, nil
}

func (r *MemoryRepository) UpdateVisit(id int, upd *models.Visit) (*models.Visit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.visits[id]
	if !ok {
		return nil, errors.New("визит не найден")
	}
	if upd.Status != "" {
		v.Status = upd.Status
	}
	if upd.DepartureTime != nil {
		v.DepartureTime = upd.DepartureTime
	}
	return v, nil
}

func (r *MemoryRepository) DeleteVisit(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.visits[id]; !ok {
		return errors.New("визит не найден")
	}
	delete(r.visits, id)
	return nil
}
