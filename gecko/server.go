package gecko

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/ACED-IDP/gecko/gecko/config"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
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
		return nil, errors.New("gecko server initialized without database")
	}
	if server.jwtApp == nil {
		return nil, errors.New("gecko server initialized without JWT app")
	}
	if server.logger == nil {
		return nil, errors.New("gecko server initialized without logger")
	}
	server.logger.Info("DB: %#v, JWTApp: %#v, Logger: %#v", server.db, server.jwtApp, server.logger)
	return server, nil
}

func (server *Server) MakeRouter() *iris.Application {
	router := iris.New()
	if router == nil {
		server.logger.Error("Failed to initialize router")
	}
	router.Use(recoveryMiddleware)
	router.OnErrorCode(iris.StatusNotFound, handleNotFound)
	router.Get("/health", server.handleHealth)
	router.Get("/config/{configId}", server.handleConfigGET)
	router.Put("/config/{configId}", server.handleConfigPUT)
	router.Delete("/config/{configId}", server.handleConfigDELETE)

	// Optionally keep UseRouter if needed, with safety checks
	router.UseRouter(func(ctx iris.Context) {
		req := ctx.Request()
		if req == nil || req.URL == nil {
			server.logger.Warning("Request or URL is nil")
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString("Internal Server Error")
			return
		}
		req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
		ctx.Next()
	})

	// Build the router to ensure it's ready for net/http
	if err := router.Build(); err != nil {
		server.logger.Error("Failed to build Iris router: %v", err)
	}
	return router
}

func recoveryMiddleware(ctx iris.Context) {
	defer func() {
		if r := recover(); r != nil {
			ctx.Application().Logger().Errorf("panic recovered: %v", r)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("Internal Server Error")
		}
	}()
	ctx.Next()
}

func (server *Server) handleConfigGET(ctx iris.Context) {
	configId := ctx.Params().Get("configId")
	doc, err := configGET(server.db, configId)
	if doc == nil && err == nil {
		msg := fmt.Sprintf("no configId found with configId: %s", configId)
		errResponse := newErrorResponse(msg, 404, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}
	if err != nil {
		msg := fmt.Sprintf("config query failed: %s", err.Error())
		errResponse := newErrorResponse(msg, 500, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}
	server.logger.Info("%#v", doc)
	_ = jsonResponseFrom(doc, http.StatusOK).write(ctx)
}

func (server *Server) handleConfigDELETE(ctx iris.Context) {
	configId := ctx.Params().Get("configId")
	doc, err := configDELETE(server.db, configId)
	if doc == false && err == nil {
		msg := fmt.Sprintf("no configId found with configId: %s", configId)
		errResponse := newErrorResponse(msg, 404, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}
	if err != nil {
		msg := fmt.Sprintf("config query failed: %s", err.Error())
		errResponse := newErrorResponse(msg, 500, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}

	okmsg := map[string]any{"code": 200, "message": fmt.Sprintf("DELETED: %s", configId)}
	server.logger.Info("%#v", okmsg)
	_ = jsonResponseFrom(okmsg, http.StatusOK).write(ctx)
}

func (server *Server) handleConfigPUT(ctx iris.Context) {
	configId := ctx.Params().Get("configId")
	data := []config.ConfigItem{}
	body, err := ctx.GetBody()
	if err != nil {
		msg := fmt.Sprintf("GetBody() failed: %s", err.Error())
		errResponse := newErrorResponse(msg, 500, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}
	if !json.Valid(body) {
		msg := "Invalid JSON format"
		errResponse := newErrorResponse(msg, 400, nil) // 400 Bad Request
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}
	errResponse := unmarshal(body, &data)
	if errResponse != nil {
		msg := fmt.Sprintf("body data unmarshal failed: %s", err.Error())
		errResponse := newErrorResponse(msg, 500, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}
	err = configPUT(server.db, configId, data)
	if err != nil {
		msg := fmt.Sprintf("configPut failed: %s", err.Error())
		errResponse := newErrorResponse(msg, 500, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}

	okmsg := map[string]any{"code": 200, "message": fmt.Sprintf("ACCEPTED: %s", configId)}
	server.logger.Info("%#v", okmsg)
	_ = jsonResponseFrom(okmsg, http.StatusOK).write(ctx)
}

func (server *Server) handleHealth(ctx iris.Context) {
	server.logger.Info("Entering handleHealth")
	err := server.db.Ping()
	if err != nil {
		server.logger.Error("Database ping failed: %v", err)
		response := newErrorResponse("database unavailable", 500, nil)
		_ = response.write(ctx)
		return
	}
	server.logger.Info("Health check passed")
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

func unmarshal(body []byte, x any) *ErrorResponse {
	if len(body) == 0 {
		return newErrorResponse("empty request body", http.StatusBadRequest, nil)
	}

	err := json.Unmarshal(body, x)
	if err != nil {
		structType := reflect.TypeOf(x)
		if structType.Kind() == reflect.Ptr {
			structType = structType.Elem()
		}

		msg := fmt.Sprintf(
			"could not parse %s from JSON; make sure input has correct types",
			structType,
		)
		response := newErrorResponse(msg, http.StatusBadRequest, &err)
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
