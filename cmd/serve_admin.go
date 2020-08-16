package cmd

import (
	"net/http"

	"github.com/kataras/muxie"
	"github.com/spf13/cobra"
)

func createServeAdminCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "admin",
		Short: "Backend server of the application.",
		Run: func(cmd *cobra.Command, args []string) {
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
					if r.Method == http.MethodOptions {
						w.WriteHeader(200)
					} else {
						next.ServeHTTP(w, r)
					}
				})
			})

			// route.SetupAdminServerV1(m.Of("/v1"))

			// if config.ProvideSingleton().Port == "" {
			// 	config.ProvideSingleton().Port = ":8080"
			// }

			// bog.Info("server started", zap.String("addr", config.ProvideSingleton().Port))
			// err := http.ListenAndServe(config.ProvideSingleton().Port, m)
			// if err != nil {
			// 	bog.Error("server error", zap.Error(err))
			// }
		},
	}
}
