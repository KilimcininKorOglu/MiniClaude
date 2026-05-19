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

func TestHubSendsToTargetSession(t *testing.T) {
	hub := NewHub()
	target, unsubscribeTarget := hub.SubscribeSession("workspace-a", "client-a", "session-a")
	defer unsubscribeTarget()
	other, unsubscribeOther := hub.SubscribeSession("workspace-a", "client-a", "session-b")
	defer unsubscribeOther()

	if !hub.SendSession("workspace-a", "session-a", Message{Type: "terminate_session"}) {
		t.Fatal("expected target session delivery")
	}

	select {
	case message := <-target:
		if message.Type != "terminate_session" {
			t.Fatalf("unexpected target message: %+v", message)
		}
	default:
		t.Fatal("expected target session message")
	}
	select {
	case message := <-other:
		t.Fatalf("unexpected other session message: %+v", message)
	default:
	}
}

func TestHubBroadcastExceptSkipsSourceSession(t *testing.T) {
	hub := NewHub()
	source, unsubscribeSource := hub.SubscribeSession("workspace-a", "client-a", "session-a")
	defer unsubscribeSource()
	other, unsubscribeOther := hub.SubscribeSession("workspace-a", "client-b", "session-b")
	defer unsubscribeOther()

	hub.BroadcastExcept("workspace-a", "session-a", Message{Type: "settings_updated"})

	select {
	case message := <-source:
		t.Fatalf("unexpected source session message: %+v", message)
	default:
	}
	select {
	case message := <-other:
		if message.Type != "settings_updated" {
			t.Fatalf("unexpected other message: %+v", message)
		}
	default:
		t.Fatal("expected other session message")
	}
}
