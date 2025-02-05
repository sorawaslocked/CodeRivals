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
