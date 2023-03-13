package model

type SysTree struct {
	Id      int
	SysId   int    `gorm:"sys_id"`
	SysName string `gorm:"sys_name"`
	Tree    string `gorm:"tree"`
}
