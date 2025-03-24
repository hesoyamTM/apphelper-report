CREATE TABLE IF NOT EXISTS reports (
    id SERIAL,
    student_id INT NOT NULL,
    trainer_id INT NOT NULL,
    description VARCHAR(500),
    date DATE NOT NULL
)