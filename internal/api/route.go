package api

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Route(route fiber.Router) {
	route.Get("/ready/", s.DoPing)

	route.Post(
		"notification",
		s.CreateNotification,
	)

	route.Get(
		"/notification/:notificationId",
		s.GetNotification,
	)
}
