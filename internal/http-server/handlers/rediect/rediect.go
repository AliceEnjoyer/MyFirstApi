package rediect

import (
	"net/http"

	resp "github.com/AliceEnjoyer/MyFirstApi/internal/lib/api/response"
	"github.com/AliceEnjoyer/MyFirstApi/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Request struct {
	Alias string `json:"alias" validate:"required,alias"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.rediect.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := r.URL.String()[1:]

		url, err := urlGetter.GetUrl(alias)
		if err != nil {
			if err == storage.ErrURLNotFound {
				log.Info("this alias does not exist", slog.String("alias", alias))

				render.JSON(w, r, resp.Error("this alias does not exist"))

				return
			}
			log.Info("cannot get url", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("cannot get url"))

			return
		}

		log.Info("got url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}
