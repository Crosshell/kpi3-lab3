package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	for i := 1; i <= 10; i++ {
		x := float64(i) * 0.1
		y := float64(i) * 0.1

		script := fmt.Sprintf(`reset
figure %.2f %.2f
update
`, x, y)

		resp, err := http.Post("http://localhost:17000/", "text/plain", bytes.NewBufferString(script))
		if err != nil {
			log.Println("❌ Failed to send request:", err)
			continue
		}
		resp.Body.Close()
		log.Printf("✅ Sent move to (%.2f, %.2f)\n", x, y)
		time.Sleep(1 * time.Second)
	}
}
