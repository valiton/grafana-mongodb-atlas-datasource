package plugin

import (
	"context"

	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/dfutil"
	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handleDatabaseMeasurementsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.DatabaseMeasurementsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleDatabaseMeasurementsQuery(ctx, query, q))
}

func (s *Server) HandleDatabaseMeasurements(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleDatabaseMeasurementsQuery),
	}, nil
}