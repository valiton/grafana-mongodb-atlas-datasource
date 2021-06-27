package datasource

import (
	"context"
	"github.com/valiton/mongodbatlas-datasource/pkg/models"

	simplejson "github.com/bitly/go-simplejson"
)

type Mongo struct {
	ID       string `json:"id"`
	Name        string `json:"name"`
}

// Projects is a list of GitHub labels
type Mongos []Mongo

// GetProjects gets all labels from a GitHub repository
func GetMongos(ctx context.Context, client *MongoDBAtlasClient, opts models.ListMongosOptions) (Mongos, error) {
	body, err := client.query(ctx, "/groups/"+opts.Project+"/processes", nil)
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