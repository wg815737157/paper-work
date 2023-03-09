package model

type SysTree struct {
	Id      int
	SysId   int
	SysName string `gorm:"syw_name"`
	Tree    string
}
