package events

// UserPreferencesEvent represents an event for user marketing preferences update
type UserPreferencesEvent struct {
	BaseBusinessAccountEvent `json:",inline"`
	UserPreferences          []UserPreference `json:"user_preferences"`
}

type UserPreference struct {
	WaId      string `json:"wa_id"`
	Detail    string `json:"detail"`
	Category  string `json:"category"` // e.g., "marketing_messages"
	Value     string `json:"value"`    // Preference value
	Timestamp int64  `json:"timestamp"`
}

// NewUserPreferencesEvent creates a new UserPreferencesEvent instance
func NewUserPreferencesEvent(baseEvent *BaseBusinessAccountEvent, preferences []UserPreference) *UserPreferencesEvent {
	return &UserPreferencesEvent{
		BaseBusinessAccountEvent: *baseEvent,
		UserPreferences:          preferences,
	}
}

func (e *UserPreferencesEvent) GetEventType() string {
	return string(UserPreferencesEventType)
}
