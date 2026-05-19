package ws

import "sync"

type Message struct {
	Type        string `json:"type"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Version     int64  `json:"version,omitempty"`
	Payload     any    `json:"payload,omitempty"`
}

type Hub struct {
	mu         sync.RWMutex
	workspaces map[string]map[chan Message]struct{}
	bufferSize int
}

func NewHub() *Hub {
	return &Hub{workspaces: map[string]map[chan Message]struct{}{}, bufferSize: 16}
}

func (h *Hub) Subscribe(workspaceID string) (<-chan Message, func()) {
	ch := make(chan Message, h.bufferSize)

	h.mu.Lock()
	if _, ok := h.workspaces[workspaceID]; !ok {
		h.workspaces[workspaceID] = map[chan Message]struct{}{}
	}
	h.workspaces[workspaceID][ch] = struct{}{}
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		if subscribers, ok := h.workspaces[workspaceID]; ok {
			delete(subscribers, ch)
			if len(subscribers) == 0 {
				delete(h.workspaces, workspaceID)
			}
		}
		h.mu.Unlock()
		close(ch)
	}
	return ch, unsubscribe
}

func (h *Hub) Broadcast(workspaceID string, message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.workspaces[workspaceID] {
		select {
		case ch <- message:
		default:
		}
	}
}

func (h *Hub) SubscriberCount(workspaceID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.workspaces[workspaceID])
}
