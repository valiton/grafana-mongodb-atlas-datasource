package models

type ListDatabasesOptions struct {
	Project string `json:"project"`

	Mongo string `json:"mongo"`
}