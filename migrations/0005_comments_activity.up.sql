CREATE TABLE IF NOT EXISTS comments (
    id CHAR(36) PRIMARY KEY,
    workspace_id CHAR(36) NOT NULL,
    project_id CHAR(36) NOT NULL,
    task_id CHAR(36) NOT NULL,
    author_user_id CHAR(36) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_comments_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    CONSTRAINT fk_comments_project_id FOREIGN KEY (project_id) REFERENCES projects(id),
    CONSTRAINT fk_comments_task_id FOREIGN KEY (task_id) REFERENCES tasks(id),
    INDEX index_comments_task_id_created_at (task_id, created_at)
);

CREATE TABLE IF NOT EXISTS activity_logs (
    id CHAR(36) PRIMARY KEY,
    workspace_id CHAR(36) NOT NULL,
    project_id CHAR(36) NOT NULL,
    task_id CHAR(36) NOT NULL,
    actor_user_id CHAR(36) NOT NULL,
    activity_type VARCHAR(128) NOT NULL,
    payload JSON NOT NULL,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_activity_logs_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    CONSTRAINT fk_activity_logs_project_id FOREIGN KEY (project_id) REFERENCES projects(id),
    CONSTRAINT fk_activity_logs_task_id FOREIGN KEY (task_id) REFERENCES tasks(id),
    INDEX index_activity_logs_task_id_created_at (task_id, created_at)
);
