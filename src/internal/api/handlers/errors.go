package handlers

import "errors"

var ErrMsgUserNotExist = errors.New("user not exist")
var ErrMsgNotPermission = errors.New("user doesnt have permission")
var ErrMsgNotFound = errors.New("not found")
var ErrNotValidParam = errors.New("invalid format")
