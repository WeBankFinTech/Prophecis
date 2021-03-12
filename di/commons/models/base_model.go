package models

type BaseModel struct {
	ID int64 `gorm:"column:id; PRIMARY_KEY" json:"id"`
	//EnableFlag *int  `json:"enable_flag"`
	EnableFlag bool `json:"enable_flag"`
}
