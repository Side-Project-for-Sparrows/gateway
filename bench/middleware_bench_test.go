package bench

import (
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

const (
	token     = "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIzIiwidHlwZSI6ImFjY2VzcyIsImlhdCI6MTc1MjMwODkxNiwiZXhwIjoxMDM5MjIyMjUxNn0.Qtdr8tx5_yEHRyllWh47-zYP9Ow2KUy0ANXCuz-bg5Jlfj_qjO2NMXQ4k_CAB0_TVCk5fPsLcq-EF7BtwcNYMqFJEJVbeeZgFdYdkOY2VIwYA-mC7tqtE1DnjAV8EGOJBsBiztgRyzpbNeG_oniZaxOTOmhhvjsLPownZ_FXfAFT6dcumcwWbUTZIAzaXbGdipXFamvLhlS2ugUHhU-I5-wy-U4IgZgrB-cm1IUM-A2kLJA8cmOhGVUOcEireozvpD7JR5KQbzV9SlExf60l5cia7JPnqiedlnrBjOflQsdSUxlD3dnJFdeit8QTHjdSTFJcWn46zuAT8xE2qi2_mA"
	queryURL  = `http://localhost:7080/board/search?query="지렁이"`
	testCount = 10000
)

func TestRealisticLatency(t *testing.T) {
	v := 2
	log.Print(v)

	client := &http.Client{}
	const iterations = testCount
	var total time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()

		req, _ := http.NewRequest("GET", queryURL, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		total += time.Since(start)
	}

	avg := total / iterations
	t.Logf("Average latency over %d requests: %s", iterations, avg)
}
