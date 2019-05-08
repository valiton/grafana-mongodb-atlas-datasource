package main

import (
	"strconv"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/grafana/grafana_plugin_model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
)

var processMetrics map[string]string

type MongoDbAtlasDatasource struct {
	plugin.NetRPCUnsupportedPlugin
	logger hclog.Logger
}

func (ds *MongoDbAtlasDatasource) Query(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	ds.logger.Debug("Query", "datasource", tsdbReq.Datasource.Name, "TimeRange", tsdbReq.TimeRange)

	queryType, err := GetQueryType(tsdbReq)
	if err != nil {
		return nil, err
	}

	ds.logger.Debug("createRequest", "queryType", queryType)

	switch queryType {
	case "test":
		return TestConnectionQuery(ctx, tsdbReq)
	case "search":
		return SearchQuery(ctx, tsdbReq)
	default:
		return MetricQuery(ctx, tsdbReq)
	}
}

func CredentialsFromReq(tsdbReq *datasource.DatasourceRequest) *AtlasCredentials {
	parameters, _ := simplejson.NewJson([]byte(tsdbReq.Datasource.JsonData))

	return &AtlasCredentials{
		Email:    parameters.Get("atlasUsername").MustString(),
		ApiToken: tsdbReq.Datasource.DecryptedSecureJsonData["atlasApiToken"],
	}
}

func TestConnectionQuery(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	credentials := CredentialsFromReq(tsdbReq)
	_, err := MakeHttpRequest(ctx, "/", credentials, nil)
	if err != nil {
		return nil, err
	}

	response := &datasource.DatasourceResponse{}
	return response, nil
}

type suggestData struct {
	Text  string
	Value string
}

// Search for possible metrics and dimensions
func SearchQuery(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	pluginLogger.Debug("SearchQuery", "Intro")
	credentials := CredentialsFromReq(tsdbReq)

	result := &datasource.DatasourceResponse{}
	firstQuery := tsdbReq.Queries[0]
	var secondQuery *datasource.Query
	if len(tsdbReq.Queries) > 1 {
		secondQuery = tsdbReq.Queries[1]
	}
	queryResult := &datasource.QueryResult{RefId: firstQuery.RefId}

	pluginLogger.Debug("SearchQuery", "First query", firstQuery.ModelJson)
	parametersJSON := "{}"
	if secondQuery != nil {
		parametersJSON = secondQuery.ModelJson
	}
	parameters, err := simplejson.NewJson([]byte(parametersJSON))
	if err != nil {
		return nil, err
	}
	pluginLogger.Debug("SearchQuery", "Parameters", parameters)
	// subType := parameters.Get("subtype").MustString()

	var data []*suggestData
	if firstQuery.RefId == "projects" {
		projects, err := GetProjects(ctx, credentials)
		if err != nil {
			return nil, err
		}
		data = make([]*suggestData, len(projects))
		for i, project := range projects {
			data[i] = &suggestData{
				Text:  project.Name,
				Value: project.ID,
			}
		}
	} else if firstQuery.RefId == "clusters" && secondQuery != nil {
		projectID := parameters.Get("projectId").MustString()
		clusters, err := GetClusters(ctx, credentials, projectID)
		if err != nil {
			return nil, err
		}
		data = make([]*suggestData, len(clusters))
		for i, cluster := range clusters {
			data[i] = &suggestData{
				Text:  cluster.Name,
				Value: cluster.ID,
			}
		}
	} else if firstQuery.RefId == "mongos" && secondQuery != nil {
		projectID := parameters.Get("projectId").MustString()
		clusterID := parameters.Get("clusterId").MustString()
		mongos, err := GetMongos(ctx, credentials, projectID, clusterID)
		if err != nil {
			return nil, err
		}
		data = make([]*suggestData, len(mongos))
		for i, mongo := range mongos {
			data[i] = &suggestData{
				Text:  mongo,
				Value: mongo,
			}
		}
	} else if firstQuery.RefId == "disks" && secondQuery != nil {
		projectID := parameters.Get("projectId").MustString()
		mongo := parameters.Get("mongo").MustString()
		disks, err := GetMongoDisks(ctx, credentials, projectID, mongo)
		if err != nil {
			return nil, err
		}
		data = make([]*suggestData, len(disks))
		for i, disk := range disks {
			data[i] = &suggestData{
				Text:  string(disk),
				Value: string(disk),
			}
		}
	} else if firstQuery.RefId == "databases" && secondQuery != nil {
		projectID := parameters.Get("projectId").MustString()
		mongo := parameters.Get("mongo").MustString()
		databases, err := GetMongoDatabases(ctx, credentials, projectID, mongo)
		if err != nil {
			return nil, err
		}
		data = make([]*suggestData, len(databases))
		for i, database := range databases {
			data[i] = &suggestData{
				Text:  database,
				Value: database,
			}
		}
	} else {
		data = make([]*suggestData, 0)
	}

	transformToTable(data, queryResult)
	result.Results = append(result.Results, queryResult)
	return result, nil
}

