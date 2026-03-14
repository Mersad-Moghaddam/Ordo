package workspace

import "errors"

var ErrWorkspaceNotFound = errors.New("workspace not found")
var ErrWorkspaceAlreadyExists = errors.New("workspace already exists")
var ErrMembershipNotFound = errors.New("membership not found")
var ErrMembershipAlreadyExists = errors.New("membership already exists")
var ErrProjectNotFound = errors.New("project not found")
var ErrProjectAlreadyExists = errors.New("project already exists")
var ErrInsufficientWorkspaceRole = errors.New("insufficient workspace role")
