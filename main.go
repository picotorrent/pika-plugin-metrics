package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type JsonRpcResponse struct {
	Version string           `json:"jsonrpc"`
	Result  map[string]int64 `json:"result"`
}

func handleMetrics(client *resty.Client) gin.HandlerFunc {
	return func(context *gin.Context) {
		result := JsonRpcResponse{}

		_, err := client.R().
			SetResult(&result).
			SetBody(map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1337,
				"method":  "session.stats",
			}).
			Post("/api/jsonrpc")

		if err != nil {
			context.Data(
				http.StatusInternalServerError,
				"text/plain",
				[]byte(err.Error()))
			return
		}

		context.Status(http.StatusOK)
		context.Writer.Header().Set("Content-Type", "text/plain; version=0.0.4")

		keys := make([]string, len(result.Result))
		i := 0

		for k := range result.Result {
			keys[i] = k
			i++
		}

		sort.Strings(keys)

		for _, k := range keys {
			v := result.Result[k]
			fmt.Fprintf(context.Writer, "pika_%s %d\n", strings.ReplaceAll(k, ".", "_"), v)
		}
	}
}

var (
	urlPtr *string = flag.String("url", "http://localhost:1337", "")
)

func main() {
	urlEnv, urlEnvFound := os.LookupEnv("PIKA_URL")

	if urlEnvFound {
		urlPtr = &urlEnv
	}

	c := resty.New().
		SetBaseURL(*urlPtr)

	r := gin.Default()
	r.SetTrustedProxies([]string{})
	r.GET("/metrics", handleMetrics(c))
	r.Run()
}
