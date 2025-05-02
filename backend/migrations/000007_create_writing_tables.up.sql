CREATE TABLE IF NOT EXISTS writing_variants (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    version integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS writing_tasks (
    id bigserial PRIMARY KEY,
    variant_id bigint NOT NULL REFERENCES writing_variants(id) ON DELETE CASCADE,
    task_number smallint NOT NULL CHECK (task_number IN (1, 2)),
    prompt text NOT NULL,
    image_url text,
    version integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS writing_submissions (
    id bigserial PRIMARY KEY,
    task_id bigint NOT NULL REFERENCES writing_tasks(id) ON DELETE CASCADE,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    answer text NOT NULL,
    submitted_at timestamp with time zone DEFAULT now(),
    version integer NOT NULL DEFAULT 1
);
