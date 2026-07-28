package api

import (
	"job4j.ru/sharetrip-notification/internal/repository"
	"job4j.ru/sharetrip-notification/internal/service"
)

type Server struct {
	Repository          repository.NotificationRepository
	NotificationService *service.NotificationService
}

func NewServer(repo repository.NotificationRepository,
) *Server {
	s := &Server{
		Repository: repo,
	}

	return s
}
