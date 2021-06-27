package datasource

import (
	"context"
	"github.com/valiton/mongodbatlas-datasource/pkg/models"

	simplejson "github.com/bitly/go-simplejson"
)

type Cluster struct {
	ID       string `json:"id"`
	Name        string `json:"name"`
}

// Projects is a list of GitHub labels
type Clusters []Cluster

// GetProjects gets all labels from a GitHub repository
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