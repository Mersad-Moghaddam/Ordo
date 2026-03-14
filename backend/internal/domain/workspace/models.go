package workspace

import "time"

type MembershipRole string

const (
	MembershipRoleOwner  MembershipRole = "owner"
	MembershipRoleAdmin  MembershipRole = "admin"
	MembershipRoleMember MembershipRole = "member"
)

type Workspace struct {
	WorkspaceID  string
	WorkspaceKey string
	DisplayName  string
	CreatedByID  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type WorkspaceMembership struct {
	WorkspaceID     string
	UserID          string
	MembershipRole  MembershipRole
	InvitedByUserID string
	JoinedAt        time.Time
	LastUpdatedAt   time.Time
}

type Project struct {
	ProjectID      string
	WorkspaceID    string
	ProjectKey     string
	DisplayName    string
	Description    string
	CreatedByID    string
	CreatedAt      time.Time
	LastUpdatedAt  time.Time
	ArchivedAtTime *time.Time
}
