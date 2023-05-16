package handler

const (
	APIRoute = "/api"
)

const (
	UUIDParam = "uuid"
)

const (
	RegisterRoute               = "/register"
	LoginRoute                  = "/login"
	QueryStationRoute           = "/station"
	QueryStationRouteWithParams = QueryStationRoute + "/:" + UUIDParam
)

type Response struct {
	Message string `json:"message"`
}

var (
	SucceedResponse      = Response{Message: "Succeed"}
	UnauthorizedResponse = Response{Message: "Unauthorized"}
)
