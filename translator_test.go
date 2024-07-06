package webhooktomqtt

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	messages []struct {
		topic   string
		payload []byte
	}
}

func (m *mockClient) Publish(topic string, payload []byte) {
	m.messages = append(m.messages, struct {
		topic   string
		payload []byte
	}{topic: topic, payload: payload})
}

func TestServeHTTP(t *testing.T) {
	client := &mockClient{}
	h := NewHandler(client)
	body := strings.NewReader(`{"hello": "world"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/mypod", body)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()

	h.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Code, http.StatusCreated)
	assert.Equal(t, recorder.Body.String(), `{"status":"success"}`)
	assert.Len(t, client.messages, 1)
	assert.Equal(t, client.messages[0].topic, "mypod")
	assert.Equal(t, client.messages[0].payload, []byte(`{"hello": "world"}`))
}
