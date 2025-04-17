package main

import (
	"bytes"
	"log"
	"net/http"
)

func main() {
	script := `
reset
white
bgrect 0.25 0.25 0.75 0.75
figure 0.5 0.5
green
figure 0.6 0.6
update 

`

	resp, err := http.Post("http://localhost:17000/", "text/plain", bytes.NewBufferString(script))
	if err != nil {
		log.Fatal("❌ Failed to send request:", err)
	}
	defer resp.Body.Close()
	log.Println("✅ Green frame script sent!")
}
