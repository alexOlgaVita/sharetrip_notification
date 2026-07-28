package domain

import "errors"

var (
	ErrNotificationNotFound             = errors.New("notification not found")
	ErrForbidden                        = errors.New("forbidden")
	ErrConflict                         = errors.New("conflict")
	ErrAlreadyExists                    = errors.New("already exists")
	ErrNotAllowedCurrentStatusToPublish = errors.New("current status isn't allowed to publish")
)
