package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
	dac "github.com/xinsnake/go-http-digest-auth-client"
	"golang.org/x/net/context"
)

type AtlasCredentials struct {
	PublicKey  string
	PrivateKey string
}

func MakeHttpRequest(ctx context.Context, path string, credentials *AtlasCredentials, query map[string]string) ([]byte, error) {
	var method = "GET"
	var baseURL = "https://cloud.mongodb.com/api/atlas/v1.0"
	var uri = baseURL + path

	pluginLogger.Debug("MakeHttpRequest", "URL", uri)

	var t = dac.NewTransport(credentials.PublicKey, credentials.PrivateKey)
	req, err := http.NewRequest(method, uri, nil)

	if query != nil {
		q := req.URL.Query()
		for key, value := range query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
		pluginLogger.Debug("MakeHttpRequest", "Full URL", req.URL.RequestURI())
	}

	if err != nil {
		return nil, err
	}

	resp, err := t.RoundTrip(req)

	if err != nil {
		pluginLogger.Debug("MakeHttpRequest", "error", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		pluginLogger.Debug("MakeHttpRequest", "io error", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("invalid status code. status: %v", resp.Status)
		pluginLogger.Debug("MakeHttpRequest", "status code error", err)
		return nil, err
	}

	pluginLogger.Debug("MakeHttpRequest", "body", string(body))

	return body, nil
}

type Project struct {
	ID   string
	Name string
}

func GetProjects(ctx context.Context, credentials *AtlasCredentials) ([]Project, error) {
	body, err := MakeHttpRequest(ctx, "/groups", credentials, nil)
	if err != nil {
		pluginLogger.Debug("GetProjects", "HTTP Error", err)
		return nil, err
	}

	jBody, err := simplejson.NewJson(body)
	if err != nil {
		pluginLogger.Debug("GetProjects", "JSON Parse Error", err)
		return nil, err
	}

	var unformattedProjects = jBody.Get("results")
	var numProjects = len(unformattedProjects.MustArray())

	pluginLogger.Debug("GetProjects", "raw projects", unformattedProjects, "num projects", numProjects)

	var projects = make([]Project, numProjects)
	for i := 0; i < numProjects; i++ {
		pluginLogger.Debug("TEST")
		var jProject = unformattedProjects.GetIndex(i)
		pluginLogger.Debug("GetProjects", "jProject", jProject, "index", i)
		project := Project{
			ID:   jProject.Get("id").MustString(),
			Name: jProject.Get("name").MustString(),
		}
		projects[i] = project
		pluginLogger.Debug("GetProjects", "project", project)
	}

	return projects, nil
}

type Cluster struct {
	ID   string
	Name string
}

func GetClusters(ctx context.Context, credentials *AtlasCredentials, groupId string) ([]Cluster, error) {
	body, err := MakeHttpRequest(ctx, "/groups/"+groupId+"/clusters", credentials, nil)
	if err != nil {
		return nil, err
	}

	jBody, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	var unformattedClusters = jBody.Get("results")
	var numClusters = len(unformattedClusters.MustArray())
	var clusters = make([]Cluster, numClusters)
	for i := 0; i < numClusters; i++ {
		var jCluster = unformattedClusters.GetIndex(i)

		clusters[i] = Cluster{
			ID:   jCluster.Get("id").MustString(),
			Name: jCluster.Get("name").MustString(),
		}
	}

	return clusters, nil
}

type Mongo struct {
	ID   string
	Name string
}

func GetMongos(ctx context.Context, credentials *AtlasCredentials, groupId string, clusterID string) ([]Mongo, error) {
	body, err := MakeHttpRequest(ctx, "/groups/"+groupId+"/processes", credentials, nil)
	if err != nil {
		return nil, err
	}

	jBody, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	var unformattedClusters = jBody.Get("results")
	var numMongos = len(unformattedClusters.MustArray())
	var mongos = make([]Mongo, numMongos)
	for i := 0; i < numMongos; i++ {
		var jMongo = unformattedClusters.GetIndex(i)

		mongos[i] = Mongo{
			ID:   jMongo.Get("id").MustString(),
			Name: jMongo.Get("hostname").MustString(),
		}
	}

	return mongos, nil
}

type DiskName string

func GetMongoDisks(ctx context.Context, credentials *AtlasCredentials, groupId string, mongo string) ([]DiskName, error) {
	body, err := MakeHttpRequest(ctx, "/groups/"+groupId+"/processes/"+mongo+"/disks", credentials, nil)
	if err != nil {
		return nil, err
	}

	jBody, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	var unformattedDisks = jBody.Get("results")
	var numDisks = len(unformattedDisks.MustArray())
	var disks = make([]DiskName, numDisks)
	for i := 0; i < numDisks; i++ {
		var jDisk = unformattedDisks.GetIndex(i)
		disks[i] = DiskName(jDisk.Get("partitionName").MustString())
	}

	return disks, nil
}

func GetMongoDatabases(ctx context.Context, credentials *AtlasCredentials, groupId string, mongo string) ([]string, error) {
	body, err := MakeHttpRequest(ctx, "/groups/"+groupId+"/processes/"+mongo+"/databases", credentials, nil)
	if err != nil {
		return nil, err
	}

	jBody, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	var unformattedDatabases = jBody.Get("results")
	var numDatabases = len(unformattedDatabases.MustArray())
	var databases = make([]string, numDatabases)
	for i := 0; i < numDatabases; i++ {
		var jDatabase = unformattedDatabases.GetIndex(i)
		databases[i] = jDatabase.Get("databaseName").MustString()
	}

	return databases, nil
}

type MeasurementOptions struct {
	Start       string
	End         string
	IntervalMs  int
	Measurement string
}

type DataPoint struct {
	Timestamp string
	Value     float64
}

func GetDatabaseMeasurements(ctx context.Context, credentials *AtlasCredentials, groupId, mongo, database string, options *MeasurementOptions) ([]*DataPoint, error) {
	url := "/groups/" + groupId + "/processes/" + mongo + "/databases/" + database + "/measurements"
	body, err := MakeHttpRequest(ctx, url, credentials, GetMeasurementOptions(options))
	if err != nil {
		return nil, err
	}
	return GetMeasurements(body, ctx)
}

func GetProcessMeasurements(ctx context.Context, credentials *AtlasCredentials, groupId, mongo string, options *MeasurementOptions) ([]*DataPoint, error) {
	url := "/groups/" + groupId + "/processes/" + mongo + "/measurements"
	body, err := MakeHttpRequest(ctx, url, credentials, GetMeasurementOptions(options))
	if err != nil {
		return nil, err
	}
	return GetMeasurements(body, ctx)
}

func GetDiskMeasurements(ctx context.Context, credentials *AtlasCredentials, groupId, mongo, disk string, options *MeasurementOptions) ([]*DataPoint, error) {
	url := "/groups/" + groupId + "/processes/" + mongo + "/disks/" + disk + "/measurements"
	body, err := MakeHttpRequest(ctx, url, credentials, GetMeasurementOptions(options))
	if err != nil {
		return nil, err
	}
	return GetMeasurements(body, ctx)
}

func GetMeasurements(body []byte, ctx context.Context) ([]*DataPoint, error) {
	jBody, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	jMeasurements := jBody.Get("measurements")
	pluginLogger.Debug("GetMeasurements", "measurements", jMeasurements)
	if len(jMeasurements.MustArray()) == 0 {
		return make([]*DataPoint, 0), nil
	}
	firstMeasurement := jMeasurements.GetIndex(0)
	pluginLogger.Debug("GetMeasurements", "first measurement", firstMeasurement)
	var rawDataPoints = firstMeasurement.Get("dataPoints")
	pluginLogger.Debug("GetMeasurements", "raw data points", rawDataPoints)
	var numDataPoints = len(rawDataPoints.MustArray())
	var dataPoints = make([]*DataPoint, 0, numDataPoints)
	for i := 0; i < numDataPoints; i++ {
		var jDataPoint = rawDataPoints.GetIndex(i)

		dataPointValue := jDataPoint.Get("value")

		// filter out all JSON null values that are sent by Atlas
		if dataPointValue.Interface() == nil {
			continue
		}

		pluginLogger.Debug("GetMeasurements", "data point", jDataPoint)

		dataPoint := &DataPoint{
			Timestamp: jDataPoint.Get("timestamp").MustString(),
			Value:     dataPointValue.MustFloat64(),
		}
		dataPoints = append(dataPoints, dataPoint)
	}

	pluginLogger.Debug("GetMeasurements", "Final data points", dataPoints)

	return dataPoints, nil
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
