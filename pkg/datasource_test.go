package main

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/grafana/grafana_plugin_model/go/datasource"
)

func TestQuery(t *testing.T) {
	ds := MongoDbAtlasDatasource{
		logger: pluginLogger,
	}

	request := &datasource.DatasourceRequest{}

	ds.Query(context.Background(), request)
}
