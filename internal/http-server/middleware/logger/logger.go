package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

// 14) создаем свой middleware для логгирования

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		) // создаем копию логгера, добавляя подсказу, что это компонент мидлвэйр логгер

		log.Info("logger middleware enabled") // выведиться только один раз при запуске

		// а от тут уже содержимое внутреней части хендлера (будет выполняться при каждом запросе)
		// r - запрос, а w -
		fn := func(w http.ResponseWriter, r *http.Request) {
			// 11 - эта та часть, которая работает до обработки запроса
			// 11
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			/* получает сведения об ответе (из него мы получаем в defer статус, bytes written)*/
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			// 11

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r) // тут мы передаем управление след. хендлеру по цепочке,
			// а после вызоветься defer, уже после всей обработки запроса.
		}

		return http.HandlerFunc(fn)
	}
}
