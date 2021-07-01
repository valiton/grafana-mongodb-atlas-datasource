import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import {} from '@emotion/core'; // required for fix issue https://github.com/grafana/grafana/issues/26512
import { SelectableValue, QueryEditorProps } from '@grafana/data';
import { SegmentAsync, Segment, Input, Icon, Tooltip } from '@grafana/ui';
import { DataSource } from './datasource';
import { defaultQuery, DataSourceOptions, Query, METRIC_TYPES, HUMANITZED_DIMENSIONS } from './types';

type Props = QueryEditorProps<DataSource, Query, DataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  isValid = () => {
    const query = defaults(this.props.query, defaultQuery);
    const { project, cluster, queryType, mongo, database, disk, dimension } = query;

    if (queryType === 'process_measurements') {
      return project && cluster && mongo && dimension;
    } else if (queryType === 'database_measurements') {
      return project && cluster && mongo && database && dimension;
    } else if (queryType === 'disk_measurements') {
      return project && cluster && disk && dimension;
    } else {
      return false;
    }
  };

  onValueChange = (field: string) => (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, [field]: value });
    if (this.isValid()) {
      onRunQuery();
    }
  };

  onQueryTypeChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, queryType: value.value || '' });
    if (this.isValid()) {
      onRunQuery();
    }
  };

  onAliasChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    if (!value.target.value) {
      return;
    }
    onChange({ ...query, alias: value.target.value });
    if (this.isValid()) {
      onRunQuery();
    }
  };

  getProjects = async (): Promise<Array<SelectableValue<string>>> => {
    const projects = await this.props.datasource.getProjects();

    return projects.map(({ id, name }) => ({
      value: id,
      label: name,
    }));
  };

  getClusters = async (): Promise<Array<SelectableValue<string>>> => {
    const clusters = await this.props.datasource.getClusters(this.props.query.project.value);

    return clusters.map(({ id, name }) => ({
      value: id,
      label: name,
    }));
  };

  getMongos = async (): Promise<Array<SelectableValue<string>>> => {
    const mongos = await this.props.datasource.getMongos(this.props.query.project.value);

    return mongos.map(({ id, name }) => ({
      value: id,
      label: name,
    }));
  };

  getDatabases = async (): Promise<Array<SelectableValue<string>>> => {
    const databases = await this.props.datasource.getDatabases(
      this.props.query.project.value,
      this.props.query.mongo.value
    );

    return databases.map(name => ({
      value: name,
      label: name,
    }));
  };

  getDisks = async (): Promise<Array<SelectableValue<string>>> => {
    const disks = await this.props.datasource.getDisks(this.props.query.project.value, this.props.query.mongo.value);

    return disks.map(name => ({
      value: name,
      label: name,
    }));
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { alias, project, cluster, queryType, mongo, database, disk, dimension } = query;
    const metricsTypeValue = queryType;

    const showQueryTypes = project.value && cluster.value;
    const showMongoField = showQueryTypes;
    const showDatabaseField = showQueryTypes && metricsTypeValue === 'database_measurements';
    const showDiskField = showQueryTypes && metricsTypeValue === 'disk_measurements';
    const showDimensionsField = showQueryTypes;

    return (
      <>
        <div className="gf-form-inline">
          <div className="gf-form">
            <span className="gf-form-label width-8 query-keyword">Metric</span>

            <SegmentAsync
              value={project}
              loadOptions={this.getProjects}
              onChange={this.onValueChange('project')}
              placeholder="select project"
            />

            {project.value && (
              <SegmentAsync
                value={cluster}
                loadOptions={this.getClusters}
                onChange={this.onValueChange('cluster')}
                placeholder="select cluster"
              />
            )}

            {showQueryTypes && (
              <Segment
                value={METRIC_TYPES.find(t => t.value === queryType) || METRIC_TYPES[0]}
                options={METRIC_TYPES}
                onChange={this.onQueryTypeChange}
                placeholder="select metric"
              />
            )}
          </div>
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <span className="gf-form-label width-8 query-keyword">Dimensions</span>

            {showMongoField && (
              <SegmentAsync
                value={mongo}
                loadOptions={this.getMongos}
                onChange={this.onValueChange('mongo')}
                placeholder="select mongo"
              />
            )}

            {showDatabaseField && cluster && (
              <SegmentAsync
                value={database}
                loadOptions={this.getDatabases}
                onChange={this.onValueChange('database')}
                placeholder="select database"
              />
            )}

            {showDiskField && cluster && (
              <SegmentAsync
                value={disk}
                loadOptions={this.getDisks}
                onChange={this.onValueChange('disk')}
                placeholder="select disk"
              />
            )}

            {showDimensionsField && cluster && (
              <Segment
                value={dimension || ''}
                options={queryType ? HUMANITZED_DIMENSIONS[queryType] : []}
                onChange={this.onValueChange('dimension')}
                placeholder="select dimension"
              />
            )}
          </div>
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <span className="gf-form-label width-8 query-keyword">Alias</span>
            <Input
              name="refId"
              value={alias}
              onChange={this.onAliasChange}
              addonAfter={
                <div style={{ display: 'flex', alignItems: 'center', padding: '5px' }}>
                  <Tooltip content="Alias for your query, you can also use the following variables: {{ projectName }}, {{ clusterName }}, {{ mongo }}, {{ disk }}, {{ database }}, {{ dimensionName }}">
                    <Icon name="question-circle" />
                  </Tooltip>
                </div>
              }
            />
          </div>
        </div>
      </>
    );
  }
}
