import { QueryCtrl } from "grafana/app/plugins/sdk";
import "./css/query-editor.css";
import { coreModule } from "grafana/app/core/core";
import humanize from "humanize-string";

export const DIMENSIONS = {
  process_measurements: [
    "ASSERT_REGULAR",
    "ASSERT_WARNING",
    "ASSERT_MSG",
    "ASSERT_USER",
    "CACHE_BYTES_READ_INTO",
    "CACHE_BYTES_WRITTEN_FROM",
    "CACHE_USAGE_DIRTY",
    "CACHE_USAGE_USED",
    "CONNECTIONS",
    "CURSORS_TOTAL_OPEN",
    "CURSORS_TOTAL_TIMED_OUT",
    "DB_STORAGE_TOTAL",
    "DB_DATA_SIZE_TOTAL",
    "DOCUMENT_METRICS_RETURNED",
    "DOCUMENT_METRICS_INSERTED",
    "DOCUMENT_METRICS_UPDATED",
    "DOCUMENT_METRICS_DELETED",
    "EXTRA_INFO_PAGE_FAULTS",
    "GLOBAL_LOCK_CURRENT_QUEUE_TOTAL",
    "GLOBAL_LOCK_CURRENT_QUEUE_READERS",
    "GLOBAL_LOCK_CURRENT_QUEUE_WRITERS",
    "LOGICAL_SIZE",
    "MEMORY_RESIDENT",
    "MEMORY_VIRTUAL",
    "MEMORY_MAPPED",
    "NETWORK_BYTES_IN",
    "NETWORK_BYTES_OUT",
    "NETWORK_NUM_REQUESTS",
    "OPCOUNTER_CMD",
    "OPCOUNTER_QUERY",
    "OPCOUNTER_UPDATE",
    "OPCOUNTER_DELETE",
    "OPCOUNTER_GETMORE",
    "OPCOUNTER_INSERT",
    "OPCOUNTER_REPL_CMD",
    "OPCOUNTER_REPL_UPDATE",
    "OPCOUNTER_REPL_DELETE",
    "OPCOUNTER_REPL_INSERT",
    "OPERATIONS_SCAN_AND_ORDER",
    "OP_EXECUTION_TIME_READS",
    "OP_EXECUTION_TIME_WRITES",
    "OP_EXECUTION_TIME_COMMANDS",
    "OPLOG_MASTER_TIME",
    "OPLOG_MASTER_LAG_TIME_DIFF",
    "OPLOG_SLAVE_LAG_MASTER_TIME",
    "OPLOG_RATE_GB_PER_HOUR",
    "QUERY_EXECUTOR_SCANNED",
    "QUERY_EXECUTOR_SCANNED_OBJECTS",
    "QUERY_TARGETING_SCANNED_PER_RETURNED",
    "QUERY_TARGETING_SCANNED_OBJECTS_PER_RETURNED",
    "TICKETS_AVAILABLE_READS",
    "TICKETS_AVAILABLE_WRITES",
    "PROCESS_CPU_USER",
    "PROCESS_CPU_KERNEL",
    "PROCESS_CPU_CHILDREN_USER",
    "PROCESS_CPU_CHILDREN_KERNEL",
    "PROCESS_NORMALIZED_CPU_USER",
    "PROCESS_NORMALIZED_CPU_KERNEL",
    "PROCESS_NORMALIZED_CPU_CHILDREN_USER",
    "PROCESS_NORMALIZED_CPU_CHILDREN_KERNEL",
    "SYSTEM_CPU_USER",
    "SYSTEM_CPU_KERNEL",
    "SYSTEM_CPU_NICE",
    "SYSTEM_CPU_IOWAIT",
    "SYSTEM_CPU_IRQ",
    "SYSTEM_CPU_SOFTIRQ",
    "SYSTEM_CPU_GUEST",
    "SYSTEM_CPU_STEAL",
    "SYSTEM_NORMALIZED_CPU_USER",
    "SYSTEM_NORMALIZED_CPU_KERNEL",
    "SYSTEM_NORMALIZED_CPU_NICE",
    "SYSTEM_NORMALIZED_CPU_IOWAIT",
    "SYSTEM_NORMALIZED_CPU_IRQ",
    "SYSTEM_NORMALIZED_CPU_SOFTIRQ",
    "SYSTEM_NORMALIZED_CPU_GUEST",
    "SYSTEM_NORMALIZED_CPU_STEAL"
  ],
  database_measurements: [
    "DATABASE_AVERAGE_OBJECT_SIZE",
    "DATABASE_COLLECTION_COUNT",
    "DATABASE_DATA_SIZE",
    "DATABASE_STORAGE_SIZE",
    "DATABASE_INDEX_SIZE",
    "DATABASE_INDEX_COUNT",
    "DATABASE_EXTENT_COUNT",
    "DATABASE_OBJECT_COUNT",
    "DATABASE_VIEW_COUNT"
  ],
  disk_measurements: [
    "DISK_PARTITION_IOPS_READ",
    "DISK_PARTITION_IOPS_WRITE",
    "DISK_PARTITION_IOPS_TOTAL",
    "DISK_PARTITION_UTILIZATION",
    "DISK_PARTITION_LATENCY_READ",
    "DISK_PARTITION_LATENCY_WRITE",
    "DISK_PARTITION_SPACE_FREE",
    "DISK_PARTITION_SPACE_USED",
    "DISK_PARTITION_SPACE_PERCENT_FREE",
    "DISK_PARTITION_SPACE_PERCENT_USED"
  ]
};

