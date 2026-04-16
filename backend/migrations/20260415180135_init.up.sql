CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       full_name TEXT NOT NULL,
                       login TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       role TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE reports (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         user_id UUID REFERENCES users(id),
                         place TEXT NOT NULL,
                         report_date DATE NOT NULL,
                         responsible_name TEXT NOT NULL,
                         created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE questions (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           text TEXT NOT NULL,
                           order_index INT NOT NULL,
                           is_active BOOLEAN DEFAULT true,
                           created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE answers (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         report_id UUID REFERENCES reports(id) ON DELETE CASCADE,
                         question_id UUID REFERENCES questions(id),
                         answer_text TEXT NOT NULL,
                         image_url TEXT,
                         created_at TIMESTAMP DEFAULT now()
);

INSERT INTO users (full_name, login, password_hash, role)
VALUES ('Admin', 'admin', '$2a$10$ODPCO0OS3NsI/YxLkUaRtuRYQnNXXSlz3VwWmF2JNdy/oZ./dpzXe', 'admin');