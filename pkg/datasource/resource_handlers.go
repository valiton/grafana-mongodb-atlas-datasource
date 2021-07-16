package datasource

import (
	"context"
	"net/http"

	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/httputil"
	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/models"
)

func handleGetProjects(ctx context.Context, client *MongoDBAtlasClient, r *http.Request) (Projects, error) {
	opts := models.ListProjectsOptions{}

	labels, err := GetProjects(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	return labels, nil
}

func (d *Datasource) HandleGetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := handleGetProjects(r.Context(), d.client, r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httputil.WriteResponse(w, projects)
}

func handleGetClusters(ctx context.Context, client *MongoDBAtlasClient, r *http.Request) (Clusters, error) {
	q := r.URL.Query()
	opts := models.ListClustersOptions{
		Project: q.Get("project"),
	}

	labels, err := GetClusters(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	return labels, nil
}

func (d *Datasource) HandleGetClusters(w http.ResponseWriter, r *http.Request) {
	clusters, err := handleGetClusters(r.Context(), d.client, r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httputil.WriteResponse(w, clusters)
}

func handleGetMongos(ctx context.Context, client *MongoDBAtlasClient, r *http.Request) (Mongos, error) {
	q := r.URL.Query()
	opts := models.ListMongosOptions{
		Project: q.Get("project"),
	}

	mongos, err := GetMongos(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	return mongos, nil
}

func (d *Datasource) HandleGetMongos(w http.ResponseWriter, r *http.Request) {
	mongos, err := handleGetMongos(r.Context(), d.client, r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httputil.WriteResponse(w, mongos)
}

func handleGetDisks(ctx context.Context, client *MongoDBAtlasClient, r *http.Request) (Disks, error) {
	q := r.URL.Query()
	opts := models.ListDisksOptions{
		Project: q.Get("project"),
		Mongo: q.Get("mongo"),
	}

	disks, err := GetDisks(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	return disks, nil
}

func (d *Datasource) HandleGetDisks(w http.ResponseWriter, r *http.Request) {
	disks, err := handleGetDisks(r.Context(), d.client, r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httputil.WriteResponse(w, disks)
}

func handleGetDatabases(ctx context.Context, client *MongoDBAtlasClient, r *http.Request) (Databases, error) {
	q := r.URL.Query()
	opts := models.ListDatabasesOptions{
		Project: q.Get("project"),
		Mongo: q.Get("mongo"),
	}

	databases, err := GetDatabases(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	return databases, nil
}

func (d *Datasource) HandleGetDatabases(w http.ResponseWriter, r *http.Request) {
	databases, err := handleGetDatabases(r.Context(), d.client, r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httputil.WriteResponse(w, databases)
}