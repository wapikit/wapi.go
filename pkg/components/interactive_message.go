package components

type InteractiveMessageType string

const (
	InteractiveMessageTypeButton          InteractiveMessageType = "button"
	InteractiveMessageTypeProduct         InteractiveMessageType = "product"
	InteractiveMessageTypeProductList     InteractiveMessageType = "product_list"
	InteractiveMessageTypeList            InteractiveMessageType = "list"
	InteractiveMessageTypeCatalog         InteractiveMessageType = "catalog_message"
	InteractiveMessageTypeLocationRequest InteractiveMessageType = "location_request_message"
)
