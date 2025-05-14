package components

// Order represents an order in the system.
type Order struct {
	CatalogID    string        `json:"catalog_id"`    // CatalogID is the ID of the catalog associated with the order.
	ProductItems []ProductItem `json:"product_items"` // ProductItems is a list of product items in the order.
	Text         *string       `json:"text"`          // Text is an additional text associated with the order.
}

// ProductItem represents a product item in an order.
type ProductItem struct {
	Currency          string  `json:"currency"`            // Currency is the currency of the product item.
	ItemPrice         float64 `json:"item_price"`          // ItemPrice is the price of the product item.
	ProductRetailerID string  `json:"product_retailer_id"` // ProductRetailerID is the ID of the retailer associated with the product item.
	Quantity          int     `json:"quantity"`            // Quantity is the quantity of the product item.
}
