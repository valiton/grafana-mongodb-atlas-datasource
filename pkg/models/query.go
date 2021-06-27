package models

const (
	// QueryTypeCommits is sent by the frontend when querying commits in a GitHub repository
	QueryDiskMeasurements = "disk_measurements"
	// QueryTypeIssues is used when querying issues in a GitHub repository
	QueryDatabaseMeasurements = "database_measurements"
	// QueryTypeContributors is used when querying contributors in a GitHub repository
	QueryProcessMeasurements = "process_measurements"
)

type QueryOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type DiskMeasurementsQuery struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Project QueryOption `json:"project"`

	Cluster QueryOption `json:"cluster"`

	// Owner is the owner of the repository (ex: grafana)
	Mongo QueryOption `json:"mongo"`

	// Owner is the owner of the repository (ex: grafana)
	Disk QueryOption `json:"disk"`

	Dimension QueryOption `json:"dimension"`

	Alias string `json:"alias"`

	RefId string `json:"refId"`
}

type ProcessMeasurementsQuery struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Project QueryOption `json:"project"`

	Cluster QueryOption `json:"cluster"`

	// Owner is the owner of the repository (ex: grafana)
	Mongo QueryOption `json:"mongo"`

	Dimension QueryOption `json:"dimension"`

	Alias string `json:"alias"`

	RefId string `json:"refId"`
}

type DatabaseMeasurementsQuery struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Project QueryOption `json:"project"`

	Cluster QueryOption `json:"cluster"`

	// Owner is the owner of the repository (ex: grafana)
	Mongo QueryOption `json:"mongo"`

	// Owner is the owner of the repository (ex: grafana)
	Database QueryOption `json:"database"`

	Dimension QueryOption `json:"dimension"`

	Alias string `json:"alias"`

	RefId string `json:"refId"`
}

type GetProcessMeasurementsQuery struct {
	Project string
	Mongo string
	Measurement string
	IntervalMs int64
	Start string
	End string
}

type GetDiskMeasurementsQuery struct {
	Project string
	Mongo string
	Disk string
	Measurement string
	IntervalMs int64
	Start string
	End string
}

type GetDatabaseMeasurementsQuery struct {
	Project string
	Mongo string
	Database string
	Measurement string
	IntervalMs int64
	Start string
	End string
}
