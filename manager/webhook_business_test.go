package manager

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/wapikit/wapi.go/pkg/events"
)

// captureEvents subscribes to the given types and records what arrives.
func captureEvents(t *testing.T, wh *WebhookManager, types ...events.EventType) func(events.EventType) []events.BaseEvent {
	t.Helper()
	var mu sync.Mutex
	got := map[events.EventType][]events.BaseEvent{}
	for _, eventType := range types {
		eventType := eventType
		wh.EventManager.On(eventType, func(event events.BaseEvent) {
			mu.Lock()
			got[eventType] = append(got[eventType], event)
			mu.Unlock()
		})
	}
	return func(eventType events.EventType) []events.BaseEvent {
		mu.Lock()
		defer mu.Unlock()
		return append([]events.BaseEvent(nil), got[eventType]...)
	}
}

func businessPayload(field, value string) string {
	return fmt.Sprintf(`{"object":"whatsapp_business_account","entry":[{"id":"waba-1","time":1700000000,"changes":[{"field":%q,"value":%s}]}]}`, field, value)
}

func postAndExpect200(t *testing.T, wh *WebhookManager, body string) {
	t.Helper()
	rec, err := postWebhook(wh, body, nil)
	if err != nil || rec.Code != http.StatusOK {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestBusinessWebhooksUseMetaFieldNamesAndDedicatedEventTypes(t *testing.T) {
	cases := []struct {
		field     string
		value     string
		eventType events.EventType
	}{
		{"message_template_status_update", `{"event":"APPROVED","message_template_id":1,"message_template_name":"n","message_template_language":"en","reason":"NONE"}`, events.MessageTemplateStatusUpdateEventType},
		{"message_template_quality_update", `{"previous_quality_score":"GREEN","new_quality_score":"RED","message_template_id":1,"message_template_name":"n","message_template_language":"en"}`, events.MessageTemplateQualityUpdateEventType},
		{"template_category_update", `{"message_template_id":1,"message_template_name":"n","message_template_language":"en","previous_category":"UTILITY","new_category":"MARKETING"}`, events.TemplateCategoryUpdateEventType},
		{"phone_number_name_update", `{"display_phone_number":"1","decision":"APPROVED","requested_verified_name":"Acme","rejection_reason":"NONE"}`, events.PhoneNumberNameUpdateEventType},
		{"phone_number_quality_update", `{"display_phone_number":"1","event":"UPGRADE","current_limit":"TIER_10K"}`, events.PhoneNumberQualityUpdateEventType},
		{"account_review_update", `{"decision":"APPROVED"}`, events.AccountReviewUpdateEventType},
		{"business_capability_update", `{"max_daily_conversation_per_phone":1000,"max_phone_numbers_per_business":2}`, events.BusinessCapabilityUpdateEventType},
		{"security", `{"display_phone_number":"1","event":"PIN_CHANGED","requester":"u"}`, events.SecurityEventType},
		// Spellings older releases matched on still work.
		{"message_template_status", `{"event":"REJECTED","message_template_id":1,"message_template_name":"n","message_template_language":"en","reason":"SCAM"}`, events.MessageTemplateStatusUpdateEventType},
		{"account_review", `{"decision":"REJECTED"}`, events.AccountReviewUpdateEventType},
	}
	for _, tc := range cases {
		t.Run(tc.field, func(t *testing.T) {
			wh := newTestWebhook(t, "")
			events_ := captureEvents(t, wh, tc.eventType, events.AccountAlertsEventType)
			postAndExpect200(t, wh, businessPayload(tc.field, tc.value))
			waitFor(t, string(tc.eventType), func() bool { return len(events_(tc.eventType)) == 1 })
			if alerts := events_(events.AccountAlertsEventType); len(alerts) != 0 {
				t.Fatalf("%s was also published as account_alerts", tc.field)
			}
		})
	}
}

func TestTemplateStatusUpdateCarriesOptionalInfo(t *testing.T) {
	wh := newTestWebhook(t, "")
	got := captureEvents(t, wh, events.MessageTemplateStatusUpdateEventType)
	postAndExpect200(t, wh, businessPayload("message_template_status_update", `{
		"event":"REJECTED","message_template_id":7,"message_template_name":"n","message_template_language":"en",
		"reason":"INVALID_FORMAT","message_template_category":"UTILITY",
		"disable_info":{"disable_date":1700000000},
		"other_info":{"title":"FIRST_PAUSE","description":"paused"},
		"rejection_info":{"reason":"bad var","recommendation":"fix it"}}`))
	waitFor(t, "status event", func() bool { return len(got(events.MessageTemplateStatusUpdateEventType)) == 1 })

	event := got(events.MessageTemplateStatusUpdateEventType)[0].(*events.MessageTemplateStatusUpdateEvent)
	if event.BusinessAccountId != "waba-1" || event.MessageTemplateId != 7 || event.MessageTemplateCategory != "UTILITY" {
		t.Fatalf("event = %+v", event)
	}
	if event.DisableInfo == nil || event.DisableInfo.DisableDate != "1700000000" {
		t.Fatalf("DisableInfo = %+v", event.DisableInfo)
	}
	if event.OtherInfo == nil || event.OtherInfo.Title != "FIRST_PAUSE" || event.OtherInfo.Description != "paused" {
		t.Fatalf("OtherInfo = %+v", event.OtherInfo)
	}
	if event.RejectionInfo == nil || event.RejectionInfo.Recommendation != "fix it" {
		t.Fatalf("RejectionInfo = %+v", event.RejectionInfo)
	}
}

func TestTemplateCategoryAndPhoneQualityCarryExtraFields(t *testing.T) {
	wh := newTestWebhook(t, "")
	got := captureEvents(t, wh, events.TemplateCategoryUpdateEventType, events.PhoneNumberQualityUpdateEventType)
	postAndExpect200(t, wh, businessPayload("template_category_update", `{"message_template_id":1,"message_template_name":"n","message_template_language":"en","previous_category":"UTILITY","new_category":"MARKETING","correct_category":"MARKETING"}`))
	postAndExpect200(t, wh, businessPayload("phone_number_quality_update", `{"display_phone_number":"1","event":"DOWNGRADE","old_limit":"TIER_10K","current_limit":"TIER_1K","max_daily_conversations_per_business":"1K"}`))
	waitFor(t, "both events", func() bool {
		return len(got(events.TemplateCategoryUpdateEventType)) == 1 && len(got(events.PhoneNumberQualityUpdateEventType)) == 1
	})

	category := got(events.TemplateCategoryUpdateEventType)[0].(*events.TemplateCategoryUpdateEvent)
	if category.CorrectCategory != "MARKETING" {
		t.Fatalf("CorrectCategory = %q", category.CorrectCategory)
	}
	quality := got(events.PhoneNumberQualityUpdateEventType)[0].(*events.PhoneNumberQualityUpdateEvent)
	if quality.OldLimit != "TIER_10K" || quality.CurrentLimit != "TIER_1K" || quality.MaxDailyConversationsPerBusiness != "1K" {
		t.Fatalf("quality = %+v", quality)
	}
}

func TestSecurityEventCarriesPayload(t *testing.T) {
	wh := newTestWebhook(t, "")
	got := captureEvents(t, wh, events.SecurityEventType)
	postAndExpect200(t, wh, businessPayload("security", `{"display_phone_number":"15550001111","event":"PIN_CHANGED","requester":"12345"}`))
	waitFor(t, "security event", func() bool { return len(got(events.SecurityEventType)) == 1 })
	event := got(events.SecurityEventType)[0].(*events.SecurityEvent)
	if event.BusinessAccountId != "waba-1" || event.Event != "PIN_CHANGED" || event.Requester != "12345" || event.DisplayPhoneNumber != "15550001111" {
		t.Fatalf("event = %+v", event)
	}
}

func TestAccountRestrictionListWithNumericExpirationIsDecoded(t *testing.T) {
	wh := newTestWebhook(t, "")
	got := captureEvents(t, wh, events.AccountUpdateEventType)
	postAndExpect200(t, wh, businessPayload("account_update", `{"event":"ACCOUNT_RESTRICTION","restriction_info":[
		{"restriction_type":"RESTRICTED_BIZ_INITIATED_MESSAGING","expiration":1641330498},
		{"restriction_type":"RESTRICTED_ADD_PHONE_NUMBER_ACTION","expiration":1641330499,"remediation":"appeal"}]}`))
	waitFor(t, "account update", func() bool { return len(got(events.AccountUpdateEventType)) == 1 })

	event := got(events.AccountUpdateEventType)[0].(*events.AccountUpdateEvent)
	if len(event.RestrictionInfoList) != 2 {
		t.Fatalf("RestrictionInfoList = %+v", event.RestrictionInfoList)
	}
	if event.RestrictionInfo == nil || event.RestrictionInfo.RestrictionType != "RESTRICTED_BIZ_INITIATED_MESSAGING" || event.RestrictionInfo.Expiration != "1641330498" {
		t.Fatalf("RestrictionInfo = %+v", event.RestrictionInfo)
	}
	if second := event.RestrictionInfoList[1]; second.Expiration != "1641330499" || second.Remediation != "appeal" {
		t.Fatalf("second = %+v", second)
	}
}

func TestAccountRestrictionSingleObjectStillDecodes(t *testing.T) {
	var value AccountUpdateValue
	if err := value.UnmarshalJSON([]byte(`{"event":"ACCOUNT_RESTRICTION","phone_number":"1","restriction_info":{"restriction_type":"X","expiration":"2026-04-19T13:59:20+0000"}}`)); err != nil {
		t.Fatal(err)
	}
	if value.PhoneNumber != "1" || value.RestrictionInfo == nil || value.RestrictionInfo.Expiration != "2026-04-19T13:59:20+0000" || len(value.RestrictionInfoList) != 1 {
		t.Fatalf("value = %+v", value)
	}
}

func TestSmbMessageEchoContentIsSurfaced(t *testing.T) {
	wh := newTestWebhook(t, "")
	got := captureEvents(t, wh, events.SmbMessageEchoesEventType)
	postAndExpect200(t, wh, businessPayload("smb_message_echoes", `{"messaging_product":"whatsapp","metadata":{"display_phone_number":"1","phone_number_id":"pn"},
		"message_echoes":[{"from":"1","to":"2","id":"wamid.1","timestamp":"3","type":"text","text":{"body":"hi"}}]}`))
	waitFor(t, "echo event", func() bool { return len(got(events.SmbMessageEchoesEventType)) == 1 })
	event := got(events.SmbMessageEchoesEventType)[0].(*events.SmbMessageEchoesEvent)
	if len(event.MessageEchoes) != 1 || event.MessageEchoes[0].Content["body"] != "hi" {
		t.Fatalf("echoes = %+v", event.MessageEchoes)
	}
}

func TestUnknownWebhookFieldIsAcknowledged(t *testing.T) {
	postAndExpect200(t, newTestWebhook(t, ""), businessPayload("some_future_field", `{"x":1}`))
}
