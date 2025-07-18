CREATE TABLE IF NOT EXISTS reports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id uuid NOT NULL,
    trainer_id uuid NOT NULL,
    group_id uuid NOT NULL,
    description VARCHAR(500),
    date DATE NOT NULL
)