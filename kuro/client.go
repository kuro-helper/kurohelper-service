package kuro

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var ErrUnavailable = errors.New("kuro AI runtime is unavailable")

type Config struct {
	URL            string
	Secret         string
	ReconnectDelay time.Duration
	RequestTimeout time.Duration
}

type envelope struct {
	Version   int             `json:"version"`
	Type      string          `json:"type"`
	RequestID string          `json:"requestId,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Error     *protocolError  `json:"error,omitempty"`
}

type protocolError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type requestResult struct {
	payload json.RawMessage
	err     error
}

type Client struct {
	config    Config
	connMu    sync.RWMutex
	conn      *websocket.Conn
	writeMu   sync.Mutex
	pendingMu sync.Mutex
	pending   map[string]chan requestResult
	connected atomic.Bool
}

func NewClient(config Config) (*Client, error) {
	parsed, err := url.Parse(config.URL)
	if err != nil || (parsed.Scheme != "ws" && parsed.Scheme != "wss") {
		return nil, fmt.Errorf("invalid Kuro runtime WebSocket URL: %q", config.URL)
	}
	if config.Secret == "" {
		return nil, errors.New("Kuro runtime bridge secret is required")
	}
	if config.ReconnectDelay <= 0 {
		config.ReconnectDelay = 5 * time.Second
	}
	if config.RequestTimeout <= 0 {
		config.RequestTimeout = 3 * time.Minute
	}
	return &Client{config: config, pending: make(map[string]chan requestResult)}, nil
}

func (client *Client) Start(ctx context.Context) {
	go client.connectLoop(ctx)
}

func (client *Client) Connected() bool {
	return client.connected.Load()
}

func (client *Client) Close() error {
	client.connMu.Lock()
	defer client.connMu.Unlock()
	client.connected.Store(false)
	if client.conn == nil {
		return nil
	}
	err := client.conn.Close()
	client.conn = nil
	return err
}

func (client *Client) connectLoop(ctx context.Context) {
	for ctx.Err() == nil {
		header := http.Header{}
		header.Set("Authorization", "Bearer "+client.config.Secret)
		conn, _, err := websocket.DefaultDialer.DialContext(ctx, client.config.URL, header)
		if err != nil {
			slog.Warn("Kuro AI runtime connection failed", "error", err)
			if !waitContext(ctx, client.config.ReconnectDelay) {
				return
			}
			continue
		}

		client.connMu.Lock()
		client.conn = conn
		client.connMu.Unlock()
		client.connected.Store(true)
		slog.Info("Kuro AI runtime connected", "url", client.config.URL)

		err = client.readLoop(ctx, conn)
		client.connected.Store(false)
		client.connMu.Lock()
		if client.conn == conn {
			client.conn = nil
		}
		client.connMu.Unlock()
		_ = conn.Close()
		client.failPending(ErrUnavailable)
		if ctx.Err() == nil {
			slog.Warn("Kuro AI runtime disconnected", "error", err)
			if !waitContext(ctx, client.config.ReconnectDelay) {
				return
			}
		}
	}
}

func (client *Client) readLoop(ctx context.Context, conn *websocket.Conn) error {
	for ctx.Err() == nil {
		var message envelope
		if err := conn.ReadJSON(&message); err != nil {
			return err
		}
		if message.Version != ProtocolVersion || message.RequestID == "" {
			continue
		}
		if message.Type == "typing" || message.Type == "images" {
			continue
		}
		client.pendingMu.Lock()
		resultChannel := client.pending[message.RequestID]
		if resultChannel != nil {
			delete(client.pending, message.RequestID)
		}
		client.pendingMu.Unlock()
		if resultChannel == nil {
			continue
		}
		if message.Error != nil {
			resultChannel <- requestResult{err: fmt.Errorf("%s: %s", message.Error.Code, message.Error.Message)}
		} else {
			resultChannel <- requestResult{payload: message.Payload}
		}
	}
	return ctx.Err()
}

func (client *Client) request(ctx context.Context, messageType string, requestID string, payload any, output any) error {
	if requestID == "" {
		requestID = randomID()
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resultChannel := make(chan requestResult, 1)
	client.pendingMu.Lock()
	client.pending[requestID] = resultChannel
	client.pendingMu.Unlock()

	client.connMu.RLock()
	conn := client.conn
	client.connMu.RUnlock()
	if conn == nil || !client.connected.Load() {
		client.removePending(requestID)
		return ErrUnavailable
	}

	client.writeMu.Lock()
	err = conn.WriteJSON(envelope{
		Version: ProtocolVersion, Type: messageType, RequestID: requestID, Payload: encoded,
	})
	client.writeMu.Unlock()
	if err != nil {
		client.removePending(requestID)
		return err
	}

	waitCtx := ctx
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, client.config.RequestTimeout)
		defer cancel()
	}
	select {
	case result := <-resultChannel:
		if result.err != nil {
			return result.err
		}
		if output == nil || len(result.payload) == 0 {
			return nil
		}
		return json.Unmarshal(result.payload, output)
	case <-waitCtx.Done():
		client.removePending(requestID)
		return waitCtx.Err()
	}
}

func (client *Client) Generate(ctx context.Context, request GenerateRequest) (GenerateResponse, error) {
	var response GenerateResponse
	err := client.request(ctx, "generate_request", request.RequestID, request, &response)
	return response, err
}

func (client *Client) ListMemories(ctx context.Context, status string, limit, offset int) (MemoryResponse, error) {
	return client.memory(ctx, MemoryRequest{Action: "list", Status: status, Limit: limit, Offset: offset})
}

func (client *Client) ForgetMemory(ctx context.Context, memoryID string) (MemoryResponse, error) {
	return client.memory(ctx, MemoryRequest{Action: "forget", MemoryID: memoryID})
}

func (client *Client) RestoreMemory(ctx context.Context, memoryID string) (MemoryResponse, error) {
	return client.memory(ctx, MemoryRequest{Action: "restore", MemoryID: memoryID})
}

func (client *Client) ClearMemories(ctx context.Context) (MemoryResponse, error) {
	return client.memory(ctx, MemoryRequest{Action: "clear"})
}

func (client *Client) ListMemoryBackups(ctx context.Context, limit, offset int) (MemoryResponse, error) {
	return client.memory(ctx, MemoryRequest{Action: "backup_list", Limit: limit, Offset: offset})
}

func (client *Client) CreateMemoryBackup(ctx context.Context) (MemoryResponse, error) {
	return client.memory(ctx, MemoryRequest{Action: "backup_create"})
}

func (client *Client) RestoreMemoryBackup(ctx context.Context, backupID string) (MemoryResponse, error) {
	return client.memory(ctx, MemoryRequest{Action: "backup_restore", BackupID: backupID})
}

func (client *Client) memory(ctx context.Context, request MemoryRequest) (MemoryResponse, error) {
	var response MemoryResponse
	err := client.request(ctx, "memory_request", "", request, &response)
	return response, err
}

func (client *Client) Health(ctx context.Context) (HealthResponse, error) {
	var response HealthResponse
	err := client.request(ctx, "health_request", "", struct{}{}, &response)
	return response, err
}

func (client *Client) removePending(requestID string) {
	client.pendingMu.Lock()
	delete(client.pending, requestID)
	client.pendingMu.Unlock()
}

func (client *Client) failPending(err error) {
	client.pendingMu.Lock()
	pending := client.pending
	client.pending = make(map[string]chan requestResult)
	client.pendingMu.Unlock()
	for _, channel := range pending {
		channel <- requestResult{err: err}
	}
}

func waitContext(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func randomID() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buffer)
}
