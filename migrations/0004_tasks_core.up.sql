CREATE TABLE IF NOT EXISTS tasks (
    id CHAR(36) PRIMARY KEY,
    workspace_id CHAR(36) NOT NULL,
    project_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(32) NOT NULL,
    priority VARCHAR(16) NOT NULL,
    assignee_user_id CHAR(36) NULL,
    created_by_user_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_tasks_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    CONSTRAINT fk_tasks_project_id FOREIGN KEY (project_id) REFERENCES projects(id),
    INDEX index_tasks_project_id (project_id),
    INDEX index_tasks_workspace_id (workspace_id),
    INDEX index_tasks_status_priority (status, priority)
);
