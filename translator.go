package webhooktomqtt

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type mqttClient interface {
	Publish(topic string, payload []byte)
}

type webhooktomqttHandler struct {
	client mqttClient
}

func NewHandler(client mqttClient) http.Handler {
	return &webhooktomqttHandler{client: client}
}

func (h *webhooktomqttHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request\tmethod=%s path=%s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, `{"status":"error"}`, http.StatusMethodNotAllowed)
		return
	}

	topic := strings.TrimPrefix(r.URL.Path, "/api/")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %s", err)
		http.Error(w, `{"status":"error"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("publishing payload\ttopic=%s bytes=%d", topic, len(payload))
	startTime := time.Now()

	h.client.Publish(topic, payload)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success"}`))

	log.Printf("published payload\ttopic=%s elapsed=%v", topic, time.Since(startTime))
}
