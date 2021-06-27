package models

// ListMilestonesOptions is provided when listing Labels in a repository
type ListClustersOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Project string `json:"project"`
}