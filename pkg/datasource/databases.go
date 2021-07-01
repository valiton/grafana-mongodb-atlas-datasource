package datasource

import (
	"context"
	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/models"

	simplejson "github.com/bitly/go-simplejson"
)

type Databases []string

func GetDatabases(ctx context.Context, client *MongoDBAtlasClient, opts models.ListDatabasesOptions) (Databases, error) {
	body, err := client.query(ctx, "/groups/"+opts.Project+"/processes/"+opts.Mongo+"/databases", nil)
	if err != nil {
		return nil, err
	}

	jBody, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	var unformattedDatabases = jBody.Get("results")
	var numDatabases = len(unformattedDatabases.MustArray())
	var databases = make(Databases, numDatabases)
	for i := 0; i < numDatabases; i++ {
		var jDatabase = unformattedDatabases.GetIndex(i)
		databases[i] = jDatabase.Get("databaseName").MustString()
	}

	return databases, nil
}