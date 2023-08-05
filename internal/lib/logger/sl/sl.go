package sl

import "golang.org/x/exp/slog"

// это вспомагательные поюшки для slog

// эта функция нужна просто из-за того, что в slog нет
// удобного огирования ошибок
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
