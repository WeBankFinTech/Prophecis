package models

type Group struct {
	ID         int64  `gorm:"column:id; PRIMARY_KEY" json:"id"`
	EnableFlag int8   `json:"enable_flag"`
	GroupId    int64  `json:"group_id"`
	Name       string `json:"name"`
	Remarks    string `json:"remarks"`
}
