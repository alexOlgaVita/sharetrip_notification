package domain

import (
	"job4j.ru/sharetrip-notification/internal/dto"
)

type Notification struct {
	Notifications []dto.Notification
}

func NewNotification() *Notification {
	return &Notification{}
}

func (sht *Notification) AddNotification(notification dto.Notification) error {
	_, ok := sht.indexOf(notification.ID)
	if ok {
		return ErrAlreadyExists
	}
	sht.Notifications = append(sht.Notifications, notification)
	return nil
}

func (sht *Notification) GetNotification() []dto.Notification {
	return sht.Notifications
}

func (sht *Notification) indexOf(id string) (int, bool) {
	for i, notification := range sht.Notifications {
		if notification.ID == id {
			return i, true
		}
	}
	return -1, false
}