export const HUMANITZED_DIMENSIONS = {
  disk_measurements: DIMENSIONS["disk_measurements"].reduce(
    (map, dim) => ({
      [humanize(dim)]: dim,
      ...map
    }),
    {}
  ),
  database_measurements: DIMENSIONS["database_measurements"].reduce(
    (map, dim) => ({
      [humanize(dim)]: dim,
      ...map
    }),
    {}
  ),
  process_measurements: DIMENSIONS["process_measurements"].reduce(
    (map, dim) => ({
      [humanize(dim)]: dim,
      ...map
    }),
    {}
  )
};

export class GrafanaMongoDbQueryParameterCtrl {
  projectNameIdMapping: Map<string, string>;
  clusterNameIdMapping: Map<string, string>;
  mongoNameIdMapping: Map<string, string>;
  metricMapping = {
    "Process Measurements": "process_measurements",
    "Database Measurements": "database_measurements",
    "Disk Measurements": "disk_measurements"
  };

  /** @ngInject */
  constructor($scope, templateSrv, uiSegmentSrv, datasourceSrv, $q) {
    this.projectNameIdMapping = undefined;
    this.mongoNameIdMapping = undefined;
    this.clusterNameIdMapping = undefined;
    $scope.init = () => {
      const target = $scope.target;
      target.projectId = target.projectId || "";
      target.projectName = target.projectName || "";
      target.clusterId = target.clusterId || "";
      target.clusterName = target.clusterName || "";
      target.mongoId = target.mongoId || "";
      target.disk = target.disk || "";
      target.database = target.database || "";
      target.metricId = target.metricId || "";
      target.metricName = target.metricName || "";

      target.dimensionId = target.dimensionId || "";
      target.dimensionName = target.dimensionName || "";

      target.alias = target.alias || "";

      $scope.projectSegment = uiSegmentSrv.getSegmentForValue(
        target.projectName,
        "select project"
      );
      $scope.clusterSegment = uiSegmentSrv.getSegmentForValue(
        target.clusterName,
        "select cluster"
      );
      $scope.mongoSegment = uiSegmentSrv.getSegmentForValue(
        target.mongoName,
        "select mongo"
      );
      $scope.databaseSegment = uiSegmentSrv.getSegmentForValue(
        target.database,
        "select database"
      );
      $scope.diskSegment = uiSegmentSrv.getSegmentForValue(
        target.disk,
        "select disk"
      );
      $scope.metricsSegment = uiSegmentSrv.getSegmentForValue(
        target.metricName,
        "select metric"
      );

      $scope.dimensionSegment = uiSegmentSrv.getSegmentForValue(
        target.dimensionName,
        "select dimension"
      );
    };

    $scope.getMetrics = () =>
      Promise.resolve(
        Object.entries(this.metricMapping).map(([value, key]) => ({
          key,
          value
        }))
      );

    $scope.getProjects = () => {
      return $scope.datasource.metricFindQuery("projects").then(projects => {
        this.projectNameIdMapping = projects.reduce(
          (map, project) => ({
            ...map,
            [project.value]: project.text
          }),
          {}
        );
        return projects;
      });
    };
    $scope.getClusters = () => {
      return $scope.datasource
        .metricFindQuery("clusters", {
          projectId: $scope.target.projectId
        })
        .then(clusters => {
          this.clusterNameIdMapping = clusters.reduce(
            (map, cluster) => ({
              ...map,
              [cluster.value]: cluster.text
            }),
            {}
          );
          return clusters;
        });
    };
    $scope.getMongos = () => {
      return $scope.datasource
        .metricFindQuery("mongos", {
          projectId: $scope.target.projectId,
          clusterId: $scope.target.clusterId
        })
        .then(mongos => {
          this.mongoNameIdMapping = mongos.reduce(
            (map, mongo) => ({
              ...map,
              [mongo.value]: mongo.text
            }),
            {}
          );
          return mongos;
        });
    };
    $scope.getDatabases = () => {
      return $scope.datasource.metricFindQuery("databases", {
        projectId: $scope.target.projectId,
        mongoId: $scope.target.mongoId
      });
    };
    $scope.getDisks = () => {
      return $scope.datasource.metricFindQuery("disks", {
        projectId: $scope.target.projectId,
        clusterId: $scope.target.clusterId,
        mongoId: $scope.target.mongoId
      });
    };
    $scope.getDimensions = () => {
      return Promise.resolve(
        Object.entries(HUMANITZED_DIMENSIONS[$scope.target.metricId]).map(
          ([value, key]) => ({
            key,
            value
          })
        )
      );
    };

    $scope.projectChanged = () => {
      $scope.target.projectName = $scope.projectSegment.value;
      $scope.target.projectId = this.projectNameIdMapping[
        $scope.projectSegment.value
      ];
      $scope.onChange();
    };
    $scope.clusterChanged = () => {
      $scope.target.clusterName = $scope.clusterSegment.value;
      $scope.target.clusterId = this.clusterNameIdMapping[
        $scope.clusterSegment.value
      ];
      $scope.onChange();
    };
    $scope.metricChanged = () => {
      $scope.target.metricName = $scope.metricsSegment.value;
      $scope.target.metricId = this.metricMapping[$scope.target.metricName];
      $scope.onChange();
    };
    $scope.mongoChanged = () => {
      $scope.target.mongoName = $scope.mongoSegment.value;
      $scope.target.mongoId = this.mongoNameIdMapping[
        $scope.mongoSegment.value
      ];
      $scope.onChange();
    };
    $scope.diskChanged = () => {
      $scope.target.disk = $scope.diskSegment.value;
      $scope.onChange();
    };
    $scope.databaseChanged = () => {
      $scope.target.database = $scope.databaseSegment.value;
      $scope.onChange();
    };
    $scope.dimensionChanged = () => {
      $scope.target.dimensionName = $scope.dimensionSegment.value;
      console.log(
        HUMANITZED_DIMENSIONS,
        HUMANITZED_DIMENSIONS[$scope.target.metricId],
        $scope.target.metricId
      );
      $scope.target.dimensionId =
        HUMANITZED_DIMENSIONS[$scope.target.metricId][
        $scope.target.dimensionName
        ];
      $scope.onChange();
    };
    $scope.isMetricSelected = metric => {
      return this.metricMapping[$scope.target.metricName] === metric;
    };
    $scope.transformToSegments = addTemplateVars => {
      return results => {
        const segments = results.map(segment =>
          uiSegmentSrv.newSegment({
            value: segment.text,
            expandable: segment.expandable
          })
        );

        if (addTemplateVars) {
          templateSrv.variables.forEach(variable => {
            segments.unshift(
              uiSegmentSrv.newSegment({
                type: "template",
                value: "$" + variable.name,
                expandable: true
              })
            );
          });
        }

        return segments;
      };
    };

    if (!$scope.onChange) {
      $scope.onChange = () => { };
    }
    $scope.init();
  }
}

export function grafanaMongodbAtlasQueryParameter() {
  return {
    templateUrl:
      "public/plugins/grafana-mongodb-atlas-datasource/partials/query.parameter.html",
    controller: GrafanaMongoDbQueryParameterCtrl,
    restrict: "E",
    scope: {
      target: "=",
      datasource: "=",
      onChange: "&"
    }
  };
}

coreModule.directive(
  "grafanaMongodbAtlasQueryParameter",
  grafanaMongodbAtlasQueryParameter
);
