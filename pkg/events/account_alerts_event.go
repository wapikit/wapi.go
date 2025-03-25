package events

type AccountAlertSeverityEnum string

const (
	AccountAlertSeverityEnumCritical AccountAlertSeverityEnum = "CRITICAL"
	AccountAlertSeverityEnumWarning  AccountAlertSeverityEnum = "WARNING"
)

type AccountAlertStatusEnum string

const (
	AccountAlertStatusEnumActive   AccountAlertStatusEnum = "ACTIVE"
	AccountAlertStatusEnumResolved AccountAlertStatusEnum = "RESOLVED"
)

type AccountAlertEvent struct {
	BaseBusinessAccountEvent
	EntityType       string
	EntityId         string
	AlertSeverity    AccountAlertSeverityEnum
	AlertStatus      AccountAlertStatusEnum
	AlertType        string
	AlertDescription string
}

func NewAccountAlertEvent(baseEvent *BaseBusinessAccountEvent, entityType string, entityId string, alertSeverity AccountAlertSeverityEnum, alertStatus AccountAlertStatusEnum, alertType string, alertDescription string) *AccountAlertEvent {
	return &AccountAlertEvent{
		BaseBusinessAccountEvent: *baseEvent,
		EntityType:               entityType,
		EntityId:                 entityId,
		AlertSeverity:            alertSeverity,
		AlertStatus:              alertStatus,
		AlertType:                alertType,
		AlertDescription:         alertDescription,
	}
}
