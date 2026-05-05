-- name: ListCompanies :many
SELECT * FROM companies ORDER BY name;

-- name: CreateCompany :one
INSERT INTO companies (name, website, industry)
VALUES ($1, $2, $3)
RETURNING *;
