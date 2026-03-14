package memory

import (
	"context"
	"sync"
	"time"

	domainauth "github.com/ordo/backend/internal/domain/auth"
	domaincollab "github.com/ordo/backend/internal/domain/collab"
	domaintask "github.com/ordo/backend/internal/domain/task"
	domainworkspace "github.com/ordo/backend/internal/domain/workspace"
)

type Store struct {
	storeMutex            sync.Mutex
	usersByID             map[string]domainauth.User
	usersByEmail          map[string]string
	refreshSessionsByID   map[string]domainauth.RefreshSession
	workspacesByID        map[string]domainworkspace.Workspace
	workspacesByKey       map[string]string
	membershipsByKey      map[string]domainworkspace.WorkspaceMembership
	projectsByID          map[string]domainworkspace.Project
	projectIDsByComposite map[string]string
	tasksByID             map[string]domaintask.Task
	commentsByID          map[string]domaincollab.Comment
	activitiesByTaskID    map[string][]domaincollab.ActivityLog
	outboxByID            map[string]domaintask.OutboxEvent
}

func NewStore() *Store {
	return &Store{
		usersByID:             map[string]domainauth.User{},
		usersByEmail:          map[string]string{},
		refreshSessionsByID:   map[string]domainauth.RefreshSession{},
		workspacesByID:        map[string]domainworkspace.Workspace{},
		workspacesByKey:       map[string]string{},
		membershipsByKey:      map[string]domainworkspace.WorkspaceMembership{},
		projectsByID:          map[string]domainworkspace.Project{},
		projectIDsByComposite: map[string]string{},
		tasksByID:             map[string]domaintask.Task{},
		commentsByID:          map[string]domaincollab.Comment{},
		activitiesByTaskID:    map[string][]domaincollab.ActivityLog{},
		outboxByID:            map[string]domaintask.OutboxEvent{},
	}
}

func membershipCompositeKey(workspaceID string, userID string) string {
	return workspaceID + ":" + userID
}

func projectCompositeKey(workspaceID string, projectKey string) string {
	return workspaceID + ":" + projectKey
}

func (store *Store) WithTransaction(requestContext context.Context, transactionWorkload func(transactionContext context.Context) error) error {
	return transactionWorkload(requestContext)
}

func (store *Store) FindByEmailAddress(requestContext context.Context, emailAddress string) (domainauth.User, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	userID, hasUser := store.usersByEmail[emailAddress]
	if !hasUser {
		return domainauth.User{}, domainauth.ErrUserNotFound
	}
	return store.usersByID[userID], nil
}

func (store *Store) FindByUserID(requestContext context.Context, userID string) (domainauth.User, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	user, hasUser := store.usersByID[userID]
	if !hasUser {
		return domainauth.User{}, domainauth.ErrUserNotFound
	}
	return user, nil
}

func (store *Store) CreateUser(requestContext context.Context, user domainauth.User) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	if _, hasUser := store.usersByEmail[user.EmailAddress]; hasUser {
		return domainauth.ErrEmailAlreadyRegistered
	}
	store.usersByID[user.UserID] = user
	store.usersByEmail[user.EmailAddress] = user.UserID
	return nil
}

func (store *Store) CreateSession(requestContext context.Context, refreshSession domainauth.RefreshSession) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	store.refreshSessionsByID[refreshSession.SessionID] = refreshSession
	return nil
}

func (store *Store) FindBySessionID(requestContext context.Context, sessionID string) (domainauth.RefreshSession, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	refreshSession, hasSession := store.refreshSessionsByID[sessionID]
	if !hasSession {
		return domainauth.RefreshSession{}, domainauth.ErrSessionNotFound
	}
	return refreshSession, nil
}

func (store *Store) RevokeSession(requestContext context.Context, sessionID string, revokedAtTime time.Time, replacementSessionID *string) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	refreshSession, hasSession := store.refreshSessionsByID[sessionID]
	if !hasSession {
		return domainauth.ErrSessionNotFound
	}
	refreshSession.RevokedAtTime = &revokedAtTime
	refreshSession.ReplacementSessionID = replacementSessionID
	store.refreshSessionsByID[sessionID] = refreshSession
	return nil
}

func (store *Store) CreateWorkspace(requestContext context.Context, workspace domainworkspace.Workspace) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	if _, hasWorkspace := store.workspacesByKey[workspace.WorkspaceKey]; hasWorkspace {
		return domainworkspace.ErrWorkspaceAlreadyExists
	}
	store.workspacesByID[workspace.WorkspaceID] = workspace
	store.workspacesByKey[workspace.WorkspaceKey] = workspace.WorkspaceID
	return nil
}

