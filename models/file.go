package models

type File struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Filename  string
	Extension string
}
