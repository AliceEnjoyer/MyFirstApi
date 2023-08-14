package delete

import (
	"errors"
	"io"
	"net/http"

	resp "github.com/AliceEnjoyer/MyFirstApi/internal/lib/api/response"
	"github.com/AliceEnjoyer/MyFirstApi/internal/lib/logger/sl"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Request struct {
	Alias string `json:"alias" validate:"required,alias"`
}

type URLDeleter interface {
	DeleteUrl(alias string) error
}

func New(log *slog.Logger, URLDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.rediect.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if len(req.Alias) == 0 {
			log.Error("request failed validation")

			render.JSON(w, r, resp.Error("request failed validation"))

			return
		}

		if err = URLDeleter.DeleteUrl(req.Alias); err != nil {
			log.Error("cannot delete alias", sl.Err(err))

			render.JSON(w, r, resp.Error("deleting alias failed"))

			return
		}

		log.Error("successfully deleted")

		render.JSON(w, r, resp.Ok())
	}
}