func (store *Store) FindWorkspaceByWorkspaceID(requestContext context.Context, workspaceID string) (domainworkspace.Workspace, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	workspace, hasWorkspace := store.workspacesByID[workspaceID]
	if !hasWorkspace {
		return domainworkspace.Workspace{}, domainworkspace.ErrWorkspaceNotFound
	}
	return workspace, nil
}

func (store *Store) FindWorkspaceByWorkspaceKey(requestContext context.Context, workspaceKey string) (domainworkspace.Workspace, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	workspaceID, hasWorkspace := store.workspacesByKey[workspaceKey]
	if !hasWorkspace {
		return domainworkspace.Workspace{}, domainworkspace.ErrWorkspaceNotFound
	}
	return store.workspacesByID[workspaceID], nil
}

func (store *Store) ListWorkspacesByUserID(requestContext context.Context, userID string, pageNumber int, pageSize int) ([]domainworkspace.Workspace, int64, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	workspaceList := make([]domainworkspace.Workspace, 0)
	for _, membership := range store.membershipsByKey {
		if membership.UserID == userID {
			workspaceList = append(workspaceList, store.workspacesByID[membership.WorkspaceID])
		}
	}
	return workspaceList, int64(len(workspaceList)), nil
}

func (store *Store) CreateMembership(requestContext context.Context, membership domainworkspace.WorkspaceMembership) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	compositeKey := membershipCompositeKey(membership.WorkspaceID, membership.UserID)
	if _, hasMembership := store.membershipsByKey[compositeKey]; hasMembership {
		return domainworkspace.ErrMembershipAlreadyExists
	}
	store.membershipsByKey[compositeKey] = membership
	return nil
}

func (store *Store) FindMembership(requestContext context.Context, workspaceID string, userID string) (domainworkspace.WorkspaceMembership, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	membership, hasMembership := store.membershipsByKey[membershipCompositeKey(workspaceID, userID)]
	if !hasMembership {
		return domainworkspace.WorkspaceMembership{}, domainworkspace.ErrMembershipNotFound
	}
	return membership, nil
}

func (store *Store) UpdateMembershipRole(requestContext context.Context, workspaceID string, userID string, membershipRole domainworkspace.MembershipRole) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	compositeKey := membershipCompositeKey(workspaceID, userID)
	membership, hasMembership := store.membershipsByKey[compositeKey]
	if !hasMembership {
		return domainworkspace.ErrMembershipNotFound
	}
	membership.MembershipRole = membershipRole
	membership.LastUpdatedAt = time.Now()
	store.membershipsByKey[compositeKey] = membership
	return nil
}

func (store *Store) CreateProject(requestContext context.Context, project domainworkspace.Project) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	compositeKey := projectCompositeKey(project.WorkspaceID, project.ProjectKey)
	if _, hasProject := store.projectIDsByComposite[compositeKey]; hasProject {
		return domainworkspace.ErrProjectAlreadyExists
	}
	store.projectsByID[project.ProjectID] = project
	store.projectIDsByComposite[compositeKey] = project.ProjectID
	return nil
}

func (store *Store) FindProjectByProjectID(requestContext context.Context, projectID string) (domainworkspace.Project, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	project, hasProject := store.projectsByID[projectID]
	if !hasProject {
		return domainworkspace.Project{}, domainworkspace.ErrProjectNotFound
	}
	return project, nil
}

func (store *Store) FindProjectByWorkspaceAndProjectKey(requestContext context.Context, workspaceID string, projectKey string) (domainworkspace.Project, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	projectID, hasProject := store.projectIDsByComposite[projectCompositeKey(workspaceID, projectKey)]
	if !hasProject {
		return domainworkspace.Project{}, domainworkspace.ErrProjectNotFound
	}
	return store.projectsByID[projectID], nil
}

func (store *Store) ListProjectsByWorkspaceID(requestContext context.Context, workspaceID string, pageNumber int, pageSize int) ([]domainworkspace.Project, int64, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	projectList := make([]domainworkspace.Project, 0)
	for _, project := range store.projectsByID {
		if project.WorkspaceID == workspaceID {
			projectList = append(projectList, project)
		}
	}
	return projectList, int64(len(projectList)), nil
}

func (store *Store) CreateTask(requestContext context.Context, task domaintask.Task) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	if _, hasTask := store.tasksByID[task.TaskID]; hasTask {
		return domaintask.ErrTaskAlreadyExists
	}
	store.tasksByID[task.TaskID] = task
	return nil
}

func (store *Store) FindTaskByTaskID(requestContext context.Context, taskID string) (domaintask.Task, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	task, hasTask := store.tasksByID[taskID]
	if !hasTask {
		return domaintask.Task{}, domaintask.ErrTaskNotFound
	}
	return task, nil
}

