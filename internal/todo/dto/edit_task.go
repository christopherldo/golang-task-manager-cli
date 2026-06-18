package dto

type EditTaskDto struct {
	Description string `json:"description"`
	IsDone      *bool  `json:"isDone,omitempty"`
}
