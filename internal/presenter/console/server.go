package console

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rumi-go/internal/config"
	"rumi-go/internal/infrastructure"
	handlers "rumi-go/internal/presenter/rest"
	eventHandlers "rumi-go/internal/presenter/rest/handlers/event"
	eventReportHandlers "rumi-go/internal/presenter/rest/handlers/event_report"
	"rumi-go/internal/usecase/event"
	"rumi-go/internal/usecase/event_report"
	"rumi-go/internal/usecase/user"

	"github.com/gorilla/mux"

	"gopkg.in/ukautz/clif.v1"
)

func (c *Console) StartServer() *clif.Command {
	return clif.NewCommand("start", "starting http server.", func(o *clif.Command, in clif.Input, out clif.Output) error {
		ctx := context.Background()
		err := config.Load(".")
		if err != nil {
			return err
		}
		conf := config.Get()

		infra, err := infrastructure.NewInfra(ctx, conf)
		if err != nil {
			return err
		}

		httpHandler := SetupHandler(conf, infra)
		server := http.Server{
			Handler: httpHandler,
			Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		}
		gracefulShutdown := make(chan os.Signal, 1)
		signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			slog.InfoContext(ctx, fmt.Sprintf("server starting on port %d", conf.Server.Port))
			if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.ErrorContext(ctx, "server error", slog.String("error", err.Error()))
			}
		}()

		<-gracefulShutdown
		slog.WarnContext(ctx, "shutting down")
		time.Sleep(5 * time.Second)
		err = server.Shutdown(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "shutdown error", slog.String("error", err.Error()))
		}

		return nil
	})
}

func SetupHandler(conf *config.Config, infra *infrastructure.Infra) *mux.Router {
	// Initialize Modules
	userUsecase := user.NewUserUsecase(infra.Database)
	eventUsecase := event.NewEventUsecase(infra.Database)
	eventReportUsecase := event_report.NewEventReportUsecase(infra.Database)

	// Create user handler
	userHandler := handlers.NewUserHandler(userUsecase)
	eventHandler := eventHandlers.NewEventHandler(eventUsecase)
	eventReportHandler := eventReportHandlers.NewEventReportHandler(eventReportUsecase)

	// Initialize HTTP Handler
	baseHandler := mux.NewRouter()
	baseHandler.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	// Public routes (no authentication required)
	baseHandler.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)
	baseHandler.HandleFunc("/v1/user/login", userHandler.LoginUser).Methods(http.MethodPost)

	// adminHandler := baseHandler.PathPrefix("").Subrouter()
	// adminHandler.Use(middleware.AuthMiddleware, middleware.RBACMiddleware(middleware.AdminOnly))
	baseHandler.HandleFunc("/v1/user", userHandler.CreateUser).Methods(http.MethodPost)
	baseHandler.HandleFunc("/v1/users", userHandler.ListUsers).Methods(http.MethodGet)
	baseHandler.HandleFunc("/v1/user", userHandler.UpdateUser).Methods(http.MethodPut)
	baseHandler.HandleFunc("/v1/user", userHandler.DeleteUser).Methods(http.MethodDelete)
	baseHandler.HandleFunc("/v1/user", userHandler.GetUserInfo).Methods(http.MethodGet)
	baseHandler.HandleFunc("/v1/user/email", userHandler.GetUserByEmail).Methods(http.MethodGet)
	baseHandler.HandleFunc("/v1/event", eventHandler.CreateEvent).Methods(http.MethodPost)
	baseHandler.HandleFunc("/v1/events", eventHandler.ListEvents).Methods(http.MethodGet)
	baseHandler.HandleFunc("/v1/event/info", eventHandler.GetEvent).Methods(http.MethodGet)
	baseHandler.HandleFunc("/v1/event", eventHandler.UpdateEvent).Methods(http.MethodPut)
	baseHandler.HandleFunc("/v1/event", eventHandler.DeleteEvent).Methods(http.MethodDelete)
	baseHandler.HandleFunc("/v1/event/report", eventReportHandler.CreateEventReport).Methods(http.MethodPost)
	baseHandler.HandleFunc("/v1/event/reports", eventReportHandler.ListEventReports).Methods(http.MethodGet)
	baseHandler.HandleFunc("/v1/event/report/info", eventReportHandler.GetEventReport).Methods(http.MethodGet)
	baseHandler.HandleFunc("/v1/event/report", eventReportHandler.UpdateEventReport).Methods(http.MethodPut)
	baseHandler.HandleFunc("/v1/event/report", eventReportHandler.DeleteEventReport).Methods(http.MethodDelete)
	baseHandler.HandleFunc("/v1/event/reports/by-event", eventReportHandler.GetEventReportsByEventID).Methods(http.MethodGet)
	baseHandler.HandleFunc("/v1/event/reports/by-user", eventReportHandler.GetEventReportsByUserID).Methods(http.MethodGet)

	return baseHandler
}

func NotFoundHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	_, _ = writer.Write([]byte("not found"))
}

func pingHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("pong"))
}
