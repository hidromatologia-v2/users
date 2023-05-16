package handler

const (
	APIRoute = "/api"
)

const (
	RegisterRoute = "/register"
)

type Response struct {
	Message string `json:"message"`
}

var (
	SucceedResponse      = Response{Message: "Succeed"}
	UnauthorizedResponse = Response{Message: "Unauthorized"}
)
