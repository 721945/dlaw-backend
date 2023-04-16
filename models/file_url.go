package models

type FileUrl struct {
	Base
	Url          string
	PublishedUrl string
	PreviewUrl   string
	FileId       *uint
}
