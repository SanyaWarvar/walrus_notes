package response

import (
	"context"
	"net/http"
	"wn/pkg/constants"
)

type meta struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Error     string `json:"error"`
	RequestId string `json:"requestId"`
}

type pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
	Pages   int `json:"pages"`
}

type Response struct {
	Meta       meta       `json:"meta"`
	Pagination pagination `json:"pagination"`
	Data       any        `json:"data"`
}

type Builder struct {
	ErrorExport bool `json:"errorExport"`
}

func NewResponseBuilder(errorExport bool) *Builder {
	return &Builder{ErrorExport: errorExport}
}

func (rb *Builder) BuildErrorResponse(ctx context.Context, message string, code string, err error) *Response {
	requestId := ""
	if ctx.Value(constants.RequestIdCtx) != nil {
		requestId = ctx.Value(constants.RequestIdCtx).(string)
	}
	m := meta{
		Message:   message,
		RequestId: requestId,
		Code:      code,
	}
	//show error
	if rb.ErrorExport {
		m.Error = err.Error()
	}

	return &Response{Meta: m}
}

func (rb *Builder) BuildSuccessPaginationResponse(ctx context.Context, page, perPage, pages int, obj any) *Response {
	r := buildResponse(ctx, obj)
	r.Pagination = pagination{
		Page:    page,
		PerPage: perPage,
		Pages:   pages,
	}
	return r
}

func (rb *Builder) BuildSuccessResponseBody(ctx context.Context, obj any) (int, *Response) {
	return http.StatusOK, buildResponse(ctx, obj)
}

func (rb *Builder) BuildSuccessResponse(ctx context.Context) (int, *Response) {
	return http.StatusOK, buildResponse(ctx, nil)
}

func buildResponse(ctx context.Context, obj any) *Response {
	requestId := ""

	if ctx.Value(constants.RequestIdCtx) != nil {
		requestId = ctx.Value(constants.RequestIdCtx).(string)
	}

	m := meta{
		Message:   "OK",
		RequestId: requestId,
	}
	return &Response{
		Meta: m,
		Data: obj,
	}
}
