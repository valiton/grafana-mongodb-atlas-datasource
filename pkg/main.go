package main

import (
	"os"

	"github.com/valiton/mongodbatlas-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
)

func main() {
	err := datasource.Serve(plugin.GetDatasourceOpts())

	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}

// package main

// import (
// 	"os"

// 	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
// 	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
// )

// func main() {
// 	// Start listening to requests sent from Grafana. This call is blocking so
// 	// it won't finish until Grafana shuts down the process or the plugin choose
// 	// to exit by itself using os.Exit. Manage automatically manages life cycle
// 	// of datasource instances. It accepts datasource instance factory as first
// 	// argument. This factory will be automatically called on incoming request
// 	// from Grafana to create different instances of SampleDatasource (per datasource
// 	// ID). When datasource configuration changed Dispose method will be called and
// 	// new datasource instance created using NewSampleDatasource factory.
// 	if err := datasource.Manage("myorgid-simple-backend-datasource", NewSampleDatasource, datasource.ManageOpts{}); err != nil {
// 		log.DefaultLogger.Error(err.Error())
// 		os.Exit(1)
// 	}
// }
