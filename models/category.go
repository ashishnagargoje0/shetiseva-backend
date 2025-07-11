package models

type Category struct {
	Slug string `bson:"slug" json:"slug"`
	Name string `bson:"name" json:"name"`
}
