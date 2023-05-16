package handler

const (
	RootRoute = "/"
	APIRoute  = "/api"
)

const (
	SessionVariable = "SESSION_VARIABLE"
)

const (
	UUIDParam = "uuid"
)

const (
	RegisterRoute               = "/register"
	LoginRoute                  = "/login"
	QueryStationRoute           = "/station"
	QueryStationRouteWithParams = QueryStationRoute + "/:" + UUIDParam
	SearchStationsRoute         = "/search/stations"
	HistoricalRoute             = "/historical"
	AccountRoute                = "/account"
	EchoRoute                   = "/echo"
	AlertRoute                  = "/alert"
)

type Response struct {
	Message string `json:"message"`
}

var (
	SucceedResponse      = Response{Message: "Succeed"}
	UnauthorizedResponse = Response{Message: "Unauthorized"}
)
