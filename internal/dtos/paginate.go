package dtos

const defaultLimit = 10

type PaginatedRequest struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"-" form:"limit"`
}

func (r *PaginatedRequest) Normalize() {
	if r.Page <= 0 {
		r.Page = 1
	}

	if r.Limit <= 0 {
		r.Limit = defaultLimit
	}
}
