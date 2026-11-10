package business

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wapikit/wapi.go/internal/request_client"
)

func newAnalyticsClient(t *testing.T, response string) (*BusinessClient, *string) {
	t.Helper()
	fields := new(string)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*fields = r.URL.Query().Get("fields")
		_, _ = io.WriteString(w, response)
	}))
	t.Cleanup(srv.Close)
	return NewBusinessClient(&BusinessClientConfig{
		BusinessAccountId: "waba",
		AccessToken:       "t",
		Requester:         request_client.NewRequestClientWithOptions("t", request_client.Options{BaseURL: srv.URL}),
	}), fields
}

func TestFetchAnalyticsSendsJSONListsAndDecodesNestedResult(t *testing.T) {
	client, fields := newAnalyticsClient(t, `{"analytics":{"phone_numbers":["1"],"granularity":"DAY","data_points":[{"start":1,"end":2,"sent":5,"delivered":4}]},"id":"waba"}`)
	res, err := client.FetchAnalytics(AccountAnalyticsOptions{
		Start:        time.Unix(1, 0),
		End:          time.Unix(2, 0),
		Granularity:  AnalyticsRequestGranularityTypeDay,
		PhoneNumbers: []string{"1", "2"},
	})
	if err != nil {
		t.Fatalf("FetchAnalytics: %v", err)
	}
	want := `analytics.country_codes([]).end(2).granularity(DAY).phone_numbers(["1","2"]).start(1)`
	if *fields != want {
		t.Fatalf("fields = %s\nwant     %s", *fields, want)
	}
	if len(res.DataPoints) != 1 || res.DataPoints[0].Sent != 5 || res.Granularity != "DAY" {
		t.Fatalf("res = %+v", res)
	}
}

const conversationAnalyticsResponse = `{"conversation_analytics":{"data":[{"data_points":[
	{"start":1,"end":2,"conversation":3,"conversation_category":"MARKETING","cost":0.75},
	{"start":2,"end":3,"conversation":1,"cost":1.5}]}]},"id":"waba"}`

func TestFetchConversationAnalyticsUsesMetaFilterNamesAndKeepsCost(t *testing.T) {
	client, fields := newAnalyticsClient(t, conversationAnalyticsResponse)
	res, err := client.FetchConversationAnalytics(ConversationAnalyticsOptions{
		Start:                time.Unix(1, 0),
		End:                  time.Unix(3, 0),
		Granularity:          ConversationAnalyticsGranularityTypeDay,
		ConversationCategory: []ConversationCategoryType{ConversationCategoryTypeMarketing},
		Dimensions:           []ConversationDimensionType{ConversationDimensionTypeConversationCategory},
	})
	if err != nil {
		t.Fatalf("FetchConversationAnalytics: %v", err)
	}
	want := `conversation_analytics.conversation_categories(["MARKETING"]).dimensions(["CONVERSATION_CATEGORY"]).end(3).granularity(DAILY).phone_numbers([]).start(1)`
	if *fields != want {
		t.Fatalf("fields = %s\nwant     %s", *fields, want)
	}
	if len(res.Data) != 1 || len(res.Data[0].DataPoints) != 2 || res.Data[0].DataPoints[0].Cost != 0.75 {
		t.Fatalf("res = %+v", res)
	}
}

func TestDeprecatedConversationAnalyticsNowDecodes(t *testing.T) {
	client, _ := newAnalyticsClient(t, conversationAnalyticsResponse)
	res, err := client.ConversationAnalytics(ConversationAnalyticsOptions{Start: time.Unix(1, 0), End: time.Unix(3, 0), Granularity: ConversationAnalyticsGranularityTypeDay})
	if err != nil {
		t.Fatalf("ConversationAnalytics: %v", err)
	}
	points := res.ConversationAnalytics[0].Data[0].DataPoints
	if len(points) != 2 || points[0].Conversation != 3 || points[1].Cost != 2 {
		t.Fatalf("points = %+v", points)
	}
}
