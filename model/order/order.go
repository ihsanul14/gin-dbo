package order

type Order struct {
	Id         string `json:"id" gorm:"id;primaryKey;uniqueIndex"`
	CustomerId string `json:"customer_id" gorm:"customer_id"`
	Name       string `json:"name,omitempty" gorm:"name"`
	Qty        int64  `json:"qty,omitempty" gorm:"qty"`
	CreatedAt  string `json:"createdAt" gorm:"createdAt"`
	UpdatedAt  string `json:"updatedAt" gorm:"updatedAt"`
}