func transformToTable(data []*suggestData, result *datasource.QueryResult) {
	table := &datasource.Table{
		Columns: make([]*datasource.TableColumn, 2),
		Rows:    make([]*datasource.TableRow, len(data)),
	}
	table.Columns[0] = &datasource.TableColumn{
		Name: "text",
	}
	table.Columns[1] = &datasource.TableColumn{
		Name: "value",
	}

	for i, r := range data {
		values := make([]*datasource.RowValue, 2)
		values[0] = &datasource.RowValue{
			Kind:        datasource.RowValue_TYPE_STRING,
			StringValue: r.Text,
		}
		values[1] = &datasource.RowValue{
			Kind:        datasource.RowValue_TYPE_STRING,
			StringValue: r.Value,
		}
		table.Rows[i] = &datasource.TableRow{
			Values: values,
		}
	}
	result.Tables = append(result.Tables, table)
}

// Query for metric values
func MetricQuery(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	pluginLogger.Debug("MetricQuery", "Hdfd", "dfsdf")
	credentials := CredentialsFromReq(tsdbReq)
	jsonQueries, err := parseJSONQueries(tsdbReq)
	if err != nil {
		return nil, err
	}

	payload := simplejson.New()
	payload.SetPath([]string{"range", "to"}, tsdbReq.TimeRange.ToRaw)
	payload.SetPath([]string{"range", "from"}, tsdbReq.TimeRange.FromRaw)
	payload.Set("targets", jsonQueries)

	queries := tsdbReq.GetQueries()
	results := make([]*datasource.QueryResult, len(queries))
	for i, query := range queries {
		pluginLogger.Debug("MetricQuery", "queryJSON", query.ModelJson, "index", i)
		queryParameters, err := simplejson.NewJson([]byte(query.ModelJson))
		if err != nil {
			pluginLogger.Debug("MetricQuery", "json parse error", err)
			return nil, err
		}

		alias := queryParameters.Get("alias").MustString()
		metric := queryParameters.Get("metric").MustString()
		projectID := queryParameters.Get("projectId").MustString()
		mongo := queryParameters.Get("mongo").MustString()
		database := queryParameters.Get("database").MustString()
		disk := queryParameters.Get("disk").MustString()
		clusterID := queryParameters.Get("clusterId").MustString()
		dimensionID := queryParameters.Get("dimensionId").MustString()
		intervalMs := queryParameters.Get("intervalMs").MustInt()

		if metric == "database_measurements" {
			fromRaw, _ := strconv.ParseInt(tsdbReq.TimeRange.FromRaw, 10, 64)
			toRaw, _ := strconv.ParseInt(tsdbReq.TimeRange.ToRaw, 10, 64)
			rawDataPoints, err := GetDatabaseMeasurements(ctx, credentials, projectID, mongo, database, &MeasurementOptions{
				Start:       time.Unix(fromRaw/1000, 0).Format(time.RFC3339),
				End:         time.Unix(toRaw/1000, 0).Format(time.RFC3339),
				IntervalMs:  intervalMs,
				Measurement: dimensionID,
			})
			if err != nil {
				return nil, err
			}

			dataPoints := make([]*datasource.Point, len(rawDataPoints))
			for j, datapoint := range rawDataPoints {
				timestamp, _ := time.Parse("2006-01-02T15:04:05Z", datapoint.Timestamp)
				dataPoints[j] = &datasource.Point{
					Timestamp: timestamp.Unix() * 1000,
					Value:     datapoint.Value,
				}
			}
			series := make([]*datasource.TimeSeries, 1)
			name := query.RefId
			if alias != "" {
				name = alias
				projectName := queryParameters.Get("projectName").MustString()
				clusterName := queryParameters.Get("clusterName").MustString()
				dimensionName := queryParameters.Get("dimensionName").MustString()
				name = strings.ReplaceAll(name, "{{projectId}}", projectID)
				name = strings.ReplaceAll(name, "{{projectName}}", projectName)
				name = strings.ReplaceAll(name, "{{clusterId}}", clusterID)
				name = strings.ReplaceAll(name, "{{clusterName}}", clusterName)
				name = strings.ReplaceAll(name, "{{database}}", database)
				name = strings.ReplaceAll(name, "{{mongo}}", mongo)
				name = strings.ReplaceAll(name, "{{disk}}", disk)
				name = strings.ReplaceAll(name, "{{dimensionName}}", dimensionName)
			}
			series[0] = &datasource.TimeSeries{
				Name:   name,
				Points: dataPoints,
			}
			results[i] = &datasource.QueryResult{
				RefId:  query.RefId,
				Series: series,
			}
		} else if metric == "process_measurements" {
			fromRaw, _ := strconv.ParseInt(tsdbReq.TimeRange.FromRaw, 10, 64)
			toRaw, _ := strconv.ParseInt(tsdbReq.TimeRange.ToRaw, 10, 64)
			rawDataPoints, err := GetProcessMeasurements(ctx, credentials, projectID, mongo, &MeasurementOptions{
				Start:       time.Unix(fromRaw/1000, 0).Format(time.RFC3339),
				End:         time.Unix(toRaw/1000, 0).Format(time.RFC3339),
				IntervalMs:  intervalMs,
				Measurement: dimensionID,
			})
			if err != nil {
				return nil, err
			}

			dataPoints := make([]*datasource.Point, len(rawDataPoints))
			for j, datapoint := range rawDataPoints {
				timestamp, _ := time.Parse("2006-01-02T15:04:05Z", datapoint.Timestamp)
				dataPoints[j] = &datasource.Point{
					Timestamp: timestamp.Unix() * 1000,
					Value:     datapoint.Value,
				}
			}
			series := make([]*datasource.TimeSeries, 1)
			name := query.RefId
			if alias != "" {
				name = alias
				projectName := queryParameters.Get("projectName").MustString()
				clusterName := queryParameters.Get("clusterName").MustString()
				dimensionName := queryParameters.Get("dimensionName").MustString()
				name = strings.ReplaceAll(name, "{{projectId}}", projectID)
				name = strings.ReplaceAll(name, "{{projectName}}", projectName)
				name = strings.ReplaceAll(name, "{{clusterId}}", clusterID)
				name = strings.ReplaceAll(name, "{{clusterName}}", clusterName)
				name = strings.ReplaceAll(name, "{{database}}", database)
				name = strings.ReplaceAll(name, "{{mongo}}", mongo)
				name = strings.ReplaceAll(name, "{{disk}}", disk)
				name = strings.ReplaceAll(name, "{{dimensionName}}", dimensionName)
			}
			series[0] = &datasource.TimeSeries{
				Name:   name,
				Points: dataPoints,
			}
			results[i] = &datasource.QueryResult{
				RefId:  query.RefId,
				Series: series,
			}
		} else if metric == "disk_measurements" {
			fromRaw, _ := strconv.ParseInt(tsdbReq.TimeRange.FromRaw, 10, 64)
			toRaw, _ := strconv.ParseInt(tsdbReq.TimeRange.ToRaw, 10, 64)
			rawDataPoints, err := GetDiskMeasurements(ctx, credentials, projectID, mongo, disk, &MeasurementOptions{
				Start:       time.Unix(fromRaw/1000, 0).Format(time.RFC3339),
				End:         time.Unix(toRaw/1000, 0).Format(time.RFC3339),
				IntervalMs:  intervalMs,
				Measurement: dimensionID,
			})
			if err != nil {
				return nil, err
			}

			dataPoints := make([]*datasource.Point, len(rawDataPoints))
			for j, datapoint := range rawDataPoints {
				timestamp, _ := time.Parse("2006-01-02T15:04:05Z", datapoint.Timestamp)
				dataPoints[j] = &datasource.Point{
					Timestamp: timestamp.Unix() * 1000,
					Value:     datapoint.Value,
				}
			}
			series := make([]*datasource.TimeSeries, 1)
			name := query.RefId
			if alias != "" {
				name = alias
				projectName := queryParameters.Get("projectName").MustString()
				clusterName := queryParameters.Get("clusterName").MustString()
				dimensionName := queryParameters.Get("dimensionName").MustString()
				name = strings.ReplaceAll(name, "{{projectId}}", projectID)
				name = strings.ReplaceAll(name, "{{projectName}}", projectName)
				name = strings.ReplaceAll(name, "{{clusterId}}", clusterID)
				name = strings.ReplaceAll(name, "{{clusterName}}", clusterName)
				name = strings.ReplaceAll(name, "{{database}}", database)
				name = strings.ReplaceAll(name, "{{mongo}}", mongo)
				name = strings.ReplaceAll(name, "{{disk}}", disk)
				name = strings.ReplaceAll(name, "{{dimensionName}}", dimensionName)
			}
			series[0] = &datasource.TimeSeries{
				Name:   name,
				Points: dataPoints,
			}
			results[i] = &datasource.QueryResult{
				RefId:  query.RefId,
				Series: series,
			}
		} else {
			pluginLogger.Debug("MetricQuery", "metric", "unknown")
			// nothing to do
		}

	}

	response := &datasource.DatasourceResponse{
		Results: results,
	}

	return response, nil
}

func GetQueryType(tsdbReq *datasource.DatasourceRequest) (string, error) {
	queryType := "query"
	if len(tsdbReq.Queries) > 0 {
		firstQuery := tsdbReq.Queries[0]
		queryJSON, err := simplejson.NewJson([]byte(firstQuery.ModelJson))
		if err != nil {
			return "", err
		}
		queryType = queryJSON.Get("queryType").MustString("query")
	}
	return queryType, nil
}

func parseJSONQueries(tsdbReq *datasource.DatasourceRequest) ([]*simplejson.Json, error) {
	jsonQueries := make([]*simplejson.Json, 0)
	for _, query := range tsdbReq.Queries {
		json, err := simplejson.NewJson([]byte(query.ModelJson))
		if err != nil {
			return nil, err
		}

		jsonQueries = append(jsonQueries, json)
	}
	return jsonQueries, nil
}
