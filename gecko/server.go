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
	server.logger.logger.Printf("DB: %#v, JWTApp: %#v, Logger: %#v", server.db, server.jwtApp, server.logger)
	return server, nil
}

func (server *Server) MakeRouter() *iris.Application {
	router := iris.New()
	if router == nil {
		log.Fatal("Failed to initialize router")
	}
	router.Use(recoveryMiddleware)
	router.Get("/", func(ctx iris.Context) {
		server.logger.logger.Println("Root handler called")
		ctx.JSON(iris.Map{"message": "Hello, World!"})
	})
	router.OnErrorCode(iris.StatusNotFound, handleNotFound)
	router.Get("/health", server.handleHealth)
	router.Get("/config/{configId}", server.handleConfigGET)
	router.Put("/config/{configId}", server.handleConfigPUT)
	// Optionally keep UseRouter if needed, with safety checks
	router.UseRouter(func(ctx iris.Context) {
		req := ctx.Request()
		if req == nil || req.URL == nil {
			log.Println("WARNING: Request or URL is nil")
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString("Internal Server Error")
			return
		}
		log.Println("REQUEST:", req)
		req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
		ctx.Next()
	})
	// Build the router to ensure it's ready for net/http
	if err := router.Build(); err != nil {
		log.Fatalf("Failed to build Iris router: %v", err)
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
	body, err := ctx.GetBody()
	if err != nil {
		msg := fmt.Sprintf("client query failed: %s", err.Error())
		errResponse := newErrorResponse(msg, 500, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return

	}
	errResponse := unmarshal(body, &data)
	if errResponse != nil {
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}
	err = configPUT(server.db, configId, data)
	if err != nil {
		msg := fmt.Sprintf("client query failed: %s", err.Error())
		errResponse := newErrorResponse(msg, 500, nil)
		errResponse.log.write(server.logger)
		_ = errResponse.write(ctx)
		return
	}

	_ = jsonResponseFrom(fmt.Sprintf("OK: %s", configId), http.StatusOK).write(ctx)
}

func (server *Server) handleHealth(ctx iris.Context) {
	server.logger.logger.Println("Entering handleHealth")
	err := server.db.Ping()
	if err != nil {
		server.logger.logger.Printf("Database ping failed: %v", err)
		response := newErrorResponse("database unavailable", 500, nil)
		_ = response.write(ctx)
		return
	}
	server.logger.logger.Println("Health check passed")
	_ = jsonResponseFrom("Healthy", http.StatusOK).write(ctx)
}

func handleNotFound(ctx iris.Context) {
	log.Println("HELLO ? ")
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
