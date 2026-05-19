package ws

import "testing"

func TestHubBroadcastsWithinWorkspace(t *testing.T) {
	hub := NewHub()
	updates, unsubscribe := hub.Subscribe("workspace-a")
	defer unsubscribe()
	hub.Subscribe("workspace-b")

	hub.Broadcast("workspace-a", Message{Type: "settings_updated", WorkspaceID: "workspace-a", Version: 2})

	select {
	case message := <-updates:
		if message.Type != "settings_updated" || message.Version != 2 {
			t.Fatalf("unexpected message: %+v", message)
		}
	default:
		t.Fatal("expected workspace update")
	}
	if count := hub.SubscriberCount("workspace-a"); count != 1 {
		t.Fatalf("expected one subscriber, got %d", count)
	}
}
