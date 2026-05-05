package db

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID        uuid.UUID
	Name      string
	Website   *string
	Industry  *string
	CreatedAt time.Time
}

type Position struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	Title       string
	Location    *string
	WorkMode    string
	SalaryRange *string
	PostURL     *string
	CreatedAt   time.Time
}

type Skill struct {
	ID   uuid.UUID
	Name string
}

type PositionWithDetails struct {
	Position
	CompanyName string
	Skills      []Skill
}

type Application struct {
	ID         uuid.UUID
	PositionID uuid.UUID
	Status     string
	Source     *string
	AppliedAt  time.Time
	Notes      *string
}

type Interview struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	StageName     string
	ScheduledAt   time.Time
	Notes         *string
	Feedback      *string
}

type ApplicationWithDetails struct {
	Application
	PositionTitle string
	CompanyName   string
	Interviews    []Interview
}
