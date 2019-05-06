## Grafana MongoDB Atlas Logs Datasource

MongoDB Atlas allows to fetch logs from their service. More information can be found here: https://docs.atlas.mongodb.com/reference/api/logs/

This plugin allows to fetch [process](https://docs.atlas.mongodb.com/reference/api/process-measurements/), [database](https://docs.atlas.mongodb.com/reference/api/process-databases-measurements/) and [disk](https://docs.atlas.mongodb.com/reference/api/process-disks-measurements/) logs from MongoDB Atlas in your Grafana dashboard. This allows you to monitor your whole MongoDB Atlas infrastructure within your grafana dashboards. 

## Installation

### Grafana Setup

TBD

### Dev setup

This plugin requires node > 8.10 and [dep](https://golang.github.io/dep/docs/installation.html)

```sh
npm install # install JavaScript dependencies
dep ensure  # install go dependencies
make        # build JavaScript frontend and Go backend
```

# Usage

## Create datasource

TBD

## Create Panel

TBD

# Contributing
Pull requests for new features, bug fixes, and suggestions are welcome! 

# Changelog

- **1.0.0** - Initial release

  Support for process, database and disk logs

# License
[MIT](./LICENSE.txt)