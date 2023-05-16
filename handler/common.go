package handler

const (
	APIRoute = "/api"
)

const (
	RegisterRoute = "/register"
	LoginRoute    = "/login"
)

type Response struct {
	Message string `json:"message"`
}

var (
	SucceedResponse      = Response{Message: "Succeed"}
	UnauthorizedResponse = Response{Message: "Unauthorized"}
)
