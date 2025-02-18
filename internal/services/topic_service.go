package services

import (
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
)

type TopicService struct {
	topicRepo repositories.TopicRepository
}

func NewTopicService(repo repositories.TopicRepository) *TopicService {
	return &TopicService{topicRepo: repo}
}

func (s *TopicService) GetAllTopics() ([]*entities.Topic, error) {
	return s.topicRepo.GetAll()
}

func (s *TopicService) GetTopic(id int) (*entities.Topic, error) {
	return s.topicRepo.Get(id)
}

func (s *TopicService) Create(name string) error {
	return s.topicRepo.Create(name)
}

func (s *TopicService) Delete(id int) error {
	return s.topicRepo.Delete(id)
}

func (s *TopicService) UpdateTopic(id int, name string) error {
	topic := &entities.Topic{
		ID:   id,
		Name: name,
	}
	return s.topicRepo.Update(topic)
}

func (s *TopicService) GetProblemCountForTopic(topicId int) (int, error) {
	var count int
	count, err := s.topicRepo.GetProblemCount(topicId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
