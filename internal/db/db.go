package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db *pgxpool.Pool
}

// Users
type CreateUserParams struct {
	Email        string
	PasswordHash string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at",
		arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.PasswordHash, &i.CreatedAt)
	return i, err
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, "SELECT id, email, password_hash, created_at FROM users WHERE email = $1", email)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.PasswordHash, &i.CreatedAt)
	return i, err
}

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, "SELECT id, email, password_hash, created_at FROM users WHERE id = $1", id)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.PasswordHash, &i.CreatedAt)
	return i, err
}

// Companies
func (q *Queries) ListCompanies(ctx context.Context) ([]Company, error) {
	rows, err := q.db.Query(ctx, "SELECT id, name, website, industry, created_at FROM companies ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Company
	for rows.Next() {
		var i Company
		if err := rows.Scan(&i.ID, &i.Name, &i.Website, &i.Industry, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (q *Queries) GetCompany(ctx context.Context, id uuid.UUID) (Company, error) {
	row := q.db.QueryRow(ctx, "SELECT id, name, website, industry, created_at FROM companies WHERE id = $1", id)
	var i Company
	err := row.Scan(&i.ID, &i.Name, &i.Website, &i.Industry, &i.CreatedAt)
	return i, err
}

type CreateCompanyParams struct {
	Name     string
	Website  *string
	Industry *string
}

func (q *Queries) CreateCompany(ctx context.Context, arg CreateCompanyParams) (Company, error) {
	row := q.db.QueryRow(ctx, "INSERT INTO companies (name, website, industry) VALUES ($1, $2, $3) RETURNING id, name, website, industry, created_at",
		arg.Name, arg.Website, arg.Industry)
	var i Company
	err := row.Scan(&i.ID, &i.Name, &i.Website, &i.Industry, &i.CreatedAt)
	return i, err
}

type UpdateCompanyParams struct {
	ID       uuid.UUID
	Name     string
	Website  *string
	Industry *string
}

func (q *Queries) UpdateCompany(ctx context.Context, arg UpdateCompanyParams) (Company, error) {
	row := q.db.QueryRow(ctx, "UPDATE companies SET name = $2, website = $3, industry = $4 WHERE id = $1 RETURNING id, name, website, industry, created_at",
		arg.ID, arg.Name, arg.Website, arg.Industry)
	var i Company
	err := row.Scan(&i.ID, &i.Name, &i.Website, &i.Industry, &i.CreatedAt)
	return i, err
}

func (q *Queries) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, "DELETE FROM companies WHERE id = $1", id)
	return err
}

// Positions
func (q *Queries) ListPositions(ctx context.Context, userID uuid.UUID) ([]PositionWithDetails, error) {
	rows, err := q.db.Query(ctx, `
		SELECT p.id, p.company_id, p.user_id, p.title, p.location, p.work_mode, p.salary_range, p.post_url, p.created_at, c.name as company_name 
		FROM positions p
		JOIN companies c ON p.company_id = c.id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PositionWithDetails
	for rows.Next() {
		var i PositionWithDetails
		if err := rows.Scan(&i.ID, &i.CompanyID, &i.UserID, &i.Title, &i.Location, &i.WorkMode, &i.SalaryRange, &i.PostURL, &i.CreatedAt, &i.CompanyName); err != nil {
			return nil, err
		}
		// Fetch skills for this position
		skills, err := q.ListSkillsForPosition(ctx, i.ID)
		if err == nil {
			i.Skills = skills
		}
		items = append(items, i)
	}
	return items, nil
}

func (q *Queries) ListSkillsForPosition(ctx context.Context, positionID uuid.UUID) ([]Skill, error) {
	rows, err := q.db.Query(ctx, `
		SELECT s.id, s.name 
		FROM skills s
		JOIN position_skills ps ON s.id = ps.skill_id
		WHERE ps.position_id = $1`, positionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Skill
	for rows.Next() {
		var i Skill
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (q *Queries) GetPosition(ctx context.Context, id uuid.UUID) (PositionWithDetails, error) {
	row := q.db.QueryRow(ctx, `
		SELECT p.id, p.company_id, p.user_id, p.title, p.location, p.work_mode, p.salary_range, p.post_url, p.created_at, c.name as company_name 
		FROM positions p
		JOIN companies c ON p.company_id = c.id
		WHERE p.id = $1`, id)
	var i PositionWithDetails
	err := row.Scan(&i.ID, &i.CompanyID, &i.UserID, &i.Title, &i.Location, &i.WorkMode, &i.SalaryRange, &i.PostURL, &i.CreatedAt, &i.CompanyName)
	if err != nil {
		return i, err
	}
	i.Skills, _ = q.ListSkillsForPosition(ctx, i.ID)
	return i, nil
}

type CreatePositionParams struct {
	CompanyID   uuid.UUID
	UserID      uuid.UUID
	Title       string
	Location    *string
	WorkMode    string
	SalaryRange *string
	PostURL     *string
}

func (q *Queries) CreatePosition(ctx context.Context, arg CreatePositionParams) (Position, error) {
	row := q.db.QueryRow(ctx, "INSERT INTO positions (company_id, user_id, title, location, work_mode, salary_range, post_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, company_id, user_id, title, location, work_mode, salary_range, post_url, created_at",
		arg.CompanyID, arg.UserID, arg.Title, arg.Location, arg.WorkMode, arg.SalaryRange, arg.PostURL)
	var i Position
	err := row.Scan(&i.ID, &i.CompanyID, &i.UserID, &i.Title, &i.Location, &i.WorkMode, &i.SalaryRange, &i.PostURL, &i.CreatedAt)
	return i, err
}

// Applications
func (q *Queries) ListApplications(ctx context.Context, userID uuid.UUID) ([]ApplicationWithDetails, error) {
	rows, err := q.db.Query(ctx, `
		SELECT a.id, a.position_id, a.user_id, a.status, a.source, a.applied_at, a.notes, p.title as position_title, c.name as company_name
		FROM applications a
		JOIN positions p ON a.position_id = p.id
		JOIN companies c ON p.company_id = c.id
		WHERE a.user_id = $1
		ORDER BY a.applied_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ApplicationWithDetails
	for rows.Next() {
		var i ApplicationWithDetails
		if err := rows.Scan(&i.ID, &i.PositionID, &i.UserID, &i.Status, &i.Source, &i.AppliedAt, &i.Notes, &i.PositionTitle, &i.CompanyName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (q *Queries) GetApplication(ctx context.Context, id uuid.UUID) (ApplicationWithDetails, error) {
	row := q.db.QueryRow(ctx, `
		SELECT a.id, a.position_id, a.user_id, a.status, a.source, a.applied_at, a.notes, p.title as position_title, c.name as company_name, c.id as company_id
		FROM applications a
		JOIN positions p ON a.position_id = p.id
		JOIN companies c ON p.company_id = c.id
		WHERE a.id = $1`, id)
	var i ApplicationWithDetails
	err := row.Scan(&i.ID, &i.PositionID, &i.UserID, &i.Status, &i.Source, &i.AppliedAt, &i.Notes, &i.PositionTitle, &i.CompanyName, &i.CompanyID)
	if err != nil {
		return i, err
	}
	i.Interviews, _ = q.ListInterviewsByApplication(ctx, i.ID)
	return i, nil
}

type CreateApplicationParams struct {
	PositionID uuid.UUID
	UserID     uuid.UUID
	Status     string
	Source     *string
	Notes      *string
}

func (q *Queries) CreateApplication(ctx context.Context, arg CreateApplicationParams) (Application, error) {
	row := q.db.QueryRow(ctx, "INSERT INTO applications (position_id, user_id, status, source, notes) VALUES ($1, $2, $3, $4, $5) RETURNING id, position_id, user_id, status, source, applied_at, notes",
		arg.PositionID, arg.UserID, arg.Status, arg.Source, arg.Notes)
	var i Application
	err := row.Scan(&i.ID, &i.PositionID, &i.UserID, &i.Status, &i.Source, &i.AppliedAt, &i.Notes)
	return i, err
}

type UpdateApplicationStatusParams struct {
	ID     uuid.UUID
	Status string
}

func (q *Queries) UpdateApplicationStatus(ctx context.Context, arg UpdateApplicationStatusParams) (Application, error) {
	row := q.db.QueryRow(ctx, "UPDATE applications SET status = $2 WHERE id = $1 RETURNING id, position_id, user_id, status, source, applied_at, notes",
		arg.ID, arg.Status)
	var i Application
	err := row.Scan(&i.ID, &i.PositionID, &i.UserID, &i.Status, &i.Source, &i.AppliedAt, &i.Notes)
	return i, err
}

// Interviews
func (q *Queries) ListInterviewsByApplication(ctx context.Context, applicationID uuid.UUID) ([]Interview, error) {
	rows, err := q.db.Query(ctx, "SELECT id, application_id, stage_name, scheduled_at, notes, feedback FROM interviews WHERE application_id = $1 ORDER BY scheduled_at ASC", applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Interview
	for rows.Next() {
		var i Interview
		if err := rows.Scan(&i.ID, &i.ApplicationID, &i.StageName, &i.ScheduledAt, &i.Notes, &i.Feedback); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

type CreateInterviewParams struct {
	ApplicationID uuid.UUID
	StageName     string
	ScheduledAt   time.Time
	Notes         *string
}

func (q *Queries) CreateInterview(ctx context.Context, arg CreateInterviewParams) (Interview, error) {
	row := q.db.QueryRow(ctx, "INSERT INTO interviews (application_id, stage_name, scheduled_at, notes) VALUES ($1, $2, $3, $4) RETURNING id, application_id, stage_name, scheduled_at, notes, feedback",
		arg.ApplicationID, arg.StageName, arg.ScheduledAt, arg.Notes)
	var i Interview
	err := row.Scan(&i.ID, &i.ApplicationID, &i.StageName, &i.ScheduledAt, &i.Notes, &i.Feedback)
	return i, err
}

// Skills
func (q *Queries) ListSkills(ctx context.Context) ([]Skill, error) {
	rows, err := q.db.Query(ctx, "SELECT id, name FROM skills ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Skill
	for rows.Next() {
		var i Skill
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (q *Queries) CreateSkill(ctx context.Context, name string) (Skill, error) {
	row := q.db.QueryRow(ctx, "INSERT INTO skills (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id, name", name)
	var i Skill
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

func (q *Queries) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, "DELETE FROM applications WHERE id = $1", id)
	return err
}

func (q *Queries) DeletePosition(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, "DELETE FROM positions WHERE id = $1", id)
	return err
}

type UpdatePositionParams struct {
	ID          uuid.UUID
	Title       string
	Location    *string
	WorkMode    string
	SalaryRange *string
	PostURL     *string
}

func (q *Queries) UpdatePosition(ctx context.Context, arg UpdatePositionParams) (Position, error) {
	row := q.db.QueryRow(ctx, "UPDATE positions SET title = $2, location = $3, work_mode = $4, salary_range = $5, post_url = $6 WHERE id = $1 RETURNING id, company_id, user_id, title, location, work_mode, salary_range, post_url, created_at",
		arg.ID, arg.Title, arg.Location, arg.WorkMode, arg.SalaryRange, arg.PostURL)
	var i Position
	err := row.Scan(&i.ID, &i.CompanyID, &i.UserID, &i.Title, &i.Location, &i.WorkMode, &i.SalaryRange, &i.PostURL, &i.CreatedAt)
	return i, err
}

func (q *Queries) GetDashboardStats(ctx context.Context, userID uuid.UUID) (DashboardStats, error) {
	var stats DashboardStats
	stats.StatusBreakdown = make(map[string]int)

	// Total Apps
	err := q.db.QueryRow(ctx, "SELECT COUNT(*) FROM applications WHERE user_id = $1", userID).Scan(&stats.TotalApplications)
	if err != nil {
		return stats, err
	}

	// Interviewing
	err = q.db.QueryRow(ctx, "SELECT COUNT(*) FROM applications WHERE user_id = $1 AND status = 'Interviewing'", userID).Scan(&stats.InterviewingCount)
	if err != nil {
		return stats, err
	}

	// Offers
	err = q.db.QueryRow(ctx, "SELECT COUNT(*) FROM applications WHERE user_id = $1 AND status = 'Offer'", userID).Scan(&stats.OfferCount)
	if err != nil {
		return stats, err
	}

	// Status Breakdown
	rows, err := q.db.Query(ctx, "SELECT status, COUNT(*) FROM applications WHERE user_id = $1 GROUP BY status", userID)
	if err != nil {
		return stats, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return stats, err
		}
		stats.StatusBreakdown[status] = count
	}

	return stats, nil
}
