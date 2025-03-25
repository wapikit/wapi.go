package events

type AccountUpdateEventEnum string

const (
	AccountUpdateEventEnumVerifiedAccount    AccountUpdateEventEnum = "VERIFIED_ACCOUNT"
	AccountUpdateEventEnumDisabledAccount    AccountUpdateEventEnum = "DISABLED_UPDATE"
	AccountUpdateEventEnumAccountViolation   AccountUpdateEventEnum = "ACCOUNT_VIOLATION"
	AccountUpdateEventEnumAccountRestriction AccountUpdateEventEnum = "ACCOUNT_RESTRICTION"
	AccountUpdateEventEnumAccountDeleted     AccountUpdateEventEnum = "ACCOUNT_DELETED"
	AccountUpdateEventEnumPartnerRemoved     AccountUpdateEventEnum = "PARTNER_REMOVED"
)

type AccountUpdateEvent struct {
	BaseBusinessAccountEvent
	StatusUpdate    AccountUpdateEventEnum
	PhoneNumber     string
	BanInfo         *BanInfo
	ViolationInfo   *ViolationInfo
	RestrictionInfo []RestrictionInfo
}

type BanInfo struct {
	WabaBanState string
	WabaBanDate  string
}

type ViolationInfo struct {
	ViolationType string
}

type RestrictionInfo struct {
	RestrictionType string
	Expiration      string
}

func NewAccountUpdateEvent(baseEvent *BaseBusinessAccountEvent, statusUpdate AccountUpdateEventEnum, phoneNumber string) *AccountUpdateEvent {
	return &AccountUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		StatusUpdate:             statusUpdate,
		PhoneNumber:              phoneNumber,
	}
}
