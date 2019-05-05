import _ from "lodash";


interface TSDBRequest {
  queries: any[];
  from?: string;
  to?: string;
}

interface TSDBQuery {
  datasourceId: string;
  target: any;
  queryType?: TSDBQueryType;
  refId?: string;
  hide?: boolean;
  type?: 'timeserie' | 'table';
}

type TSDBQueryType = 'query' | 'search' | 'test';

interface TSDBRequestOptions {
  range?: {
    from: any;
    to: any;
  };
  targets: TSDBQuery[];
}

export class GenericDatasource {
  name: string;
  type: string;
  id: string;
  url: string;
  withCredentials: boolean;
  instanceSettings: any;

  /** @ngInject */
  constructor(instanceSettings, private backendSrv, private templateSrv) {
    this.type = instanceSettings.type;
    this.url = instanceSettings.url;
    this.name = instanceSettings.name;
    this.id = instanceSettings.id;
  }

  query(options) {
    const query = this.buildQueryParameters(options);
    query.targets = query.targets.filter(t => !t.hide);

    if (query.targets.length <= 0) {
      return Promise.resolve({data: []});
    }

    return this.doTsdbRequest(query).then(handleTsdbResponse);
  }

  testDatasource() {
    return this.doTsdbRequest({targets: [{
      target: 'test',
      datasourceId: this.id,
      queryType: "test"
    }]}).then(response => {
      if (response.status === 200) {
        return { status: "success", message: "Data source is working", title: "Success" };
      } else {
        return { status: "failed", message: "Data source is not working", title: "Error" };
      }
    }).catch(error => {
      return { status: "failed", message: "Data source is not working", title: "Error" };
    });
  }

  annotationQuery(options) {
    var query = this.templateSrv.replace(options.annotation.query, {}, 'glob');
    var annotationQuery = {
      range: options.range,
      annotation: {
        name: options.annotation.name,
        datasource: options.annotation.datasource,
        enable: options.annotation.enable,
        iconColor: options.annotation.iconColor,
        query: query
      },
      rangeRaw: options.rangeRaw
    };

    return Promise.resolve()
  }

  metricFindQuery(query, params) {
    const interpolated: TSDBQuery = {
      target: query,
      datasourceId: this.id,
      queryType: "search",
      refId: query
    };

    return this.doTsdbRequest({
      targets: params ? [interpolated, params] : [interpolated]
    }).then(response => {
      return response.data.results[query].tables[0].rows.map(row => ({
        value: row[0],
        text: row[1]
      }))
    })
  }

  doTsdbRequest(options: TSDBRequestOptions) {
    const tsdbRequestData: TSDBRequest = {
      queries: options.targets,
    };

    if (options.range) {
      tsdbRequestData.from = options.range.from.valueOf().toString();
      tsdbRequestData.to = options.range.to.valueOf().toString();
    }

    return this.backendSrv.datasourceRequest({
      url: '/api/tsdb/query',
      method: 'POST',
      data: tsdbRequestData
    });
  }

  buildQueryParameters(options: any): TSDBRequestOptions {
    //remove placeholder targets
    options.targets = _.filter(options.targets, target => {
      return target.target !== 'select metric';
    });

    const targets = _.map(options.targets, target => {
      return {
        queryType: 'query',
        target: this.templateSrv.replace(target.target, options.scopedVars, 'regex'),
        refId: target.refId,
        hide: target.hide,
        type: target.type || 'timeserie',
        datasourceId: this.id,
        projectId: target.projectId,
        projectName: target.projectName,
        clusterId: target.clusterId,
        clusterName: target.clusterName,
        metric: target.metricId,
        dimensionId: target.dimensionId,
        dimensionName: target.dimensionName,
        database: target.database,
        mongo: target.mongo,
        disk: target.disk,
        intervalMs: target.intervalMs,
        maxDataPoints: target.maxDataPoints,
        alias: target.alias
      };
    });

    options.targets = targets;

    return options;
  }
}

export function handleTsdbResponse(response) {
  const res= [];
  _.forEach(response.data.results, r => {
    _.forEach(r.series, s => {
      res.push({target: s.name, datapoints: s.points});
    });
    _.forEach(r.tables, t => {
      t.type = 'table';
      t.refId = r.refId;
      res.push(t);
    });
  });

  response.data = res;
  return response;
}

export function mapToTextValue(result) {
  return _.map(result, (d, i) => {
    if (d && d.text && d.value) {
      return { text: d.text, value: d.value };
    } else if (_.isObject(d)) {
      return { text: d, value: i};
    }
    return { text: d, value: d };
  });
}
