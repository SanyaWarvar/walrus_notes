package constants

// Headers
const (
	AuthorizationHeader = "Authorization"
	RequestIdHeader     = "X-REQUEST-ID"
	RefreshHeader       = "X-REFRESH-TOKEN"
)

// Roles
const (
	ClientRole = "CLIENT"
	AdminRole  = "ADMIN"
)

var RolesToIntMap = map[string]int{
	ClientRole: 1,
	AdminRole:  2,
}

// Context
const (
	UserIdCtx    = "userId"
	UserRoleCtx  = "userRole"
	RequestIdCtx = "requestId"
	TraceIdCtx   = "traceId"
	SpanIdCtx    = "spanId"
	ApiNameCtx   = "apiName"
)

// Errors
const (
	BindBodyError string = "bind_body"
	BindPathError string = "bind_path"
)
