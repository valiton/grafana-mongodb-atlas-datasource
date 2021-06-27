package plugin

import (
	"context"

	"github.com/valiton/mongodbatlas-datasource/pkg/dfutil"
	"github.com/valiton/mongodbatlas-datasource/pkg/datasource"
	"github.com/valiton/mongodbatlas-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// Instance is the root Datasource implementation that wraps a Datasource
type Instance struct {
	Datasource Datasource
	Handlers   Handlers
}

// HandleDatabaseMeasurementsQuery ...
func (i *Instance) HandleDatabaseMeasurementsQuery(ctx context.Context, q *models.DatabaseMeasurementsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleDatabaseMeasurementsQuery(ctx, q, req)
}

// HandleProcessMeasurementsQuery ...
func (i *Instance) HandleProcessMeasurementsQuery(ctx context.Context, q *models.ProcessMeasurementsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleProcessMeasurementsQuery(ctx, q, req)
}

// HandleDiskMeasurementsQuery ...
func (i *Instance) HandleDiskMeasurementsQuery(ctx context.Context, q *models.DiskMeasurementsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleDiskMeasurementsQuery(ctx, q, req)
}

// CheckHealth ...
func (i *Instance) CheckHealth(ctx context.Context) error {
	return i.Datasource.CheckHealth(ctx)
}

// NewMongoDbAtlasInstance creates a new MongoDbAtlas using the settings to determine if things like the Caching Wrapper should be enabled
func NewMongoDbAtlasInstance(ctx context.Context, settings *models.Settings) *Instance {
	var (
		ds = datasource.NewDatasource(ctx, settings)
	)

	var d Datasource = ds

	return &Instance{
		Datasource: d,
		Handlers: Handlers{
			Projects:     ds.HandleGetProjects,
			Clusters: ds.HandleGetClusters,
			Databases: ds.HandleGetDatabases,
			Mongos: ds.HandleGetMongos,
			Disks: ds.HandleGetDisks,
		},
	}
}

func newDataSourceInstance(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	datasourceSettings, err := models.LoadSettings(settings)
	if err != nil {
		return nil, err
	}

	return NewMongoDbAtlasInstance(context.Background(), datasourceSettings), nil
}
