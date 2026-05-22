package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func NewLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.String(a.Key, a.Value.Time().Format(time.DateTime))
			}
			return a
		},
	}))

	slog.SetDefault(log)
	return log
}


func MiddleWareLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		ctx.Next()

		latency := time.Since(start)

		logs := []slog.Attr{
			slog.Int("status", ctx.Writer.Status()),
			slog.String("method", ctx.Request.Method),
			slog.String("path", path),
			slog.String("ip", ctx.ClientIP()),
			slog.String("duration", latency.String()),
			slog.String("user_agent", ctx.Request.UserAgent()),
		}

		if query != "" {
			logs = append(logs, slog.String("query", query))
		}

		if len(ctx.Errors) > 0 {
			for _, err := range ctx.Errors {
				slog.Error(err.Error(), slog.String("path", path))
			}
			return
		}

		if ctx.Writer.Status() >= 500 {
			slog.Error("HTTP запрос неуспешный", ToAny(logs)...)
		} else {
			slog.Info("HTTP запрос", ToAny(logs)...)
		}
	}
}

func ToAny(attributes []slog.Attr) []any {
	result := make([]any, len(attributes))

	for i, v := range attributes {
		result[i] = v
	}
	return result
}
