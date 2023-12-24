package form

type CreateWorkerRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type ModifyWorkerRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}