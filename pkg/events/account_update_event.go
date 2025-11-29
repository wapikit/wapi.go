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
	StatusUpdate                     AccountUpdateEventEnum
	PhoneNumber                      string
	Country                          string
	WabaInfo                         *WabaInfo
	ViolationInfo                    *ViolationInfo
	AuthInternationalRateEligibility *AuthInternationalRateEligibility
	BanInfo                          *BanInfo
	RestrictionInfo                  *RestrictionInfo
	PartnerClientCertificationInfo   *PartnerClientCertificationInfo
}

type WabaInfo struct {
	WabaId                     string
	OwnerBusinessId            string
	AdAccountLinked            string
	PartnerAppId               string
	SolutionId                 string
	SolutionPartnerBusinessIds []string
}

type ViolationInfo struct {
	ViolationType string
}

type AuthInternationalRateEligibility struct {
	ExceptionCountries []ExceptionCountry
	StartTime          int64
}

type ExceptionCountry struct {
	CountryCode string
	StartTime   int64
}

type BanInfo struct {
	WabaBanState []string
	WabaBanDate  string
}

type RestrictionInfo struct {
	RestrictionType string
	Expiration      string
}

type PartnerClientCertificationInfo struct {
	ClientBusinessId string
	Status           string
	RejectionReasons []string
}

func NewAccountUpdateEvent(
	baseEvent *BaseBusinessAccountEvent,
	statusUpdate AccountUpdateEventEnum,
	phoneNumber string,
	country string,
	wabaInfo *WabaInfo,
	violationInfo *ViolationInfo,
	authInternationalRateEligibility *AuthInternationalRateEligibility,
	banInfo *BanInfo,
	restrictionInfo *RestrictionInfo,
	partnerClientCertificationInfo *PartnerClientCertificationInfo,
) *AccountUpdateEvent {
	return &AccountUpdateEvent{
		BaseBusinessAccountEvent:         *baseEvent,
		StatusUpdate:                     statusUpdate,
		PhoneNumber:                      phoneNumber,
		Country:                          country,
		WabaInfo:                         wabaInfo,
		ViolationInfo:                    violationInfo,
		AuthInternationalRateEligibility: authInternationalRateEligibility,
		BanInfo:                          banInfo,
		RestrictionInfo:                  restrictionInfo,
		PartnerClientCertificationInfo:   partnerClientCertificationInfo,
	}
}
