package models

const (
	QueryDiskMeasurements = "disk_measurements"
	QueryDatabaseMeasurements = "database_measurements"
	QueryProcessMeasurements = "process_measurements"
)

type QueryOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type DiskMeasurementsQuery struct {
	Project QueryOption `json:"project"`

	Cluster QueryOption `json:"cluster"`

	Mongo QueryOption `json:"mongo"`

	Disk QueryOption `json:"disk"`

	Dimension QueryOption `json:"dimension"`

	Alias string `json:"alias"`

	RefId string `json:"refId"`
}

type ProcessMeasurementsQuery struct {
	Project QueryOption `json:"project"`

	Cluster QueryOption `json:"cluster"`

	Mongo QueryOption `json:"mongo"`

	Dimension QueryOption `json:"dimension"`

	Alias string `json:"alias"`

	RefId string `json:"refId"`
}

type DatabaseMeasurementsQuery struct {
	Project QueryOption `json:"project"`

	Cluster QueryOption `json:"cluster"`

	Mongo QueryOption `json:"mongo"`

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
