package events

type PhoneNumberUpdateEventEnum string

const (
	PhoneNumberUpdateEventEnumDowngrade  PhoneNumberUpdateEventEnum = "DOWNGRADE"
	PhoneNumberUpdateEventEnumFlagged    PhoneNumberUpdateEventEnum = "FLAGGED"
	PhoneNumberUpdateEventEnumOnboarding PhoneNumberUpdateEventEnum = "ONBOARDING"
	PhoneNumberUpdateEventEnumUnflagged  PhoneNumberUpdateEventEnum = "UNFLAGGED"
	PhoneNumberUpdateEventEnumUpgrade    PhoneNumberUpdateEventEnum = "UPGRADE"
)

type PhoneNumberQualityUpdateCurrentLimitEnum string

const (
	PhoneNumberQualityUpdateCurrentLimitEnumTier50        PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_50"
	PhoneNumberQualityUpdateCurrentLimitEnumTier250       PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_250"
	PhoneNumberQualityUpdateCurrentLimitEnumTier1K        PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_1K"
	PhoneNumberQualityUpdateCurrentLimitEnumTier10K       PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_10K"
	PhoneNumberQualityUpdateCurrentLimitEnumTier100K      PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_100K"
	PhoneNumberQualityUpdateCurrentLimitEnumTierUnlimited PhoneNumberQualityUpdateCurrentLimitEnum = "TIER_UNLIMITED"
)

type PhoneNumberQualityUpdateEvent struct {
	BaseBusinessAccountEvent
	DisplayPhoneNumber string
	Event              PhoneNumberUpdateEventEnum
	CurrentLimit       PhoneNumberQualityUpdateCurrentLimitEnum
}

func NewPhoneNumberQualityUpdateEvent(baseEvent *BaseBusinessAccountEvent, displayPhoneNumber string, event PhoneNumberUpdateEventEnum, currentLimit PhoneNumberQualityUpdateCurrentLimitEnum) *PhoneNumberQualityUpdateEvent {
	return &PhoneNumberQualityUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		DisplayPhoneNumber:       displayPhoneNumber,
		Event:                    event,
		CurrentLimit:             currentLimit,
	}
}
