package basic

import (
	"log/slog"
)

func main() {
	slog.Info("Starting server")  // want "log message should start with a lowercase letter"
	slog.Info("starting server")  // OK
	slog.Info("запуск сервера")   // want "log message must contain only English letters, digits, and spaces"
	slog.Info("server started!🚀") // want "log message contains special characters or emojis"
}
