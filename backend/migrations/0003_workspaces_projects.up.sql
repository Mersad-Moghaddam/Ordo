CREATE TABLE IF NOT EXISTS workspaces (
    id CHAR(36) PRIMARY KEY,
    workspace_key VARCHAR(64) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    created_by_user_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS workspace_memberships (
    workspace_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    membership_role VARCHAR(32) NOT NULL,
    invited_by_user_id CHAR(36) NOT NULL,
    joined_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (workspace_id, user_id),
    CONSTRAINT fk_workspace_memberships_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    CONSTRAINT fk_workspace_memberships_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX index_workspace_memberships_user_id (user_id)
);

CREATE TABLE IF NOT EXISTS projects (
    id CHAR(36) PRIMARY KEY,
    workspace_id CHAR(36) NOT NULL,
    project_key VARCHAR(64) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_by_user_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    archived_at TIMESTAMP NULL,
    UNIQUE KEY unique_workspace_project_key (workspace_id, project_key),
    CONSTRAINT fk_projects_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    INDEX index_projects_workspace_id (workspace_id)
);
