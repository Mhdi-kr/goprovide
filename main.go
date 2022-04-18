package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type context interface {
	*resty.Client
}

type authenticate func(c *resty.Client) bool
type available func() bool
type prebook func() bool
type book func() bool
type confirm func() bool
type cancel func() bool
type refund func() bool

type Provider struct {
	name    string
	actions ProviderActions
}

func newProvider(name string, actions ProviderActions) *Provider {
	client := resty.New()
	checkActionsInstance(actions)
	provider := &Provider{
		name:    name,
		actions: actions,
	}
	provider.actions.authenticate(client)
	return provider
}

type ProviderActions struct {
	authenticate authenticate
	available    available
	prebook      prebook
	book         book
	confirm      confirm
	cancel       cancel
	refund       refund
}

func checkActionsInstance(a ProviderActions) {
	if a.authenticate == nil {
		panic("authenticate method in passed actions is not defined")
	}
	if a.available == nil {
		panic("available method in passed actions is not defined")
	}
	if a.prebook == nil {
		panic("prebook method in passed actions is not defined")
	}
	if a.book == nil {
		panic("book method in passed actions is not defined")
	}
	if a.confirm == nil {
		panic("confirm method in passed actions is not defined")
	}
	if a.cancel == nil {
		panic("cancel method in passed actions is not defined")
	}
	if a.refund == nil {
		panic("refund method in passed actions is not defined")
	}
}

func newContainer() {

}

func NewZapLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	defer logger.Sync()
	return logger
}

func compactJson(jsonBytes []byte) string {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, jsonBytes); err != nil {
		fmt.Println(err)
	}
	return buffer.String()
}

func NewRestyClient() *resty.Client {
	logger := NewZapLogger()
	client := resty.New()
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		logger.Info("HTTP request",
			zap.String("endpoint", req.URL),
			zap.Int("attempt", 3),
			zap.Duration("backoff", time.Second),
		)
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		switch resp.StatusCode() {
		case http.StatusOK:
			logger.Info("HTTP response",
				zap.String("status", resp.Status()),
				zap.String("body", compactJson(resp.Body())),
				zap.Int("attempt", 3),
				zap.Duration("backoff", time.Second),
			)
		case http.StatusNotFound:
			logger.Info("HTTP response",
				zap.String("status", resp.Status()),
				zap.String("body", compactJson(resp.Body())),
				zap.Int("attempt", 3),
				zap.Duration("backoff", time.Second),
			)
		}

		return nil
	})

	client.OnError(func(req *resty.Request, err error) {
		fmt.Println("asd")
		if _, ok := err.(*resty.ResponseError); ok {
			// fmt.Print(v.Error())
			// v.Response contains the last response from the server
			// v.Err contains the original error
		}
		// Log the error, increment a metric, etc...
	})
	return client
}

func main() {
	// url := "google.com"
	httpClient := NewRestyClient()
	httpClient.R().Get("https://jsonplaceholder.typicode.com/todoss/1")
	// logger.Error("failed to fetch URL",
	// 	// Structured context as strongly typed Field values.
	// 	zap.String("url", url),
	// 	zap.Int("attempt", 3),
	// 	zap.Duration("backoff", time.Second),
	// )
	// logger.Fatal("failed to fetch URL",
	// 	// Structured context as strongly typed Field values.
	// 	zap.String("url", url),
	// 	zap.Int("attempt", 3),
	// 	zap.Duration("backoff", time.Second),
	// )
	// actions := ProviderActions{
	// 	authenticate: func(client *resty.Client) bool {
	// 		resp, err := client.R().Get("https://www.example.com")
	// 		fmt.Println(resp, err)
	// 		return true
	// 	},
	// 	available: func() bool {
	// 		return true
	// 	},
	// 	prebook: func() bool {
	// 		return true
	// 	},
	// 	book: func() bool {
	// 		return true
	// 	},
	// 	confirm: func() bool {
	// 		return true
	// 	},
	// 	cancel: func() bool {
	// 		return true
	// 	},
	// 	refund: func() bool {
	// 		return true
	// 	},
	// }
	// newProvider("ADOTEL", actions)
}
