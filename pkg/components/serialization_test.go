package components

import (
	"encoding/json"
	"strings"
	"testing"
)

func mustJSON(t *testing.T, v interface{}) string {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(raw)
}

func TestContactUsesCountryCodeAndOmitsEmptyOrg(t *testing.T) {
	contact := NewContact(ContactName{FormattedName: "Ada"})
	contact.AddAddress(ContactAddress{CountryCode: "IN", Type: "HOME"})
	got := mustJSON(t, contact)
	if !strings.Contains(got, `"country_code":"IN"`) || strings.Contains(got, "countryCode") {
		t.Fatalf("address not sent as country_code: %s", got)
	}
	if strings.Contains(got, `"org"`) {
		t.Fatalf("empty org was sent: %s", got)
	}

	contact.SetOrg(ContactOrg{Company: "Acme"})
	if got := mustJSON(t, contact); !strings.Contains(got, `"org":{"company":"Acme"}`) {
		t.Fatalf("org missing: %s", got)
	}

	var inbound Contact
	if err := json.Unmarshal([]byte(`{"name":{"formatted_name":"A"},"addresses":[{"country_code":"US","type":"WORK"}]}`), &inbound); err != nil {
		t.Fatal(err)
	}
	if inbound.Addresses[0].CountryCode != "US" {
		t.Fatalf("inbound country_code not decoded: %+v", inbound.Addresses)
	}
}

func TestCtaAndProductListOmitUnsetHeaders(t *testing.T) {
	cta, _ := NewCtaMessage(CtaMessageParams{BodyText: "hi"})
	if got := mustJSON(t, cta); strings.Contains(got, `"header"`) {
		t.Fatalf("empty CTA header was sent: %s", got)
	}
	cta.SetHeader("Title")
	if got := mustJSON(t, cta); !strings.Contains(got, `"header":{"type":"text","text":"Title"}`) {
		t.Fatalf("CTA header missing: %s", got)
	}

	list, err := NewProductListMessage(ProductListMessageParams{CatalogId: "c", ProductRetailerId: "sku", BodyText: "b"})
	if err != nil {
		t.Fatal(err)
	}
	got := mustJSON(t, list)
	if strings.Contains(got, `"header"`) || strings.Contains(got, `"product_retailer_id"`) {
		t.Fatalf("product list sent an empty header or a stray product_retailer_id: %s", got)
	}
	list.SetHeader("Picks")
	if got := mustJSON(t, list); !strings.Contains(got, `"header":{"type":"text","text":"Picks"}`) {
		t.Fatalf("product list header missing: %s", got)
	}
	if _, err := NewProductListMessage(ProductListMessageParams{CatalogId: "c", BodyText: "b"}); err != nil {
		t.Fatalf("ProductRetailerId should be optional: %v", err)
	}
}

func TestNewVideoMessageKeepsCaption(t *testing.T) {
	video, err := NewVideoMessage(VideoMessageConfigs{Link: "https://x/v.mp4", Caption: "look"})
	if err != nil {
		t.Fatal(err)
	}
	if video.Caption != "look" {
		t.Fatalf("caption dropped: %+v", video)
	}
}

func TestListSetBodyUsesText(t *testing.T) {
	list := &listMessage{}
	list.SetBody("Pick one")
	if list.Body.Text != "Pick one" {
		t.Fatalf("body = %q", list.Body.Text)
	}
}

func TestTemplateMediaParameterAcceptsId(t *testing.T) {
	got := mustJSON(t, TemplateMessageParameterMedia{Id: "media-1"})
	if got != `{"id":"media-1"}` {
		t.Fatalf("got %s", got)
	}
	if got := mustJSON(t, TemplateMessageParameterMedia{Link: "https://x"}); got != `{"link":"https://x"}` {
		t.Fatalf("got %s", got)
	}
}
