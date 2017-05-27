// Code generated by "stringer -type ItemAction,TaskStatus"; DO NOT EDIT.

package thingscloud

import "fmt"

const _ItemAction_name = "ActionCreatedActionModifiedActionDeleted"

var _ItemAction_index = [...]uint8{0, 13, 27, 40}

func (i ItemAction) String() string {
	if i < 0 || i >= ItemAction(len(_ItemAction_index)-1) {
		return fmt.Sprintf("ItemAction(%d)", i)
	}
	return _ItemAction_name[_ItemAction_index[i]:_ItemAction_index[i+1]]
}

const (
	_TaskStatus_name_0 = "TaskStatusPending"
	_TaskStatus_name_1 = "TaskStatusCanceledTaskStatusCompleted"
)

var (
	_TaskStatus_index_0 = [...]uint8{0, 17}
	_TaskStatus_index_1 = [...]uint8{0, 18, 37}
)

func (i TaskStatus) String() string {
	switch {
	case i == 0:
		return _TaskStatus_name_0
	case 2 <= i && i <= 3:
		i -= 2
		return _TaskStatus_name_1[_TaskStatus_index_1[i]:_TaskStatus_index_1[i+1]]
	default:
		return fmt.Sprintf("TaskStatus(%d)", i)
	}
}