import { DataSourceWithBackend } from '@grafana/runtime';
import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceOptions, Query, Project, Cluster, Mongo } from './types';

export class DataSource extends DataSourceWithBackend<Query, DataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<DataSourceOptions>) {
    super(instanceSettings);
  }

  async getProjects(): Promise<Project[]> {
    return this.getResource('projects', {});
  }

  async getClusters(projectId: string | undefined): Promise<Cluster[]> {
    if (projectId === undefined) {
      return [];
    }
    return this.getResource('clusters', {
      project: projectId,
    });
  }

  async getMongos(projectId: string | undefined): Promise<Mongo[]> {
    if (projectId === undefined) {
      return [];
    }
    return this.getResource('mongos', {
      project: projectId,
    });
  }

  async getDisks(projectId: string | undefined, mongoId: string | undefined): Promise<string[]> {
    if (projectId === undefined || mongoId === undefined) {
      return [];
    }
    return this.getResource('disks', {
      project: projectId,
      mongo: mongoId,
    });
  }

  async getDatabases(projectId: string | undefined, mongoId: string | undefined): Promise<string[]> {
    if (projectId === undefined || mongoId === undefined) {
      return [];
    }
    return this.getResource('databases', {
      project: projectId,
      mongo: mongoId,
    });
  }
}
