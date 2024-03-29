package handler

const (
	RootRoute = "/"
)

const (
	SessionVariable = "SESSION_VARIABLE"
)

const (
	UUIDParam = "uuid"
)

const (
	RegisterRoute          = "/register"
	LoginRoute             = "/login"
	SensorRoute            = "/sensor"
	StationRoute           = "/station"
	StationRouteWithParams = StationRoute + "/:" + UUIDParam
	HistoricalRoute        = "/historical"
	AccountRoute           = "/account"
	EchoRoute              = "/echo"
	AlertRoute             = "/alert"
	AlertRouteWithParam    = AlertRoute + "/:" + UUIDParam
	ResetPasswordRoute     = "/reset/password"
	ConfirmAccountRoute    = "/confirm/account"
)

type Response struct {
	Message string `json:"message"`
}

var (
	SucceedResponse      = Response{Message: "Succeed"}
	UnauthorizedResponse = Response{Message: "Unauthorized"}
)
