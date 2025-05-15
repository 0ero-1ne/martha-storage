package models

type File struct {
	Id        uint `gorm:"primaryKey"`
	Filename  string
	Extension string
	Filetype  Filetype
}

type Filetype int8

const (
	Audio Filetype = iota
	Chapter
	Image
)
