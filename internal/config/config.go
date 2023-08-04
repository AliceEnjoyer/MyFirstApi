package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv" //  Minimalistic configuration reader
)

/*
6) создаем папку internal, здесь будет лежать всякая внутреняя машинерия,
например парсинг конфиа (в этом файле)
*/

/*
 7. создаем структуру конфига первым делом и он будет полностю

соответствовать yaml файлу, который сделаный пару шахов назад
*/
type Config struct {
	/* Тэг yaml определяет какое имя будет у соответствующего параметра в yaml файле.
	env-default - значение по умолчанию, если в конфиге значение этого параметра отсутствует,
	но это не безопасно, так как при какой-то причине конфиг файл может потеряться и установиться
	не верное для текущего случая (деплоя на рподакшн например) значение из env-default, по этому
	для конечного проекта можно указать env-required: "true", что бы приложение не запустилось
	при отсутствии параметра в конфиг файле.*/
	Env        string `yaml"env" env-required: "true"` // `yaml"env" env-default"local"` - было
	StorgePath string `yaml"storage_path" env-required: "true"`
	HTTPServer `yaml"http_server"`
}

/* 8) */
type HTTPServer struct {
	Address     string        `yaml"address" env-default"localhost:8080"`
	Timeout     time.Duration `yaml"timeout" env-required: "true"` // можно env-default"10s"
	IdleTimeout time.Duration `yaml"idle_timeout" env-required: "true"`
}

/*
 9. теперь нам нужно написать функцию, которая

прочитает файл с конфигом и создаст и дополнит объект конфиг,
который был только что написан.
*/
func MustLoad(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}
	var cnfg Config
	if err := cleanenv.ReadConfig(configPath, &cnfg); err != nil {
		log.Fatalf("cannot read config: %s", configPath)
	}

	return cnfg
}
