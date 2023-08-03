package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/AliceEnjoyer/MyFirstApi/internal/config"
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

func main() {
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

	fmt.Println(cnfg.HTTPServer.Address)
}
