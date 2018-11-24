package models

type Member struct {
	PublicId int64  `json:"id"`
	Name     string `json:"name"`
	Guid     string `json:"guid"`
}
