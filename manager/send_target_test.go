package manager

import (
	"encoding/json"
	"testing"

	"github.com/wapikit/wapi.go/pkg/components"
)

func newTextMsg(t *testing.T) components.BaseMessage {
	t.Helper()
	msg, err := components.NewTextMessage(components.TextMessageConfigs{Text: "hi"})
	if err != nil {
		t.Fatalf("NewTextMessage: %v", err)
	}
	return msg
}

// Phone targets must serialize `to` and never `recipient` (backward compatible).
func TestPhoneTargetSerializesTo(t *testing.T) {
	body, err := newTextMsg(t).ToJson(NewPhoneTarget("919999999999").configs(""))
	if err != nil {
		t.Fatalf("ToJson: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload["to"] != "919999999999" {
		t.Fatalf("expected to=919999999999, got %v", payload["to"])
	}
	if _, ok := payload["recipient"]; ok {
		t.Fatalf("phone target must not include recipient: %v", payload["recipient"])
	}
}

// BSUID targets must serialize `recipient` and never `to`.
func TestBSUIDTargetSerializesRecipient(t *testing.T) {
	body, err := newTextMsg(t).ToJson(NewBSUIDTarget("IN.1013756321578695").configs(""))
	if err != nil {
		t.Fatalf("ToJson: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload["recipient"] != "IN.1013756321578695" {
		t.Fatalf("expected recipient=IN.1013756321578695, got %v", payload["recipient"])
	}
	if _, ok := payload["to"]; ok {
		t.Fatalf("bsuid target must not include to: %v", payload["to"])
	}
}

// Parent-BSUID targets also serialize `recipient`.
func TestParentBSUIDTargetSerializesRecipient(t *testing.T) {
	body, _ := newTextMsg(t).ToJson(NewParentBSUIDTarget("IN.parent.42").configs(""))
	var payload map[string]any
	_ = json.Unmarshal(body, &payload)
	if payload["recipient"] != "IN.parent.42" {
		t.Fatalf("expected recipient=IN.parent.42, got %v", payload["recipient"])
	}
}
