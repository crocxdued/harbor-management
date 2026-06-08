package repository

import (
	"context"
	"errors"
	"fmt"
	"harbor/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(host, port, user, password, dbname string) (*PostgresRepository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil { return nil, err }
	if err := pool.Ping(ctx); err != nil { return nil, err }
	return &PostgresRepository{pool: pool}, nil
}

func (r *PostgresRepository) Close() { r.pool.Close() }

// ── Users ────────────────────────────────────────────────────

func (r *PostgresRepository) GetAllUsers() ([]models.User, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, name, email, age, role, created_at FROM users ORDER BY id`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Role, &u.CreatedAt)
		out = append(out, u)
	}
	return out, nil
}

func (r *PostgresRepository) GetUserByID(id int) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(context.Background(),
		`SELECT id, name, email, age, role, created_at FROM users WHERE id=$1`, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Role, &u.CreatedAt)
	if err != nil { return nil, errors.New("пользователь не найден") }
	return &u, nil
}

func (r *PostgresRepository) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(context.Background(),
		`SELECT id, name, email, age, role, password_hash, created_at FROM users WHERE email=$1`, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Role, &u.Password, &u.CreatedAt)
	if err != nil { return nil, errors.New("пользователь не найден") }
	return &u, nil
}

func (r *PostgresRepository) CreateUser(u *models.User) (*models.User, error) {
	err := r.pool.QueryRow(context.Background(),
		`INSERT INTO users(name,email,age,role,password_hash) VALUES($1,$2,$3,$4,$5) RETURNING id,created_at`,
		u.Name, u.Email, u.Age, u.Role, u.Password).Scan(&u.ID, &u.CreatedAt)
	if err != nil { return nil, fmt.Errorf("ошибка создания пользователя: %w", err) }
	return u, nil
}

func (r *PostgresRepository) UpdateUser(id int, upd *models.User) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(context.Background(),
		`UPDATE users SET
			name = COALESCE(NULLIF($1,''), name),
			role = COALESCE(NULLIF($2,''), role),
			age  = COALESCE($3, age)
		 WHERE id=$4 RETURNING id,name,email,age,role,created_at`,
		upd.Name, upd.Role, upd.Age, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Role, &u.CreatedAt)
	if err != nil { return nil, errors.New("пользователь не найден") }
	return &u, nil
}

func (r *PostgresRepository) DeleteUser(id int) error {
	res, err := r.pool.Exec(context.Background(), `DELETE FROM users WHERE id=$1`, id)
	if err != nil { return err }
	if res.RowsAffected() == 0 { return errors.New("пользователь не найден") }
	return nil
}

// ── Ships ────────────────────────────────────────────────────

func (r *PostgresRepository) GetAllShips() ([]models.Ship, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id,name,imo_number,ship_type,flag_country,gross_tonnage,year_built,created_at FROM ships ORDER BY id`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []models.Ship
	for rows.Next() {
		var s models.Ship
		rows.Scan(&s.ID, &s.Name, &s.IMONumber, &s.ShipType, &s.FlagCountry, &s.GrossTonnage, &s.YearBuilt, &s.CreatedAt)
		out = append(out, s)
	}
	return out, nil
}

func (r *PostgresRepository) GetShipByID(id int) (*models.Ship, error) {
	var s models.Ship
	err := r.pool.QueryRow(context.Background(),
		`SELECT id,name,imo_number,ship_type,flag_country,gross_tonnage,year_built,created_at FROM ships WHERE id=$1`, id).
		Scan(&s.ID, &s.Name, &s.IMONumber, &s.ShipType, &s.FlagCountry, &s.GrossTonnage, &s.YearBuilt, &s.CreatedAt)
	if err != nil { return nil, errors.New("судно не найдено") }
	return &s, nil
}

func (r *PostgresRepository) CreateShip(s *models.Ship) (*models.Ship, error) {
	err := r.pool.QueryRow(context.Background(),
		`INSERT INTO ships(name,imo_number,ship_type,flag_country,gross_tonnage,year_built)
		 VALUES($1,$2,$3,$4,$5,$6) RETURNING id,created_at`,
		s.Name, s.IMONumber, s.ShipType, s.FlagCountry, s.GrossTonnage, s.YearBuilt).
		Scan(&s.ID, &s.CreatedAt)
	if err != nil { return nil, fmt.Errorf("ошибка создания судна: %w", err) }
	return s, nil
}

