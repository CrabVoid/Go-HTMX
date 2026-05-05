package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

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
	UserID      uuid.UUID
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

type PositionSkill struct {
	PositionID uuid.UUID
	SkillID    uuid.UUID
}

type Application struct {
	ID         uuid.UUID
	PositionID uuid.UUID
	UserID     uuid.UUID
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

type Contact struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	Name        string
	Role        *string
	Email       *string
	LinkedinURL *string
}

// Custom types for joined queries
type PositionWithDetails struct {
	Position
	CompanyName string
	Skills      []Skill
}

type ApplicationWithDetails struct {
	Application
	PositionTitle string
	CompanyName   string
	CompanyID     uuid.UUID
	Interviews    []Interview
}
