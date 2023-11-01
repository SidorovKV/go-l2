package service

import (
	"dev11/model"
	"time"
)

type Repository interface {
	Add(event model.Event) (model.Event, error)
	Update(event model.Event) (model.Event, error)
	Delete(id uint) error
	GetByDay(day time.Time) ([]model.Event, error)
	GetByWeek(startDay time.Time) ([]model.Event, error)
	GetByMonth(month time.Month, year int) ([]model.Event, error)
}

type FakeEventService struct {
}

func (f FakeEventService) CreateEvent(event model.Event) (model.Event, error) {
	return event, nil
}

func (f FakeEventService) UpdateEvent(event model.Event) (model.Event, error) {
	return event, nil
}

func (f FakeEventService) DeleteEvent(event model.Event) error {
	return nil
}

func (f FakeEventService) EventsForDay(userId uint, day time.Time) ([]model.Event, error) {
	return make([]model.Event, 0), nil
}

func (f FakeEventService) EventsForWeek(userId uint, startDay time.Time) ([]model.Event, error) {
	return make([]model.Event, 0), nil
}

func (f FakeEventService) EventsForMonth(userId uint, month time.Month, year int) ([]model.Event, error) {
	return make([]model.Event, 0), nil
}
