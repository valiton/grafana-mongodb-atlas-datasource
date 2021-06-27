package datasource

import (
	"context"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/valiton/mongodbatlas-datasource/pkg/models"

	simplejson "github.com/bitly/go-simplejson"
)

type Project struct {
	ID       string `json:"id"`
	Name        string `json:"name"`
}

// Projects is a list of GitHub labels
type Projects []Project

// GetProjects gets all labels from a GitHub repository
func GetProjects(ctx context.Context, client *MongoDBAtlasClient, opts models.ListProjectsOptions) (Projects, error) {
	body, err := client.query(ctx, "/groups", nil)
	if err != nil {
		log.DefaultLogger.Debug("GetProjects", "HTTP Error", err)
		return nil, err
	}

	jBody, err := simplejson.NewJson(body)
	if err != nil {
		log.DefaultLogger.Debug("GetProjects", "JSON Parse Error", err)
		return nil, err
	}

	var unformattedProjects = jBody.Get("results")
	var numProjects = len(unformattedProjects.MustArray())

	log.DefaultLogger.Debug("GetProjects", "raw projects", unformattedProjects, "num projects", numProjects)

	var projects = make(Projects, numProjects)
	for i := 0; i < numProjects; i++ {
		var jProject = unformattedProjects.GetIndex(i)
		log.DefaultLogger.Debug("GetProjects", "jProject", jProject, "index", i)
		project := Project{
			ID:   jProject.Get("id").MustString(),
			Name: jProject.Get("name").MustString(),
		}
		projects[i] = project
		log.DefaultLogger.Debug("GetProjects", "project", project)
	}

	return projects, nil
}