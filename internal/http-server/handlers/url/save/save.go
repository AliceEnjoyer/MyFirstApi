package save

import (
	"errors"
	"io"
	"net/http"

	resp "github.com/AliceEnjoyer/MyFirstApi/internal/lib/api/response"
	"github.com/AliceEnjoyer/MyFirstApi/internal/lib/logger/sl"
	"github.com/AliceEnjoyer/MyFirstApi/internal/lib/random"
	"github.com/AliceEnjoyer/MyFirstApi/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render" // The render package helps manage HTTP request / response payloads.
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

/*
15) в папке handlers будут те хэндлеры, которые обрабатывают поступившие запросы.
в папке  handlers/url будут находиться те хендлеры, которые работают с url.

Первый метод в этом апи, который реализируеться в видео - это save

Последующая конструкция являеться типичным конструктором для хендлера:
---
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
---

*/

/*
"Тут мы описываем запрос, которые будут поступать к нам
в виде json (будет приходить POST запрос, где будет находиться
json объект, который описывает url, который нужно сохранить)"
*/
type Request struct {
	/* validate будет давать инфу пакету github.com/go-playground/validator/v10,
	он говорит о том, что  это обязательное поле, то есть если этого поля нет,
	то получим при валидации ошибку.*/
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"` // omitempty - это просто опциональное поле в json файлике
}

// а это прост ответ на данный запрос
type Response struct {
	resp.Response        // прочитай в файле github.com/AliceEnjoyer/MyFirstApi/internal/lib/api/response
	Alias         string `json:"alias,omitempty"`
}

/*
Почему мы не сделали обобщающий интерфейс в 12, как в тг боте,
а сделали отдельный интерфейс для конкретного хендлера?

Мы определили интерфейс к месту использования (как сказал Николай Тузов)

А вообще, как я понял, это само по себе очень красиво и удобно:
нет лишних импортов - тупо интерфейс с онли интересующими нас методами,
возможно даже в этом есть некая оптимизация.
*/
type URLSaver interface {
	SaveUrl(urlToSave, alias string) (int64, error)
	IfExists(alias string) (bool, error)
}

const aliasLength = 6 // длинна псевдонимов для сокращений

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.save.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		/*
			r.Body - это сам json файл
		*/
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request")) // возвращает ответ серверу (клиенту)

			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		// делаем валидацию запроса
		if err := validator.New().Struct(req); err != nil {
			log.Error("request failed validation", sl.Err(err))

			render.JSON(w, r, resp.Error("request failed validation"))

			return
		}

		alias := req.Alias

		if alias == "" {
			for {
				alias = random.NewRandomAlias(aliasLength)
				check, err := urlSaver.IfExists(alias)
				if err != nil {
					log.Error("cannot check if alias exsts", slog.String("url", req.URL))

					render.JSON(w, r, resp.Error("cannot check if alias exsts"))

					return
				}
				if !check {
					break
				}
			}

		}

		id, err := urlSaver.SaveUrl(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Alias:    alias,
		})
	}
}
