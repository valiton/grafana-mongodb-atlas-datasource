package datasource

import (
	"context"
	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/models"

	simplejson "github.com/bitly/go-simplejson"
)

type Cluster struct {
	ID       string `json:"id"`
	Name        string `json:"name"`
}

type Clusters []Cluster

func GetClusters(ctx context.Context, client *MongoDBAtlasClient, opts models.ListClustersOptions) (Clusters, error) {
	body, err := client.query(ctx, "/groups/"+opts.Project+"/clusters", nil)
	if err != nil {
		return nil, err
	}

	jBody, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	var unformattedClusters = jBody.Get("results")
	var numClusters = len(unformattedClusters.MustArray())
	var clusters = make(Clusters, numClusters)
	for i := 0; i < numClusters; i++ {
		var jCluster = unformattedClusters.GetIndex(i)

		clusters[i] = Cluster{
			ID:   jCluster.Get("id").MustString(),
			Name: jCluster.Get("name").MustString(),
		}
	}

	return clusters, nil
}