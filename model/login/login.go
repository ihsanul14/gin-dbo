package login

type User struct {
	Username   string `json:"username" gorm:"username;primaryKey;uniqueIndex"`
	Password   string `json:"password,omitempty" gorm:"password" swaggerignore:"true"`
	Role       string `json:"role" gorm:"role"`
	CustomerId string `json:"customerId,omitempty" gorm:"customer_id"`
	CreatedAt  string `json:"createdAt" gorm:"createdAt"`
	UpdatedAt  string `json:"updatedAt" gorm:"updatedAt"`
}
