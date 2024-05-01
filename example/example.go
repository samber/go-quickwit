package main

import (
	"fmt"
	"time"

	"github.com/samber/go-quickwit"
)

func main() {
	// docker-compose up -d
	// curl -X POST \
	//     'http://localhost:7280/api/v1/indexes' \
	//     -H 'Content-Type: application/yaml' \
	//     --data-binary @test-config.yaml

	client := quickwit.NewWithDefault("http://localhost:7280")
	defer client.Stop() // flush and stop

	for i := 0; i < 10; i++ {
		msg := map[string]any{
			"timestamp": time.Now().Unix(),
			"message":   fmt.Sprintf("hello %d", i),
		}
		fmt.Println(msg)

		client.Push(msg)
		time.Sleep(1 * time.Second)
	}
}
