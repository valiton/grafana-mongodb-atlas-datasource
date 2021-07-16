package plugin

import (
	"context"

	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/dfutil"
	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Datasource interface {
	HandleDatabaseMeasurementsQuery(context.Context, *models.DatabaseMeasurementsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleProcessMeasurementsQuery(context.Context, *models.ProcessMeasurementsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleDiskMeasurementsQuery(context.Context, *models.DiskMeasurementsQuery, backend.DataQuery) (dfutil.Framer, error)
	CheckHealth(context.Context) error
}

func HandleQueryData(ctx context.Context, d Datasource, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	m := GetQueryHandlers(&Server{
		Datasource: d,
	})

	return m.QueryData(ctx, req)
}

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