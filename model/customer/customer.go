package customer

type Customer struct {
	Id        string `json:"id" gorm:"id;primaryKey;uniqueIndex"`
	Name      string `json:"name,omitempty" gorm:"name"`
	CreatedAt string `json:"createdAt" gorm:"createdAt"`
	UpdatedAt string `json:"updatedAt" gorm:"updatedAt"`
}
