package model

type NavigateMsg struct {
	To         Page
	ActiveTask string
}

type TaskSavedMsg struct {
	Task interface{}
}

type ErrMsg struct {
	Err error
}

func (e ErrMsg) Error() string {
	return e.Err.Error()
}
