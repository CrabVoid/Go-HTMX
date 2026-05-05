-- Companies
-- name: ListCompanies :many
SELECT * FROM companies ORDER BY name;

-- name: GetCompany :one
SELECT * FROM companies WHERE id = $1;

-- name: CreateCompany :one
INSERT INTO companies (name, website, industry)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCompany :one
UPDATE companies 
SET name = $2, website = $3, industry = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCompany :exec
DELETE FROM companies WHERE id = $1;

-- Positions
-- name: ListPositions :many
SELECT p.*, c.name as company_name 
FROM positions p
JOIN companies c ON p.company_id = c.id
ORDER BY p.created_at DESC;

-- name: GetPosition :one
SELECT p.*, c.name as company_name 
FROM positions p
JOIN companies c ON p.company_id = c.id
WHERE p.id = $1;

-- name: CreatePosition :one
INSERT INTO positions (company_id, title, location, work_mode, salary_range, post_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdatePosition :one
UPDATE positions
SET title = $2, location = $3, work_mode = $4, salary_range = $5, post_url = $6
WHERE id = $1
RETURNING *;

-- name: DeletePosition :exec
DELETE FROM positions WHERE id = $1;

-- Applications
-- name: ListApplications :many
SELECT a.*, p.title as position_title, c.name as company_name
FROM applications a
JOIN positions p ON a.position_id = p.id
JOIN companies c ON p.company_id = c.id
ORDER BY a.applied_at DESC;

-- name: GetApplication :one
SELECT a.*, p.title as position_title, c.name as company_name, c.id as company_id
FROM applications a
JOIN positions p ON a.position_id = p.id
JOIN companies c ON p.company_id = c.id
WHERE a.id = $1;

-- name: CreateApplication :one
INSERT INTO applications (position_id, status, source, notes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateApplicationStatus :one
UPDATE applications SET status = $2 WHERE id = $1 RETURNING *;

-- name: DeleteApplication :exec
DELETE FROM applications WHERE id = $1;

-- Interviews
-- name: ListInterviewsByApplication :many
SELECT * FROM interviews WHERE application_id = $1 ORDER BY scheduled_at ASC;

-- name: CreateInterview :one
INSERT INTO interviews (application_id, stage_name, scheduled_at, notes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- Skills
-- name: ListSkills :many
SELECT * FROM skills ORDER BY name;

-- name: CreateSkill :one
INSERT INTO skills (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING *;
