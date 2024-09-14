package models

// ProductLabelReverseResolver returns all projects for a given label
type ProductLabelReverseResolver func(l string) []string

type Product struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	SmallName   string `json:"smallName"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func (p *Product) IsValid() bool {
	return p.Name != "" && p.Price > 0 && p.Description != ""
}
