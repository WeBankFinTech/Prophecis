package models

type BaseModel struct {
	ID int64 `gorm:"column:id; PRIMARY_KEY" json:"id"`
}
