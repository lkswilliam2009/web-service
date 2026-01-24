package models

import "github.com/google/uuid"

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GroupArr struct {
	GroupID   uuid.UUID `json:"group_id"`
	GroupName string    `json:"group_name"`
}