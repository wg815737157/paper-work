package model

type SysTree struct {
	Id      int
	SysId   int
	SysName string `gorm:"sys_name"`
	Tree    string
}
