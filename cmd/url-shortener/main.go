package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/AliceEnjoyer/MyFirstApi/internal/config"
	"github.com/AliceEnjoyer/MyFirstApi/internal/lib/logger/sl"
	"github.com/AliceEnjoyer/MyFirstApi/internal/storage/sqlite"
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

	// бд (12 - )
	storage, err := sqlite.NewDatabase(cnfg.StoragePath)
	if err != nil {
		log.Error("can not init db", sl.Err(err))
		os.Exit(1)
	}
	log.Info("starting database", slog.String("env", cnfg.Env))
	_ = storage
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