func (store *Store) UpdateTaskStatus(requestContext context.Context, taskID string, taskStatus domaintask.TaskStatus) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	task, hasTask := store.tasksByID[taskID]
	if !hasTask {
		return domaintask.ErrTaskNotFound
	}
	task.Status = taskStatus
	task.UpdatedAt = time.Now()
	store.tasksByID[taskID] = task
	return nil
}

func (store *Store) ListTasksByProjectID(requestContext context.Context, projectID string, pageNumber int, pageSize int) ([]domaintask.Task, int64, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	taskList := make([]domaintask.Task, 0)
	for _, task := range store.tasksByID {
		if task.ProjectID == projectID {
			taskList = append(taskList, task)
		}
	}
	return taskList, int64(len(taskList)), nil
}

func (store *Store) CreateComment(requestContext context.Context, comment domaincollab.Comment) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	store.commentsByID[comment.CommentID] = comment
	return nil
}

func (store *Store) FindCommentByCommentID(requestContext context.Context, commentID string) (domaincollab.Comment, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	comment, hasComment := store.commentsByID[commentID]
	if !hasComment {
		return domaincollab.Comment{}, domaincollab.ErrCommentNotFound
	}
	return comment, nil
}

func (store *Store) UpdateCommentBody(requestContext context.Context, commentID string, body string) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	comment, hasComment := store.commentsByID[commentID]
	if !hasComment {
		return domaincollab.ErrCommentNotFound
	}
	comment.Body = body
	comment.UpdatedAt = time.Now()
	store.commentsByID[commentID] = comment
	return nil
}

func (store *Store) SoftDeleteComment(requestContext context.Context, commentID string) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	comment, hasComment := store.commentsByID[commentID]
	if !hasComment {
		return domaincollab.ErrCommentNotFound
	}
	deletionTime := time.Now()
	comment.DeletedAt = &deletionTime
	comment.UpdatedAt = deletionTime
	store.commentsByID[commentID] = comment
	return nil
}

func (store *Store) ListCommentsByTaskID(requestContext context.Context, taskID string, pageNumber int, pageSize int) ([]domaincollab.Comment, int64, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	commentList := make([]domaincollab.Comment, 0)
	for _, comment := range store.commentsByID {
		if comment.TaskID == taskID {
			commentList = append(commentList, comment)
		}
	}
	return commentList, int64(len(commentList)), nil
}

func (store *Store) CreateActivity(requestContext context.Context, activityLog domaincollab.ActivityLog) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	store.activitiesByTaskID[activityLog.TaskID] = append(store.activitiesByTaskID[activityLog.TaskID], activityLog)
	return nil
}

func (store *Store) ListActivitiesByTaskID(requestContext context.Context, taskID string, pageNumber int, pageSize int) ([]domaincollab.ActivityLog, int64, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	activityList := store.activitiesByTaskID[taskID]
	return activityList, int64(len(activityList)), nil
}

func (store *Store) CreateOutboxEvent(requestContext context.Context, event domaintask.OutboxEvent) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	store.outboxByID[event.EventID] = event
	return nil
}

func (store *Store) ListPendingOutboxEvents(requestContext context.Context, batchSize int) ([]domaintask.OutboxEvent, error) {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	outboxList := make([]domaintask.OutboxEvent, 0)
	for _, event := range store.outboxByID {
		if event.Status == "pending" || event.Status == "retry" {
			outboxList = append(outboxList, event)
		}
	}
	if batchSize > 0 && len(outboxList) > batchSize {
		return outboxList[:batchSize], nil
	}
	return outboxList, nil
}

func (store *Store) MarkOutboxEventPublished(requestContext context.Context, eventID string) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	event, hasEvent := store.outboxByID[eventID]
	if !hasEvent {
		return nil
	}
	event.Status = "published"
	event.UpdatedAt = time.Now()
	store.outboxByID[eventID] = event
	return nil
}

func (store *Store) MarkOutboxEventRetry(requestContext context.Context, eventID string, attempts int, nextRetryUnixTimestamp int64) error {
	store.storeMutex.Lock()
	defer store.storeMutex.Unlock()
	event, hasEvent := store.outboxByID[eventID]
	if !hasEvent {
		return nil
	}
	event.Status = "retry"
	event.Attempts = attempts
	nextRetryTime := time.Unix(nextRetryUnixTimestamp, 0)
	event.NextRetryAt = &nextRetryTime
	event.UpdatedAt = time.Now()
	store.outboxByID[eventID] = event
	return nil
}