func (r *PostgresRepository) UpdateShip(id int, upd *models.Ship) (*models.Ship, error) {
	var s models.Ship
	err := r.pool.QueryRow(context.Background(),
		`UPDATE ships SET
			name          = COALESCE(NULLIF($1,''), name),
			ship_type     = COALESCE(NULLIF($2,''), ship_type),
			flag_country  = COALESCE(NULLIF($3,''), flag_country),
			gross_tonnage = CASE WHEN $4>0 THEN $4 ELSE gross_tonnage END,
			year_built    = CASE WHEN $5>0 THEN $5 ELSE year_built END
		 WHERE id=$6 RETURNING id,name,imo_number,ship_type,flag_country,gross_tonnage,year_built,created_at`,
		upd.Name, upd.ShipType, upd.FlagCountry, upd.GrossTonnage, upd.YearBuilt, id).
		Scan(&s.ID, &s.Name, &s.IMONumber, &s.ShipType, &s.FlagCountry, &s.GrossTonnage, &s.YearBuilt, &s.CreatedAt)
	if err != nil { return nil, errors.New("судно не найдено") }
	return &s, nil
}

func (r *PostgresRepository) DeleteShip(id int) error {
	res, err := r.pool.Exec(context.Background(), `DELETE FROM ships WHERE id=$1`, id)
	if err != nil { return err }
	if res.RowsAffected() == 0 { return errors.New("судно не найдено") }
	return nil
}

// ── Visits ───────────────────────────────────────────────────

func (r *PostgresRepository) GetAllVisits() ([]models.Visit, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT v.id,v.ship_id,v.berth_id,v.arrival_time,v.departure_time,
		        v.status,v.purpose,v.created_at,s.name,b.number
		 FROM visits v
		 JOIN ships  s ON s.id=v.ship_id
		 JOIN berths b ON b.id=v.berth_id
		 ORDER BY v.arrival_time DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []models.Visit
	for rows.Next() {
		var v models.Visit
		rows.Scan(&v.ID, &v.ShipID, &v.BerthID, &v.ArrivalTime, &v.DepartureTime,
			&v.Status, &v.Purpose, &v.CreatedAt, &v.ShipName, &v.BerthNumber)
		out = append(out, v)
	}
	return out, nil
}

func (r *PostgresRepository) GetVisitByID(id int) (*models.Visit, error) {
	var v models.Visit
	err := r.pool.QueryRow(context.Background(),
		`SELECT v.id,v.ship_id,v.berth_id,v.arrival_time,v.departure_time,
		        v.status,v.purpose,v.created_at,s.name,b.number
		 FROM visits v
		 JOIN ships  s ON s.id=v.ship_id
		 JOIN berths b ON b.id=v.berth_id
		 WHERE v.id=$1`, id).
		Scan(&v.ID, &v.ShipID, &v.BerthID, &v.ArrivalTime, &v.DepartureTime,
			&v.Status, &v.Purpose, &v.CreatedAt, &v.ShipName, &v.BerthNumber)
	if err != nil { return nil, errors.New("визит не найден") }
	return &v, nil
}

func (r *PostgresRepository) CreateVisit(v *models.Visit) (*models.Visit, error) {
	err := r.pool.QueryRow(context.Background(),
		`INSERT INTO visits(ship_id,berth_id,arrival_time,status,purpose)
		 VALUES($1,$2,$3,'planned',$4) RETURNING id,created_at`,
		v.ShipID, v.BerthID, v.ArrivalTime, v.Purpose).Scan(&v.ID, &v.CreatedAt)
	if err != nil { return nil, fmt.Errorf("ошибка создания визита: %w", err) }
	v.Status = "planned"
	return v, nil
}

func (r *PostgresRepository) UpdateVisit(id int, upd *models.Visit) (*models.Visit, error) {
	var v models.Visit
	err := r.pool.QueryRow(context.Background(),
		`UPDATE visits SET
			status         = COALESCE(NULLIF($1,''), status),
			departure_time = COALESCE($2, departure_time)
		 WHERE id=$3
		 RETURNING id,ship_id,berth_id,arrival_time,departure_time,status,purpose,created_at`,
		upd.Status, upd.DepartureTime, id).
		Scan(&v.ID, &v.ShipID, &v.BerthID, &v.ArrivalTime, &v.DepartureTime,
			&v.Status, &v.Purpose, &v.CreatedAt)
	if err != nil { return nil, errors.New("визит не найден") }
	// освобождаем причал если визит закрыт
	if upd.Status == "completed" || upd.Status == "cancelled" {
		r.pool.Exec(context.Background(),
			`UPDATE berths SET is_available=true WHERE id=$1`, v.BerthID)
	}
	return &v, nil
}

func (r *PostgresRepository) DeleteVisit(id int) error {
	res, err := r.pool.Exec(context.Background(), `DELETE FROM visits WHERE id=$1`, id)
	if err != nil { return err }
	if res.RowsAffected() == 0 { return errors.New("визит не найден") }
	return nil
}
