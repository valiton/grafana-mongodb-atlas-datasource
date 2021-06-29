import { DataQuery, DataSourceJsonData, SelectableValue } from '@grafana/data';

import humanize from 'humanize-string';

export const DIMENSIONS = {
  process_measurements: [
    'ASSERT_REGULAR',
    'ASSERT_WARNING',
    'ASSERT_MSG',
    'ASSERT_USER',
    'CACHE_BYTES_READ_INTO',
    'CACHE_BYTES_WRITTEN_FROM',
    'CACHE_USAGE_DIRTY',
    'CACHE_USAGE_USED',
    'COMPUTED_MEMORY',
    'CONNECTIONS',
    'CONNECTIONS_MAX',
    'CONNECTIONS_PERCENT',
    'CURSORS_TOTAL_OPEN',
    'CURSORS_TOTAL_TIMED_OUT',
    'DB_STORAGE_TOTAL',
    'DB_DATA_SIZE_TOTAL',
    'DOCUMENT_METRICS_RETURNED',
    'DOCUMENT_METRICS_INSERTED',
    'DOCUMENT_METRICS_UPDATED',
    'DOCUMENT_METRICS_DELETED',
    'EXTRA_INFO_PAGE_FAULTS',
    'GLOBAL_LOCK_CURRENT_QUEUE_TOTAL',
    'GLOBAL_LOCK_CURRENT_QUEUE_READERS',
    'GLOBAL_LOCK_CURRENT_QUEUE_WRITERS',
    'LOGICAL_SIZE',
    'MEMORY_RESIDENT',
    'MEMORY_VIRTUAL',
    'MEMORY_MAPPED',
    'NETWORK_BYTES_IN',
    'NETWORK_BYTES_OUT',
    'NETWORK_NUM_REQUESTS',
    'OPCOUNTER_CMD',
    'OPCOUNTER_QUERY',
    'OPCOUNTER_UPDATE',
    'OPCOUNTER_DELETE',
    'OPCOUNTER_GETMORE',
    'OPCOUNTER_INSERT',
    'OPCOUNTER_REPL_CMD',
    'OPCOUNTER_REPL_UPDATE',
    'OPCOUNTER_REPL_DELETE',
    'OPCOUNTER_REPL_INSERT',
    'OPERATIONS_SCAN_AND_ORDER',
    'OP_EXECUTION_TIME_READS',
    'OP_EXECUTION_TIME_WRITES',
    'OP_EXECUTION_TIME_COMMANDS',
    'OPLOG_MASTER_TIME',
    'OPLOG_MASTER_LAG_TIME_DIFF',
    'OPLOG_SLAVE_LAG_MASTER_TIME',
    'OPLOG_RATE_GB_PER_HOUR',
    'QUERY_EXECUTOR_SCANNED',
    'QUERY_EXECUTOR_SCANNED_OBJECTS',
    'QUERY_TARGETING_SCANNED_PER_RETURNED',
    'QUERY_TARGETING_SCANNED_OBJECTS_PER_RETURNED',
    'TICKETS_AVAILABLE_READS',
    'TICKETS_AVAILABLE_WRITES',
    'PROCESS_CPU_USER',
    'PROCESS_CPU_KERNEL',
    'PROCESS_CPU_CHILDREN_USER',
    'PROCESS_CPU_CHILDREN_KERNEL',
    'PROCESS_NORMALIZED_CPU_USER',
    'PROCESS_NORMALIZED_CPU_KERNEL',
    'PROCESS_NORMALIZED_CPU_CHILDREN_USER',
    'PROCESS_NORMALIZED_CPU_CHILDREN_KERNEL',
    'SWAP_USAGE_USED',
    'SYSTEM_CPU_USER',
    'SYSTEM_CPU_KERNEL',
    'SYSTEM_CPU_NICE',
    'SYSTEM_CPU_IOWAIT',
    'SYSTEM_CPU_IRQ',
    'SYSTEM_CPU_SOFTIRQ',
    'SYSTEM_CPU_GUEST',
    'SYSTEM_CPU_STEAL',
    'SYSTEM_MEMORY_AVAILABLE',
    'SYSTEM_NORMALIZED_CPU_USER',
    'SYSTEM_NORMALIZED_CPU_KERNEL',
    'SYSTEM_NORMALIZED_CPU_NICE',
    'SYSTEM_NORMALIZED_CPU_IOWAIT',
    'SYSTEM_NORMALIZED_CPU_IRQ',
    'SYSTEM_NORMALIZED_CPU_SOFTIRQ',
    'SYSTEM_NORMALIZED_CPU_GUEST',
    'SYSTEM_NORMALIZED_CPU_STEAL',
  ],
  database_measurements: [
    'DATABASE_AVERAGE_OBJECT_SIZE',
    'DATABASE_COLLECTION_COUNT',
    'DATABASE_DATA_SIZE',
    'DATABASE_STORAGE_SIZE',
    'DATABASE_INDEX_SIZE',
    'DATABASE_INDEX_COUNT',
    'DATABASE_EXTENT_COUNT',
    'DATABASE_OBJECT_COUNT',
    'DATABASE_VIEW_COUNT',
  ],
  disk_measurements: [
    'DISK_PARTITION_IOPS_READ',
    'DISK_PARTITION_IOPS_WRITE',
    'DISK_PARTITION_IOPS_TOTAL',
    'DISK_PARTITION_UTILIZATION',
    'DISK_PARTITION_LATENCY_READ',
    'DISK_PARTITION_LATENCY_WRITE',
    'DISK_PARTITION_SPACE_FREE',
    'DISK_PARTITION_SPACE_USED',
    'DISK_PARTITION_SPACE_PERCENT_FREE',
    'DISK_PARTITION_SPACE_PERCENT_USED',
  ],
};

