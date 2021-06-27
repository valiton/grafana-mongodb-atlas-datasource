package datasource

import (
	"context"
	"time"
	"strings"

	simplejson "github.com/bitly/go-simplejson"

	"github.com/valiton/mongodbatlas-datasource/pkg/dfutil"
	"github.com/valiton/mongodbatlas-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)



// Datasource handles requests to GitHub
type Datasource struct {
	client *MongoDBAtlasClient
}

// CheckHealth calls frequently used endpoints to determine if the client has sufficient privileges
func (d *Datasource) CheckHealth(ctx context.Context) error {
	_, err := d.client.query(ctx, "/", nil)

	return err
}

func unixToRFC3339(date int64) string {
	return time.Unix(date, 0).UTC().Format(time.RFC3339)
}

func formatInterval(interval time.Duration) int64 {
	return int64(interval / time.Millisecond)
}

type MeasurementOptions struct {
	Start       string
	End         string
	IntervalMs  int64
	Measurement string
}

type DataPoint struct {
	Timestamp string
	Value     float64
}

type DataFrame struct {
	Name string
	Points []DataPoint
}

func (df DataFrame) Frames() data.Frames {
	frame := data.NewFrame(
		"logs",
		data.NewField("time", nil, []time.Time{}),
		data.NewField(df.Name, nil, []float64{}),
	)

	for _, v := range df.Points {
		timestamp, _ := time.Parse(time.RFC3339, v.Timestamp)
		frame.AppendRow(
			timestamp,
			v.Value,
		)
	}

	return data.Frames{frame}
}


func GetMeasurementOptions(options *MeasurementOptions) map[string]string {
	var granularity string

	if options.IntervalMs <= 60000 {
		granularity = "PT1M"
	} else if options.IntervalMs <= 500000 {
		granularity = "PT5M"
	} else if options.IntervalMs <= 3600000 {
		granularity = "PT1H"
	} else {
		granularity = "PT1D"
	}

	return map[string]string{
		"start":       options.Start,
		"end":         options.End,
		"m":           options.Measurement,
		"granularity": granularity,
	}
}

func GetMeasurements(name string, body []byte, ctx context.Context) (*DataFrame, error) {
	jBody, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	df := &DataFrame{
		Name: name,
		Points: make([]DataPoint, 0),
	}

	jMeasurements := jBody.Get("measurements")
	log.DefaultLogger.Debug("GetMeasurements", "measurements", jMeasurements)
	if len(jMeasurements.MustArray()) == 0 {
		return df, nil
	}
	firstMeasurement := jMeasurements.GetIndex(0)
	log.DefaultLogger.Debug("GetMeasurements", "first measurement", firstMeasurement)
	var rawDataPoints = firstMeasurement.Get("dataPoints")
	log.DefaultLogger.Debug("GetMeasurements", "raw data points", rawDataPoints)
	var numDataPoints = len(rawDataPoints.MustArray())
	var dataPoints = make([]DataPoint, 0, numDataPoints)
	for i := 0; i < numDataPoints; i++ {
		var jDataPoint = rawDataPoints.GetIndex(i)

		dataPointValue := jDataPoint.Get("value")

		// filter out all JSON null values that are sent by Atlas
		if dataPointValue.Interface() == nil {
			continue
		}

		log.DefaultLogger.Debug("GetMeasurements", "data point", jDataPoint)

		dataPoint := DataPoint{
			Timestamp: jDataPoint.Get("timestamp").MustString(),
			Value:     dataPointValue.MustFloat64(),
		}
		dataPoints = append(dataPoints, dataPoint)
	}

	df.Points = dataPoints;

	log.DefaultLogger.Debug("GetMeasurements", "Final data points", dataPoints)

	return df, nil
}

func getName(refId, alias, projectID, projectName, clusterID, clusterName, database, mongo, disk, dimensionName string) string {
	name := refId
	if alias != "" {
		name = alias
		name = strings.ReplaceAll(name, "{{projectId}}", projectID)
		name = strings.ReplaceAll(name, "{{projectName}}", projectName)
		name = strings.ReplaceAll(name, "{{clusterId}}", clusterID)
		name = strings.ReplaceAll(name, "{{clusterName}}", clusterName)
		name = strings.ReplaceAll(name, "{{database}}", database)
		name = strings.ReplaceAll(name, "{{mongo}}", mongo)
		name = strings.ReplaceAll(name, "{{disk}}", disk)
		name = strings.ReplaceAll(name, "{{dimensionName}}", dimensionName)
	}
	return name
}

// HandleRepositoriesQuery is the query handler for listing GitHub Repositories
func (d *Datasource) HandleProcessMeasurementsQuery(ctx context.Context, query *models.ProcessMeasurementsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	project := query.Project.Value
	mongo := query.Mongo.Value
	options := &MeasurementOptions{
		Start: unixToRFC3339(req.TimeRange.From.Unix()),
		End: unixToRFC3339(req.TimeRange.To.Unix()),
		IntervalMs: formatInterval(req.Interval),
		Measurement: query.Dimension.Value,
	}
	
	name := getName(query.RefId, query.Alias, query.Project.Value, query.Project.Label, query.Cluster.Value, query.Cluster.Label, "", mongo, "", query.Dimension.Value)

	body, err := d.client.query(ctx, "/groups/"+project+"/processes/"+mongo+"/measurements", GetMeasurementOptions(options))
	if err != nil {
		return nil, err
	}
	return GetMeasurements(name, body, ctx)
}

func (d *Datasource) HandleDatabaseMeasurementsQuery(ctx context.Context, query *models.DatabaseMeasurementsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	project := query.Project.Value
	mongo := query.Mongo.Value
	database := query.Database.Value
	options := &MeasurementOptions{
		Start: unixToRFC3339(req.TimeRange.From.Unix()),
		End: unixToRFC3339(req.TimeRange.To.Unix()),
		IntervalMs: formatInterval(req.Interval),
		Measurement: query.Dimension.Value,
	}

	name := getName(query.RefId, query.Alias, query.Project.Value, query.Project.Label, query.Cluster.Value, query.Cluster.Label, database, mongo, "", query.Dimension.Value)

	body, err := d.client.query(ctx, "/groups/"+project+"/processes/"+mongo+"/databases/"+database+"/measurements", GetMeasurementOptions(options))
	if err != nil {
		return nil, err
	}
	return GetMeasurements(name, body, ctx)
}

func (d *Datasource) HandleDiskMeasurementsQuery(ctx context.Context, query *models.DiskMeasurementsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	project := query.Project.Value
	mongo := query.Mongo.Value
	disk := query.Disk.Value
	options := &MeasurementOptions{
		Start: unixToRFC3339(req.TimeRange.From.Unix()),
		End: unixToRFC3339(req.TimeRange.To.Unix()),
		IntervalMs: formatInterval(req.Interval),
		Measurement: query.Dimension.Value,
	}

	name := getName(query.RefId, query.Alias, query.Project.Value, query.Project.Label, query.Cluster.Value, query.Cluster.Label, "", mongo, disk, query.Dimension.Value)

	body, err := d.client.query(ctx, "/groups/"+project+"/processes/"+mongo+"/disks/"+disk+"/measurements", GetMeasurementOptions(options))
	if err != nil {
		return nil, err
	}
	return GetMeasurements(name, body, ctx)
}

// NewDatasource creates a new datasource for handling queries
func NewDatasource(ctx context.Context, settings *models.Settings) *Datasource {
	client := &MongoDBAtlasClient{
		settings: settings,
	}
	return &Datasource{
		client: client,
	}
}