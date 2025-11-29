package events

// PaymentConfigurationUpdateEvent represents an event for payment configuration update
type PaymentConfigurationUpdateEvent struct {
	BaseBusinessAccountEvent `json:",inline"`
	ConfigurationName        string `json:"configuration_name"`
	ProviderName             string `json:"provider_name"`
	ProviderMid              string `json:"provider_mid"`
	Status                   string `json:"status"`
	CreatedTimestamp         int64  `json:"created_timestamp"`
	UpdatedTimestamp         int64  `json:"updated_timestamp"`
}

// NewPaymentConfigurationUpdateEvent creates a new PaymentConfigurationUpdateEvent instance
func NewPaymentConfigurationUpdateEvent(
	baseEvent *BaseBusinessAccountEvent,
	configName string,
	providerName string,
	providerMid string,
	status string,
	createdTs int64,
	updatedTs int64,
) *PaymentConfigurationUpdateEvent {
	return &PaymentConfigurationUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		ConfigurationName:        configName,
		ProviderName:             providerName,
		ProviderMid:              providerMid,
		Status:                   status,
		CreatedTimestamp:         createdTs,
		UpdatedTimestamp:         updatedTs,
	}
}

func (e *PaymentConfigurationUpdateEvent) GetEventType() string {
	return string(PaymentConfigurationUpdateEventType)
}
