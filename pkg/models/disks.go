package models

// ListMilestonesOptions is provided when listing Labels in a repository
type ListDisksOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Project string `json:"project"`

	// Owner is the owner of the repository (ex: grafana)
	Mongo string `json:"mongo"`
}