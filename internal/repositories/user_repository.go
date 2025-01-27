package repositories

import (
	"database/sql"
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
