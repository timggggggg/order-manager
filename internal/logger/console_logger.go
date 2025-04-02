// nolint
package logger

import (
	"context"
	"fmt"
	"io"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
)

func colorStr(color string, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, reset)
}

type ConsoleLogger struct {
	output io.Writer
}

func NewConsoleLogger(output io.Writer) *ConsoleLogger {
	return &ConsoleLogger{
		output: output,
	}
}

func (l *ConsoleLogger) Print(v ...any) {
	fmt.Fprint(l.output, v...)
}

func (l *ConsoleLogger) Printf(format string, v ...any) {
	fmt.Fprintf(l.output, format, v...)
}

func (l *ConsoleLogger) Println(v ...any) {
	fmt.Fprintln(l.output, v...)
}

func (l *ConsoleLogger) PrintError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("[ERROR]: %s\n", err)
}

func (l *ConsoleLogger) LogStatusChange(ctx context.Context, ts time.Time, id int64, statusFrom, statusTo models.OrderStatus) {
	fmt.Fprintf(
		l.output,
		"[INFO]: Order ID=%s status changed: %s -> %s. Time: %s\n",
		fmt.Sprintf("%s%d%s", yellow, id, reset),
		colorStr(red, string(statusFrom)),
		colorStr(green, string(statusTo)),
		models.FormatTime(ts),
	)
}

func (l *ConsoleLogger) LogRequest(ctx context.Context, ts time.Time, method, url, request_body string) {
	fmt.Fprintf(
		l.output,
		"[INFO]: http request: method=%s\n url=%s\n body=%s\n Time: %s\n\n",
		method,
		url,
		request_body,
		models.FormatTime(ts),
	)
}

func (l *ConsoleLogger) LogResponse(ctx context.Context, ts time.Time, code int64, body string) {
	fmt.Fprintf(
		l.output,
		"[INFO]: http response: code=%d\n body=%s\n Time: %s\n\n",
		code,
		body,
		models.FormatTime(ts),
	)
}
