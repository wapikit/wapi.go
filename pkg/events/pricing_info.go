package events

// MessagePricingInfo carries Meta's per-message pricing-block fields from the
// `messages.statuses[].pricing` shape of the WhatsApp Business webhook.
// Meta populates this for billable status updates (typically `sent` or
// `delivered` depending on category × country × pricing-model rollout).
//
// Note: Meta does NOT include the per-message amount in this block — the amount
// must be looked up from Meta's published rate card keyed by
// (Category × PricingModel × Country × Currency).
type MessagePricingInfo struct {
	// Billable is whether Meta will charge us for this send. Falsey for
	// user-initiated service-window replies, free entry-points, etc.
	Billable bool `json:"billable"`

	// PricingModel is "CBP" (conversation-based, legacy 24h window) or
	// "PMP" (per-message). Meta is mid-migration across markets in 2025+.
	PricingModel string `json:"pricingModel"`

	// Category is Meta's classification of the message: "MARKETING" |
	// "UTILITY" | "AUTHENTICATION" | "SERVICE" | "REFERRAL_CONVERSION".
	// This is the same category surfaced on the template, but the webhook
	// is the authoritative source — Meta can override at send time.
	Category string `json:"category"`
}
