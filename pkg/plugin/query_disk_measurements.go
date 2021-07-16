package plugin

import (
	"context"

	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/dfutil"
	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handleDiskMeasurementsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.DiskMeasurementsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleDiskMeasurementsQuery(ctx, query, q))
}

func (s *Server) HandleDiskMeasurements(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleDiskMeasurementsQuery),
	}, nil
}