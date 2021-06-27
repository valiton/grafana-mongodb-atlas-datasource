package plugin

import (
	"context"

	"github.com/valiton/mongodbatlas-datasource/pkg/dfutil"
	"github.com/valiton/mongodbatlas-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// The Datasource type handles the requests sent to the datasource backend
type Datasource interface {
	HandleDatabaseMeasurementsQuery(context.Context, *models.DatabaseMeasurementsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleProcessMeasurementsQuery(context.Context, *models.ProcessMeasurementsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleDiskMeasurementsQuery(context.Context, *models.DiskMeasurementsQuery, backend.DataQuery) (dfutil.Framer, error)
	CheckHealth(context.Context) error
}

// HandleQueryData handles the `QueryData` request for the Github datasource
func HandleQueryData(ctx context.Context, d Datasource, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	m := GetQueryHandlers(&Server{
		Datasource: d,
	})

	return m.QueryData(ctx, req)
}

// CheckHealth ensures that the datasource settings are able to retrieve data from GitHub
func CheckHealth(ctx context.Context, d Datasource, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if err := d.CheckHealth(ctx); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	return &backend.CheckHealthResult{
		Status: backend.HealthStatusOk,
		Message: backend.HealthStatusOk.String(),
	}, nil
}