package models

type Portfolio struct {
	ID          int      `json:"portfolio_id" bson:"_id"`
	ProfileID   int      `json:"profile_id" bson:"profile_id"`
	Name        string   `json:"name" bson:"name, omitempty"`
	Category    Category `json:"category" bson:"category, omitempty"`
	Description string   `json:"description" bson:"description, omitempty"`
	Crafts      []Craft  `json:"crafts" bson:"crafts, omitempty"`
}

type Category struct {
	ID   int    `json:"category_id" bson:"_id"`
	Name string `json:"category_name" bson:"category_name"`
}

type Craft struct {
	ID          int       `json:"craft_id" bson:"_id"`
	Name        string    `json:"craft_name" bson:"craft_name"`
	Tags        []Tag     `json:"tags" bson:"tags, omitempty"`
	Description string    `json:"craft_description" bson:"craft_description, omitempty"`
	Contents    []Content `json:"contents" bson:"contents"`
}

type Tag struct {
	ID   int    `json:"tag_id" bson:"_id"`
	Name string `json:"tag_name" bson:"tag_name"`
}

type Content struct {
	ID          int    `json:"content_id" bson:"_id"`
	Description string `json:"content_description" bson:"content_description, omitempty"`
	Data        []byte `json:"data" bson:"data"`
}
