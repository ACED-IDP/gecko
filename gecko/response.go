package gecko

import (
	"encoding/json"
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/uc-cdis/arborist/arborist"
)

type jsonResponse struct {
	content any
	code    int
}

type ErrorResponse struct {
	HTTPError arborist.HTTPError `json:"error"`
	// err stores an internal representation of an error in case it needs to be
	// tracked along with the http-ish version in `HTTPError`.
	err error
	// log embeds a LogCache so we can log things to the response and write it
	// out later to the server's logger.
	log LogCache
}

func jsonResponseFrom(content any, code int) *jsonResponse {
	return &jsonResponse{
		content: content,
		code:    code,
	}
}

func (response *jsonResponse) write(ctx iris.Context) error {
	ctx.ContentType("application/json")
	if response.code > 0 {
		ctx.StatusCode(response.code)
	} else {
		ctx.StatusCode(http.StatusOK)
	}
	var bytes []byte
	var err error
	if wantPrettyJSON(ctx.Request()) {
		bytes, err = json.MarshalIndent(response.content, "", "    ")
	} else {
		bytes, err = json.Marshal(response.content)
	}
	if err != nil {
		return err
	}
	_, err = ctx.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func wantPrettyJSON(r *http.Request) bool {
	prettyJSON := false
	if r.Method == "GET" {
		prettyJSON = prettyJSON || r.URL.Query().Get("pretty") == "true"
		prettyJSON = prettyJSON || r.URL.Query().Get("prettyJSON") == "true"
	}
	return prettyJSON
}

func newErrorResponse(message string, code int, err *error) *ErrorResponse {
	response := &ErrorResponse{
		HTTPError: arborist.HTTPError{
			Message: message,
			Code:    code,
		},
	}
	if err != nil {
		response.err = *err
	}
	if code >= 500 {
		response.log.Error(message)
	} else {
		response.log.Info(message)
	}
	return response
}

func (errorResponse *ErrorResponse) write(ctx iris.Context) error {
	var bytes []byte
	var err error

	prettyJSON := false
	if ctx.Method() == "GET" {
		prettyJSON = prettyJSON || ctx.URLParamDefault("pretty", "false") == "true"
		prettyJSON = prettyJSON || ctx.URLParamDefault("pretty", "false") == "true"
	}

	if prettyJSON {
		bytes, err = json.MarshalIndent(errorResponse, "", "    ")
	} else {
		bytes, err = json.Marshal(errorResponse)
	}
	if err != nil {
		return err
	}
	ctx.ContentType("application/json")
	ctx.StatusCode(errorResponse.HTTPError.Code)
	_, err = ctx.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
