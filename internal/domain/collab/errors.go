package collab

import "errors"

var ErrCommentNotFound = errors.New("comment not found")
var ErrCommentForbidden = errors.New("comment forbidden")
var ErrActivityLogWriteFailure = errors.New("activity log write failure")
var ErrCommentAlreadyDeleted = errors.New("comment already deleted")
