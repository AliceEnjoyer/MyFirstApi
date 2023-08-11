package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AliceEnjoyer/MyFirstApi/internal/config"
	"github.com/AliceEnjoyer/MyFirstApi/internal/http-server/handlers/rediect"
	"github.com/AliceEnjoyer/MyFirstApi/internal/http-server/handlers/url/save"
	"github.com/AliceEnjoyer/MyFirstApi/internal/http-server/middleware/logger"
	"github.com/AliceEnjoyer/MyFirstApi/internal/lib/logger/sl"
	"github.com/AliceEnjoyer/MyFirstApi/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog" // это просто обертка текстовых и json логгеров
)

/*

rest api - это не более чем архитектурный подход
к построению веб приложений, которые общаються по протоколу http.

В данном дизайне системы все заимодействие из сервером сводиться к 4ем операциям:
- получение данных (HTTP запрос GET)
- добавление данных (POST)
- модификация сущесвтующих данных (PUT/PATCH)
- удаление данных (DELETE)

Эти операции еще называют CRUD - create, read update and delete.

1) Первым делом создаем файл, который будет запускать приложение (апи),
ми его разместили в папке cmd - это довольно распространеный подход,
здесь будут лежать различные команды, например одна из комманд - это запуск
этого приложения, но в будущем могут появиться какие-то другие комманды, типа
очистки кеша, выполнения каких-то операций и тд. Пока достаточно запуска приложения.

*/

// 11) сделаем текстовые константы что бы в дальнейшем можно было удобно
// менять значение env
const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// читаем конфиг (2 - 9)
	ConfigPath := flag.String(
		"ConfigPath",
		"",
		"",
	)

	flag.Parse()

	if *ConfigPath == "" {
		log.Fatal("ConfigPath is not specified")
	}

	cnfg := config.MustLoad(*ConfigPath)
	fmt.Println(cnfg)

	// инициализируем логгер (10 - 11)
	log := setupLogger(cnfg.Env)
	log.Info("starting url-shortner", slog.String("env", cnfg.Env))
	log.Debug("debug messages are enabled") // есди в &slog.HandlerOptions{ будет slog.LevelInfo, то log.Debug не будет работать

	// бд (12 - 13)
	storage, err := sqlite.NewDatabase(cnfg.StoragePath)
	if err != nil {
		log.Error("can not init db", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	/* инициализируем роутер */
	router := chi.NewRouter()

	/* подключааем к роутеру middleware. Что это такоей?
	Есть основные хендлеры запросов, типа  добавления или удаления url,
	а есть серединный хендлер, без которого не будут выполняться основные.
	К примеру мы хотим добавить ссылку, но middleware нам говорит авторизироваться,
	иначе мы НЕ ПРОЙДЕМ!!!*/

	router.Use(middleware.RequestID) // этот middleware добавляет к каждому запросу request id,
	// что бы при ошибке знать на каком request id произошла ошибка.

	router.Use(middleware.RealIP) // этот middleware добавляет к каждому запросу ip пользователя этого сервера

	// подключаем созданый мидлвэйр логгер, который логирует запросы (14)
	router.Use(logger.New(log))

	// если случаеться какая-то паникавнутри хендлера, то сервер востанавливаеться (в исзходном коде все понятно)
	router.Use(middleware.Recoverer)

	/* как я понял, оно убирает продолжения (типа json) из
	пути, созраняя его в контексте (дальше будеи понятно). Пишет
	красивые url при подкоючении к роутеру*/
	router.Use(middleware.URLFormat)

	// подключаем хендлер, который обрабатывает ссылки для сокращения
	// 15: http-server.handlers.url.save
	router.Post("/", save.New(log, storage))

	// подключаем наш хендлер редиект
	router.Get("/{alias}", rediect.New(log, storage))

	log.Info("starting server", slog.String("address", cnfg.Address))

	// создаем сам сервер
	srv := &http.Server{
		Addr:         cnfg.Address,
		Handler:      router,
		ReadTimeout:  cnfg.HTTPServer.Timeout,
		WriteTimeout: cnfg.HTTPServer.Timeout,
		IdleTimeout:  cnfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("server stoped")
}

/*
 10. выносим инициализацию логгера из-за того, что

его установка будет зависить от параметра env. Он от него
зависит потому что на local мы бы хотели видеть логи в виде
текста в консольке, на dev - json файлики из удаленного сервака,
которые являт собой характер debug, а на prod - json
файлики, которые просто имеют функцию инфрмации.
*/
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		// тут делаем простой текстовый логгер (с параметрами и сам разберешся)
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo, // тут уже просто инфо, так как это прод
				}),
		)
	}

	return log
}
