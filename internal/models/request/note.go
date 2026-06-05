package request

type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type PatchNoteRequest struct {
	Title   string `json:"title" validate:"omitempty"`
	Content string `json:"content" validate:"omitempty"`
}
