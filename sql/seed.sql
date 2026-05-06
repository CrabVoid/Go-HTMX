-- 1. Create a default user (password: 'password123')
-- User ID: 00000000-0000-0000-0000-000000000001
INSERT INTO users (id, email, password_hash) 
VALUES ('00000000-0000-0000-0000-000000000001', 'test@example.com', '$2a$10$8K1p/aC2l9uP2vJ0H5Y6Oe5o0Y9v7V.nF8f9S9zWvX7kY8h9u/m1e')
ON CONFLICT (email) DO NOTHING;

-- 2. Insert Companies
INSERT INTO companies (id, name, website, industry) VALUES
('11111111-1111-1111-1111-111111111111', 'Tech Corp', 'https://techcorp.com', 'Technology'),
('22222222-2222-2222-2222-222222222222', 'Finance Solutions', 'https://financesolutions.com', 'Finance'),
('33333333-3333-3333-3333-333333333333', 'Health Plus', 'https://healthplus.com', 'Healthcare'),
('deadbeef-1111-1111-1111-111111111111', 'Cloud Systems', 'https://cloudsystems.io', 'Technology'),
('deadbeef-2222-2222-2222-222222222222', 'Retail Giant', 'https://retailgiant.com', 'Retail')
ON CONFLICT (name) DO NOTHING;

-- 3. Insert Skills
INSERT INTO skills (id, name) VALUES
('44444444-4444-4444-4444-444444444444', 'Go'),
('55555555-5555-5555-5555-555555555555', 'PostgreSQL'),
('66666666-6666-6666-6666-666666666666', 'React'),
('77777777-7777-7777-7777-777777777777', 'Docker'),
('deadbeef-4444-4444-4444-444444444444', 'Kubernetes'),
('deadbeef-5555-5555-5555-555555555555', 'Python'),
('deadbeef-6666-6666-6666-666666666666', 'TypeScript')
ON CONFLICT (name) DO NOTHING;

-- 4. Insert Positions
INSERT INTO positions (id, company_id, user_id, title, location, work_mode, salary_range, post_url) VALUES
('88888888-8888-8888-8888-888888888888', '11111111-1111-1111-1111-111111111111', '00000000-0000-0000-0000-000000000001', 'Backend Developer Intern', 'Remote', 'Remote', '$3000 - $4000', 'https://techcorp.com/jobs/1'),
('99999999-9999-9999-9999-999999999999', '22222222-2222-2222-2222-222222222222', '00000000-0000-0000-0000-000000000001', 'Fullstack Engineer', 'New York, NY', 'Hybrid', '$100k - $120k', 'https://financesolutions.com/careers/2'),
('beefbeef-1111-1111-1111-111111111111', 'deadbeef-1111-1111-1111-111111111111', '00000000-0000-0000-0000-000000000001', 'DevOps Engineer', 'Austin, TX', 'Hybrid', '$110k - $140k', 'https://cloudsystems.io/jobs/devops'),
('beefbeef-2222-2222-2222-222222222222', '33333333-3333-3333-3333-333333333333', '00000000-0000-0000-0000-000000000001', 'Frontend Developer', 'London, UK', 'Onsite', '£50k - £70k', 'https://healthplus.com/careers/frontend'),
('beefbeef-3333-3333-3333-333333333333', 'deadbeef-2222-2222-2222-222222222222', '00000000-0000-0000-0000-000000000001', 'Data Scientist', 'Remote', 'Remote', '$130k - $160k', 'https://retailgiant.com/jobs/data-science')
ON CONFLICT (id) DO NOTHING;

-- 5. Link Positions to Skills
INSERT INTO position_skills (position_id, skill_id) VALUES
('88888888-8888-8888-8888-888888888888', '44444444-4444-4444-4444-444444444444'), -- Backend Intern -> Go
('88888888-8888-8888-8888-888888888888', '55555555-5555-5555-5555-555555555555'), -- Backend Intern -> PG
('99999999-9999-9999-9999-999999999999', '66666666-6666-6666-6666-666666666666'), -- Fullstack -> React
('99999999-9999-9999-9999-999999999999', '44444444-4444-4444-4444-444444444444'), -- Fullstack -> Go
('beefbeef-1111-1111-1111-111111111111', '77777777-7777-7777-7777-777777777777'), -- DevOps -> Docker
('beefbeef-1111-1111-1111-111111111111', 'deadbeef-4444-4444-4444-444444444444'), -- DevOps -> K8s
('beefbeef-2222-2222-2222-222222222222', '66666666-6666-6666-6666-666666666666'), -- Frontend -> React
('beefbeef-2222-2222-2222-222222222222', 'deadbeef-6666-6666-6666-666666666666'), -- Frontend -> TS
('beefbeef-3333-3333-3333-333333333333', 'deadbeef-5555-5555-5555-555555555555')  -- Data Scientist -> Python
ON CONFLICT DO NOTHING;

-- 6. Insert Applications
INSERT INTO applications (id, position_id, user_id, status, source, notes) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '88888888-8888-8888-8888-888888888888', '00000000-0000-0000-0000-000000000001', 'Interviewing', 'LinkedIn', 'First contact via recruiter'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '99999999-9999-9999-9999-999999999999', '00000000-0000-0000-0000-000000000001', 'Applied', 'Referral', 'Referred by a friend')
ON CONFLICT (id) DO NOTHING;

-- 7. Insert Interviews
INSERT INTO interviews (id, application_id, stage_name, scheduled_at, notes) VALUES
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Technical Round', NOW() + INTERVAL '2 days', 'Coding challenge on Go fundamentals'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'System Design', NOW() + INTERVAL '5 days', 'Architecture discussion')
ON CONFLICT (id) DO NOTHING;
