package cmd

import (
	"appsrv/app"
	"appsrv/pkg/config"
	"appsrv/pkg/runtime"
	"appsrv/sql"
	"context"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/kataras/muxie"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Runner func(context.Context)

func createApplicationRunner(run Runner) *cobra.Command {
	return &cobra.Command{
		Use:   "app",
		Short: "Frontend server of the application.",
		Run: func(cmd *cobra.Command, args []string) {
			run(context.Background())
		},
	}
}

func createServerMux(wrapper muxie.Wrapper, cs []runtime.Controller) *muxie.Mux {
	m := muxie.NewMux()
	m.PathCorrection = true
	m.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Add("Access-Control-Allow-Methods", "*")
			w.Header().Add("Access-Control-Allow-Headers", "Authorization,Content-Type")
			w.Header().Add("Access-Control-Max-Age", "600")
			w.Header().Add("Access-Control-Expose-Headers", "X-Refresh-Token")
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			if r.Method == http.MethodOptions {
				w.WriteHeader(200)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	})

	m.Use(wrapper)

	for _, ctrl := range cs {
		ctrl.RegisterRoute(m.Of("v1"))
	}

	return m
}

func createServerRunner(m *muxie.Mux, db *pg.DB, log *zap.Logger, conf *config.ServerConfig) Runner {
	sql.Run(db, log)

	if conf.Port == "" {
		conf.Port = ":8080"
	}

	return func(ctx context.Context) {
		log.Info("server started", zap.String("addr", conf.Port))
		err := http.ListenAndServe(conf.Port, m)
		if err != nil {
			log.Error("server error", zap.Error(err))
		}
	}
}

func resolveControllerSet() []runtime.Controller {
	return []runtime.Controller{
		app.NewUserController(),
		app.NewTopicController(),
		app.NewNotificationController(),
		app.NewCommentController(),
		app.NewAnnounceController(),
	}
}
