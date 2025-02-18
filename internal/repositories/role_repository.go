package repositories

import (
	"database/sql"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type RoleRepository interface {
	Count() (int, error)
	Get(id int) (*entities.Role, error)
	GetAll() ([]*entities.Role, error)
	GetUserRoles(userId int) ([]*entities.Role, error)
	IsUserAdmin(userId int) (bool, error)
	AssignRole(userId int, roleId int) error
	RemoveRole(userId int, roleId int) error
}

type PGRoleRepository struct {
	db *sql.DB
}

func NewPGRoleRepository(db *sql.DB) *PGRoleRepository {
	return &PGRoleRepository{db: db}
}

func (repo *PGRoleRepository) Count() (int, error) {
	var count int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM roles").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *PGRoleRepository) Get(id int) (*entities.Role, error) {
	role := &entities.Role{}
	err := repo.db.QueryRow(`
        SELECT id, name, created_at 
        FROM roles 
        WHERE id = $1
    `, id).Scan(&role.ID, &role.Name, &role.CreatedAt)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (repo *PGRoleRepository) GetAll() ([]*entities.Role, error) {
	rows, err := repo.db.Query(`
        SELECT id, name, created_at 
        FROM roles
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entities.Role
	for rows.Next() {
		role := &entities.Role{}
		err = rows.Scan(&role.ID, &role.Name, &role.CreatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (repo *PGRoleRepository) GetUserRoles(userId int) ([]*entities.Role, error) {
	rows, err := repo.db.Query(`
        SELECT r.id, r.name, r.created_at
        FROM roles r
        JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = $1
    `, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entities.Role
	for rows.Next() {
		role := &entities.Role{}
		err = rows.Scan(&role.ID, &role.Name, &role.CreatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (repo *PGRoleRepository) IsUserAdmin(userId int) (bool, error) {
	var exists bool
	err := repo.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = $1 AND r.name = 'admin'
        )
    `, userId).Scan(&exists)

	return exists, err
}

func (repo *PGRoleRepository) AssignRole(userId int, roleId int) error {
	_, err := repo.db.Exec(`
        INSERT INTO user_roles (user_id, role_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, role_id) DO NOTHING
    `, userId, roleId)

	return err
}

func (repo *PGRoleRepository) RemoveRole(userId int, roleId int) error {
	_, err := repo.db.Exec(`
        DELETE FROM user_roles
        WHERE user_id = $1 AND role_id = $2
    `, userId, roleId)

	return err
}
