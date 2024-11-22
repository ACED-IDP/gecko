package chameleon

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris"
	"github.com/uc-cdis/arborist/arborist"
)

type LogHandler struct {
	logger *log.Logger
}

type Server struct {
	iris   *iris.Application
	db     *sqlx.DB
	jwtApp arborist.JWTDecoder
	logger *LogHandler
	stmts  *arborist.CachedStmts
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) WithLogger(logger *log.Logger) *Server {
	server.logger = &LogHandler{logger: logger}
	return server
}

func (server *Server) WithJWTApp(jwtApp arborist.JWTDecoder) *Server {
	server.jwtApp = jwtApp
	return server
}

func (server *Server) WithDB(db *sqlx.DB) *Server {
	server.db = db
	server.stmts = arborist.NewCachedStmts(db)
	return server
}

func (server *Server) Init() (*Server, error) {
	if server.db == nil {
		return nil, errors.New("chameleon server initialized without database")
	}
	if server.jwtApp == nil {
		return nil, errors.New("chameleon server initialized without JWT app")
	}
	if server.logger == nil {
		return nil, errors.New("chameleon server initialized without logger")
	}

	return server, nil
}

func (server *Server) MakeRouter() *iris.Application {
	router := iris.New()
	router.Get("/health", server.handleHealth)
	router.Get("/config/{configId}", server.handleConfigGET)
	router.Put("/config/{configId}", server.handleConfigPUT)
	router.OnErrorCode(iris.StatusNotFound, handleNotFound)
	router.UseRouter(func(ctx iris.Context) {
		ctx.Request().URL.Path = strings.TrimSuffix(ctx.Request().URL.Path, "/")
		ctx.Next()
	})
	return router
}

func (server *Server) handleConfigGET(ctx iris.Context) {
	configId := ctx.Params().Get("configId")
	doc, err := configGET(server.db, configId)
	if doc == nil {
		msg := fmt.Sprintf("no client found with clientID: %s", configId)
		errResponse := newErrorResponse(msg, 404, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}
	if err != nil {
		msg := fmt.Sprintf("client query failed: %s", err.Error())
		errResponse := newErrorResponse(msg, 500, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}

	_ = jsonResponseFrom(doc, http.StatusOK).write(ctx)
}

func (server *Server) handleConfigPUT(ctx iris.Context) {
	configId := ctx.Params().Get("configId")
	data := map[string]any{}
	body := ctx.Recorder().Body()
	errResponse := unmarshal(body, data)
	if errResponse != nil {
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}
	err := configPUT(server.db, configId, data)
	if err != nil {
		msg := fmt.Sprintf("client query failed: %s", err.Error())
		errResponse := newErrorResponse(msg, 500, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}

	//_ = jsonResponseFrom(policy, http.StatusOK).write(w, r)
}

func (server *Server) handleHealth(ctx iris.Context) {
	err := server.db.Ping()
	if err != nil {
		server.logger.Error("database ping failed; returning unhealthy")
		response := newErrorResponse("database unavailable", 500, nil)
		_ = response.write(ctx)
		return
	}
	_ = jsonResponseFrom("Healthy", http.StatusOK).write(ctx)
}

func handleNotFound(ctx iris.Context) {
	response := struct {
		Error struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"error"`
	}{
		Error: struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "not found",
			Code:    404,
		},
	}
	_ = jsonResponseFrom(response, 404).write(ctx)
}

func unmarshal(body []byte, x interface{}) *ErrorResponse {
	var structValue reflect.Value = reflect.ValueOf(x)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}
	var structType reflect.Type = structValue.Type()
	err := json.Unmarshal(body, x)
	if err != nil {
		msg := fmt.Sprintf(
			"could not parse %s from JSON; make sure input has correct types",
			structType,
		)
		response := newErrorResponse(msg, 400, &err)
		response.log.Info(
			"tried to create %s but input was invalid; offending JSON: %s",
			structType,
			loggableJSON(body),
		)
		return response
	}
	return nil
}

func loggableJSON(bytes []byte) []byte {
	return regWhitespace.ReplaceAll(bytes, []byte(""))
}

var regWhitespace *regexp.Regexp = regexp.MustCompile(`\s`)
