package order

type StatusOrder string

const (
	StatusOrderPaid   StatusOrder = "PAID"
	StatusOrderFailed StatusOrder = "FAILED"
)

type Orders struct {
	ID     int    `json:"id,omitempty"`
	UserID int    `json:"user_id,omitempty"`
	Status string `json:"status,omitempty"`
}

type OrderItems struct {
	ID      int `json:"id,omitempty"`
	UserID  int `json:"user_id,omitempty"`
	ItemID  int `json:"item_id,omitempty"`
	OrderID int `json:"order_id,omitempty"`
	Qty     int `json:"qty,omitempty"`
}

type MasterItemStocks struct {
	ID         int `json:"id,omitempty"`
	ItemID     int `json:"itemID,omitempty"`
	StockQty   int `json:"stock_qty,omitempty"`
	DesiredQty int
}

type Item struct {
	ItemID int `json:"item_id,omitempty"`
	Qty    int `json:"qty,omitempty"`
}

type CheckoutOrderRequest struct {
	UserID int    `json:"user_id,omitempty"`
	Items  []Item `json:"items,omitempty"`
}

type ReservedItem struct {
	UserID int
	ItemID int
	Qty    int
}

type CreateOrderRequest struct {
	UserID int         `json:"user_id,omitempty"`
	Status StatusOrder `json:"status,omitempty"`
}
