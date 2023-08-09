package save

import (
	"errors"
	"io"
	"net/http"

	resp "github.com/AliceEnjoyer/MyFirstApi/internal/lib/api/response"
	"github.com/AliceEnjoyer/MyFirstApi/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render" // The render package helps manage HTTP request / response payloads.
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
}

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

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}
		/*

			я закончил тута(скопировал хуйню с 77ой строки потому что и так понятно)!!!!
			1:16:35!!!

		*/
	}
}
