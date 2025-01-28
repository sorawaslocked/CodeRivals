package repositories

import (
	"database/sql"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type TopicRepository interface {
	Count() (uint64, error)
	Get(id uint64) (*entities.Topic, error)
	GetAll() ([]*entities.Topic, error)
	GetAllForProblem(problemId uint64) ([]*entities.Topic, error)
	Create(name string) error
	Update(newName string) error
	Delete(id uint64) error
}

type PGTopicRepository struct {
	db *sql.DB
}

func NewPGTopicRepository(db *sql.DB) *PGTopicRepository {
	return &PGTopicRepository{db: db}
}

func (repo *PGTopicRepository) Count() (uint64, error) {
	var count uint64

	stmt := "SELECT COUNT(*) FROM topics"

	err := repo.db.QueryRow(stmt).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *PGTopicRepository) Get(id uint64) (*entities.Topic, error) {
	topic := &entities.Topic{
		ID: id,
	}

	stmt := "SELECT name FROM topics WHERE id = $1"

	err := repo.db.QueryRow(stmt, id).Scan(&topic.Name)

	if err != nil {
		return nil, err
	}

	return topic, nil
}

func (repo *PGTopicRepository) Create(name string) error {
	stmt := "INSERT INTO topics (name) VALUES ($1)"

	_, err := repo.db.Exec(stmt, name)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PGTopicRepository) GetAll() ([]*entities.Topic, error) {
	stmt := "SELECT id, name FROM topics"

	rows, err := repo.db.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var topics []*entities.Topic

	for rows.Next() {
		topic := &entities.Topic{}

		err = rows.Scan(&topic.ID, &topic.Name)

		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

func (repo *PGTopicRepository) GetAllForProblem(problemId uint64) ([]*entities.Topic, error) {
	stmt := `SELECT t.name, t.id
	FROM problem_topics pt
	JOIN topics t ON pt.topic_id = t.id
	WHERE pt.problem_id = $1`

	rows, err := repo.db.Query(stmt, problemId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var topics []*entities.Topic

	for rows.Next() {
		topic := &entities.Topic{}

		err = rows.Scan(&topic.Name, &topic.ID)

		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

func (repo *PGTopicRepository) Update(newName string) error {
	stmt := "UPDATE topics SET name = $1 WHERE id = $2"

	_, err := repo.db.Exec(stmt, newName)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PGTopicRepository) Delete(id uint64) error {
	stmt := "DELETE FROM topics WHERE id = $1"

	_, err := repo.db.Exec(stmt, id)

	if err != nil {
		return err
	}

	return nil
}