export const HUMANITZED_DIMENSIONS: Record<string, Array<SelectableValue<string>>> = {
  disk_measurements: DIMENSIONS['disk_measurements'].map(dim => ({
    value: dim,
    label: humanize(dim),
  })),
  database_measurements: DIMENSIONS['database_measurements'].map(dim => ({
    value: dim,
    label: humanize(dim),
  })),
  process_measurements: DIMENSIONS['process_measurements'].map(dim => ({
    value: dim,
    label: humanize(dim),
  })),
};

export const METRIC_TYPES = Object.keys(DIMENSIONS).map(metric => ({
  value: metric,
  label: humanize(metric),
}));

export const DEFAULT_PROJECT = {
  label: 'select project',
  value: '',
};

export const DEFAULT_CLUSTER = {
  label: 'select cluster',
  value: '',
};

export const DEFAULT_MONGO = {
  label: 'select mongo',
  value: '',
};

export const DEFAULT_DISK = {
  label: 'select disk',
  value: '',
};

export const DEFAULT_DATABASE = {
  label: '',
  value: '',
};

export interface Query extends DataQuery {
  project: SelectableValue<string>;
  cluster: SelectableValue<string>;
  queryType: string;
  mongo: SelectableValue<string>;
  database: SelectableValue<string>;
  disk: SelectableValue<string>;
  metric: SelectableValue<string>;
  dimension: SelectableValue<string>;
  alias: string;
}

export const defaultQuery: Partial<Query> = {
  project: DEFAULT_PROJECT,
  cluster: DEFAULT_CLUSTER,
  queryType: METRIC_TYPES[0].value,
  alias: '',
};

/**
 * These are options configured for each DataSource instance
 */
export interface DataSourceOptions extends DataSourceJsonData {
  publicKey?: string;
  apiType?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface SecureJsonData {
  privateKey?: string;
}

export interface Project {
  id: string;
  name: string;
}

export interface Cluster {
  id: string;
  name: string;
}

export interface Mongo {
  id: string;
  name: string;
}
