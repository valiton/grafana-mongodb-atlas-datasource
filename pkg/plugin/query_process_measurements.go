package plugin

import (
	"context"

	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/dfutil"
	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handleProcessMeasurementsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.ProcessMeasurementsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleProcessMeasurementsQuery(ctx, query, q))
}

func (s *Server) HandleProcessMeasurements(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleProcessMeasurementsQuery),
	}, nil
}