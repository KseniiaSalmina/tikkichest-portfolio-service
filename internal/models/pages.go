package models

type CategoriesPage struct {
	Categories  []Category `json:"categories"`
	PageNo      int        `json:"page_number"`
	Limit       int        `json:"limit"`
	PagesAmount int        `json:"pages_amount"`
}

type PortfoliosPage struct {
	Portfolios  []Portfolio `json:"portfolios"`
	PageNo      int         `json:"page_number"`
	Limit       int         `json:"limit"`
	PagesAmount int         `json:"pages_amount"`
}

type CraftsPage struct {
	Crafts      []Craft `json:"portfolios"`
	PageNo      int     `json:"page_number"`
	Limit       int     `json:"limit"`
	PagesAmount int     `json:"pages_amount"`
}

type TagsPage struct {
	Tags        []Tag `json:"portfolios"`
	PageNo      int   `json:"page_number"`
	Limit       int   `json:"limit"`
	PagesAmount int   `json:"pages_amount"`
}
