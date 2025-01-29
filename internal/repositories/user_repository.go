package repositories

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type UserRepository interface {
	// Basic CRUD
	Count() (uint64, error)
	Create(user *entities.User) error
	Get(id uint64) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	GetByUsername(username string) (*entities.User, error)
	GetAll() ([]*entities.User, error)
	Update(user *entities.User) error
	UpdatePoints(id uint64, points int) error
	Delete(id uint64) error

	// User engagement
	//GetComments(userID uint64) ([]*entities.Comment, error)
}

type PGUserRepository struct {
	db *sql.DB
}

func NewPGUserRepository(db *sql.DB) *PGUserRepository { return &PGUserRepository{db: db} }

func (repo *PGUserRepository) Count() (uint64, error) {
	var count uint64
	err := repo.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *PGUserRepository) Create(user *entities.User) error {
	stmt := `INSERT INTO users (username, email, hashed_password, points, created_at, updated_at)
	VALUES ($1, $2, $3, 0, NOW(), NOW()) RETURNING id`

	err := repo.db.QueryRow(
		stmt,
		user.Username,
		user.Email,
		user.HashedPassword,
	).Scan(&user.ID)

	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)

	// Check for duplicate username and email error
	if ok && pqErr.Code == "23505" {
		if pqErr.Constraint == "users_username_key" {
			return ErrDuplicateUsername
		}
		if pqErr.Constraint == "users_email_key" {
			return ErrDuplicateEmail
		}
	}

	return err
}

func (repo *PGUserRepository) Get(id uint64) (*entities.User, error) {
	user := &entities.User{}
	stmt := `SELECT id, username, email, hashed_password, points, created_at, updated_at
	FROM users WHERE id = $1`

	err := repo.db.QueryRow(stmt, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *PGUserRepository) GetByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	stmt := `SELECT id, username, email, hashed_password, points, created_at, updated_at
	FROM users WHERE email = $1`

	err := repo.db.QueryRow(stmt, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *PGUserRepository) GetByUsername(username string) (*entities.User, error) {
	user := &entities.User{}
	stmt := `SELECT id, username, email, hashed_password, points, created_at, updated_at
	FROM users WHERE username = $1`

	err := repo.db.QueryRow(stmt, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *PGUserRepository) GetAll() ([]*entities.User, error) {
	stmt := `SELECT id, username, email, hashed_password, points, created_at, updated_at
	FROM users`

	rows, err := repo.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User

	for rows.Next() {
		user := &entities.User{}
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.HashedPassword,
			&user.Points,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repo *PGUserRepository) Update(user *entities.User) error {
	stmt := `UPDATE users
	SET username = $1, email = $2, hashed_password = $3, points = $4, updated_at = NOW()
	WHERE id = $5`

	_, err := repo.db.Exec(
		stmt,
		user.Username,
		user.Email,
		user.HashedPassword,
		user.Points,
		user.ID,
	)

	return err
}

func (repo *PGUserRepository) UpdatePoints(id uint64, points int) error {
	stmt := `UPDATE users SET points = points + $1, updated_at = NOW() WHERE id = $2`
	_, err := repo.db.Exec(stmt, points, id)
	return err
}

func (repo *PGUserRepository) Delete(id uint64) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete related records first
	stmts := []string{
		`DELETE FROM comments WHERE user_id = $1`,
		`DELETE FROM problem_submissions WHERE user_id = $1`,
		`DELETE FROM users WHERE id = $1`,
	}

	for _, stmt := range stmts {
		_, err = tx.Exec(stmt, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
