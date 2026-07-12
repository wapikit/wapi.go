package manager

import (
	"encoding/json"
	"testing"
)

// The supplied dual-identifier status fixture must parse recipient_id AND
// recipient_user_id together, and preserve the raw pricing fields.
func TestStatusDualIdentifierAndPricingParse(t *testing.T) {
	raw := `{
		"id":"wamid.HBgMOTE5ODMxODA3NDU1FQ==",
		"status":"sent",
		"timestamp":"1783743596",
		"recipient_id":"919831807455",
		"recipient_user_id":"IN.1013756321578695",
		"pricing":{"billable":true,"pricing_model":"PMP","category":"utility","type":"regular"}
	}`
	var s Status
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("unmarshal status: %v", err)
	}
	if s.RecipientId != "919831807455" {
		t.Fatalf("recipient_id=%q", s.RecipientId)
	}
	if s.RecipientUserId != "IN.1013756321578695" {
		t.Fatalf("recipient_user_id=%q", s.RecipientUserId)
	}
	if !s.Pricing.Billable || s.Pricing.PricingModel != "PMP" || string(s.Pricing.Category) != "utility" {
		t.Fatalf("pricing not preserved: %+v", s.Pricing)
	}
}

// The phone-omitted (BSUID-only) status variant must still parse and remain
// processable — recipient_id may be empty.
func TestStatusPhoneOmittedVariantParses(t *testing.T) {
	raw := `{"id":"wamid.X","status":"delivered","timestamp":"1","recipient_user_id":"IN.999"}`
	var s Status
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.RecipientUserId != "IN.999" {
		t.Fatalf("recipient_user_id=%q", s.RecipientUserId)
	}
	if s.RecipientId != "" {
		t.Fatalf("expected empty recipient_id, got %q", s.RecipientId)
	}
}

// Inbound contact identity: user_id, parent_user_id, and profile.username parse.
func TestSenderContactIdentityParse(t *testing.T) {
	raw := `{"wa_id":"919831807455","user_id":"IN.1013","parent_user_id":"IN.parent","profile":{"name":"Ravi","username":"ravi.k"}}`
	var c SenderContact
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		t.Fatalf("unmarshal contact: %v", err)
	}
	if c.UserId != "IN.1013" || c.ParentUserId != "IN.parent" {
		t.Fatalf("identity not parsed: %+v", c)
	}
	if c.Profile.Username != "ravi.k" || c.Profile.Name != "Ravi" {
		t.Fatalf("profile not parsed: %+v", c.Profile)
	}
}

// Message-level from_user_id / from_parent_user_id parse.
func TestMessageFromIdentityParse(t *testing.T) {
	raw := `{"id":"wamid.Y","from":"919831807455","from_user_id":"IN.1013","from_parent_user_id":"IN.parent","timestamp":"1","type":"text"}`
	var m Message
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("unmarshal message: %v", err)
	}
	if m.FromUserId != "IN.1013" || m.FromParentUserId != "IN.parent" {
		t.Fatalf("from identity not parsed: from_user_id=%q from_parent_user_id=%q", m.FromUserId, m.FromParentUserId)
	}
}

// Send response must expose contacts[].user_id for BSUID sends; phone sends
// omit it (absence is not a failure).
func TestSendResponseUserIdParse(t *testing.T) {
	raw := `{"messaging_product":"whatsapp","contacts":[{"input":"IN.123","user_id":"IN.123"}],"messages":[{"id":"wamid.Z"}]}`
	var r MessageSendResponse
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(r.Contacts) != 1 || r.Contacts[0].UserId != "IN.123" {
		t.Fatalf("contacts[].user_id not parsed: %+v", r.Contacts)
	}

	phoneRaw := `{"messaging_product":"whatsapp","contacts":[{"input":"919999999999","wa_id":"919999999999"}],"messages":[{"id":"wamid.P"}]}`
	var pr MessageSendResponse
	if err := json.Unmarshal([]byte(phoneRaw), &pr); err != nil {
		t.Fatalf("unmarshal phone response: %v", err)
	}
	if pr.Contacts[0].WaID != "919999999999" || pr.Contacts[0].UserId != "" {
		t.Fatalf("phone response parsed wrong: %+v", pr.Contacts[0])
	}
}

// Business-username webhook value parses under the correct field shape.
func TestBusinessUsernameUpdateValueParse(t *testing.T) {
	raw := `{"metadata":{"display_phone_number":"918788920939","phone_number_id":"1114345261760613"},"username":"acme.store","previous_username":"acme","event":"updated"}`
	var v BusinessUsernameUpdateValue
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v.Username != "acme.store" || v.Metadata.PhoneNumberId != "1114345261760613" {
		t.Fatalf("username update not parsed: %+v", v)
	}
}
