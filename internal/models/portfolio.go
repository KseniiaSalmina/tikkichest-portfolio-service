package models

type Portfolio struct {
	ID          string
	ProfileID   int
	Category    string
	Description string
	Crafts      []Craft
}

type Craft struct {
	ID          string
	Name        string
	Tags        []string
	Description string
	Contents    []Content
}

type Content struct {
	ID          string
	Description string
	Data        interface{}
}
