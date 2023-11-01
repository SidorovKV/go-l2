package model

import "time"

type Event struct {
	Id          uint      `json:"id,omitempty"`
	Date        time.Time `json:"date"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatorId   uint      `json:"user_id,omitempty"`
}
