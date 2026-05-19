package ws

import "sync"

type Message struct {
	Type        string `json:"type"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Version     int64  `json:"version,omitempty"`
	Payload     any    `json:"payload,omitempty"`
}

type subscriber struct {
	workspaceID string
	clientID    string
	sessionID   string
	ch          chan Message
}

type Hub struct {
	mu         sync.RWMutex
	workspaces map[string]map[chan Message]subscriber
	clients    map[string]map[chan Message]struct{}
	sessions   map[string]chan Message
	bufferSize int
}

func NewHub() *Hub {
	return &Hub{
		workspaces: map[string]map[chan Message]subscriber{},
		clients:    map[string]map[chan Message]struct{}{},
		sessions:   map[string]chan Message{},
		bufferSize: 16,
	}
}

func (h *Hub) Subscribe(workspaceID string) (<-chan Message, func()) {
	return h.subscribe(workspaceID, "", "")
}

func (h *Hub) SubscribeSession(workspaceID, clientID, sessionID string) (<-chan Message, func()) {
	return h.subscribe(workspaceID, clientID, sessionID)
}

func (h *Hub) subscribe(workspaceID, clientID, sessionID string) (<-chan Message, func()) {
	ch := make(chan Message, h.bufferSize)
	sub := subscriber{workspaceID: workspaceID, clientID: clientID, sessionID: sessionID, ch: ch}

	h.mu.Lock()
	if _, ok := h.workspaces[workspaceID]; !ok {
		h.workspaces[workspaceID] = map[chan Message]subscriber{}
	}
	h.workspaces[workspaceID][ch] = sub
	if clientID != "" {
		key := clientKey(workspaceID, clientID)
		if _, ok := h.clients[key]; !ok {
			h.clients[key] = map[chan Message]struct{}{}
		}
		h.clients[key][ch] = struct{}{}
	}
	if sessionID != "" {
		h.sessions[sessionKey(workspaceID, sessionID)] = ch
	}
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		if subscribers, ok := h.workspaces[workspaceID]; ok {
			delete(subscribers, ch)
			if len(subscribers) == 0 {
				delete(h.workspaces, workspaceID)
			}
		}
		if clientID != "" {
			key := clientKey(workspaceID, clientID)
			if subscribers, ok := h.clients[key]; ok {
				delete(subscribers, ch)
				if len(subscribers) == 0 {
					delete(h.clients, key)
				}
			}
		}
		if sessionID != "" {
			delete(h.sessions, sessionKey(workspaceID, sessionID))
		}
		h.mu.Unlock()
		close(ch)
	}
	return ch, unsubscribe
}

func (h *Hub) Broadcast(workspaceID string, message Message) {
	h.broadcast(workspaceID, "", message)
}

func (h *Hub) BroadcastExcept(workspaceID, sessionID string, message Message) {
	h.broadcast(workspaceID, sessionID, message)
}

func (h *Hub) SendClient(workspaceID, clientID string, message Message) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	delivered := false
	for ch := range h.clients[clientKey(workspaceID, clientID)] {
		if send(ch, message) {
			delivered = true
		}
	}
	return delivered
}

func (h *Hub) SendSession(workspaceID, sessionID string, message Message) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ch, ok := h.sessions[sessionKey(workspaceID, sessionID)]
	if !ok {
		return false
	}
	return send(ch, message)
}

func (h *Hub) broadcast(workspaceID, exceptSessionID string, message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch, sub := range h.workspaces[workspaceID] {
		if exceptSessionID != "" && sub.sessionID == exceptSessionID {
			continue
		}
		send(ch, message)
	}
}

func (h *Hub) SubscriberCount(workspaceID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.workspaces[workspaceID])
}

func send(ch chan Message, message Message) bool {
	select {
	case ch <- message:
		return true
	default:
		return false
	}
}

func clientKey(workspaceID, clientID string) string {
	return workspaceID + ":" + clientID
}

func sessionKey(workspaceID, sessionID string) string {
	return workspaceID + ":" + sessionID
}
