-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL COLLATE NOCASE CHECK(username<>'') UNIQUE,
    password TEXT NOT NULL CHECK(password<>''),
    role TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name<>''),

    owner_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE boards (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name<>''),

    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    order_number INT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL CHECK(title<>''),

    -- TODO(patrik): Might change this to nullable
    board_id TEXT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE tags (
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    slug TEXT NOT NULL,

    PRIMARY KEY(project_id, slug)
);

CREATE TABLE tasks_tags  (
    task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    tag_slug TEXT NOT NULL,

    FOREIGN KEY(project_id, tag_slug) REFERENCES tags(project_id, slug),
    PRIMARY KEY(task_id, project_id, tag_slug)
);

CREATE TABLE users_settings (
    id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name TEXT
);

CREATE TABLE api_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

-- +goose Down
DROP TABLE api_tokens; 

DROP TABLE users_settings; 

DROP TABLE tasks_tags; 
DROP TABLE tags; 

DROP TABLE tasks; 

DROP TABLE users; 
