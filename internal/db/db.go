package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Query(context.Context, string, ...interface{}) (any, error)
	QueryRow(context.Context, string, ...interface{}) any
	Exec(context.Context, string, ...interface{}) (any, error)
}

func New(db *pgxpool.Pool) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db *pgxpool.Pool
}

func (q *Queries) ListCompanies(ctx context.Context) ([]Company, error) {
	return []Company{}, nil
}

func (q *Queries) CreateCompany(ctx context.Context, arg CreateCompanyParams) (Company, error) {
	return Company{
		Name:     arg.Name,
		Website:  arg.Website,
		Industry: arg.Industry,
	}, nil
}

func (q *Queries) ListPositions(ctx context.Context) ([]PositionWithDetails, error) {
	return []PositionWithDetails{}, nil
}

func (q *Queries) ListSkills(ctx context.Context) ([]Skill, error) {
	return []Skill{}, nil
}

func (q *Queries) GetApplication(ctx context.Context, id uuid.UUID) (ApplicationWithDetails, error) {
	return ApplicationWithDetails{}, nil
}

func (q *Queries) CreateInterview(ctx context.Context, arg CreateInterviewParams) (Interview, error) {
	return Interview{}, nil
}

type CreateInterviewParams struct {
	ApplicationID uuid.UUID
	StageName     string
	ScheduledAt   time.Time
	Notes         *string
}

type CreateCompanyParams struct {
	Name     string
	Website  *string
	Industry *string
}
