## Grafana MongoDB Atlas Logs Datasource

MongoDB Atlas allows to fetch logs from their service. More information can be found here: https://docs.atlas.mongodb.com/reference/api/logs/

This plugin allows to fetch [process](https://docs.atlas.mongodb.com/reference/api/process-measurements/), [database](https://docs.atlas.mongodb.com/reference/api/process-databases-measurements/) and [disk](https://docs.atlas.mongodb.com/reference/api/process-disks-measurements/) logs from MongoDB Atlas in your Grafana dashboard. This allows you to monitor your whole MongoDB Atlas infrastructure within your grafana dashboards.

## Installation

### Grafana Setup

You can load the latest plugin version with the following command:

```bash
grafana-cli --pluginUrl https://github.com/valiton/grafana-mongodb-atlas-datasource/releases/latest/download/grafana-mongodb-atlas-datasource.zip plugins install grafana-mongodb-atlas-datasource
```

> Please note that we currently only build for linux. If you have a windows machine, then you have to update the Makefile accordingly

For docker setup add the following environment variable to automatically install the plugin:

```bash
docker run -p 3000:3000 -e GF_INSTALL_PLUGINS="https://github.com/valiton/grafana-mongodb-atlas-datasource/releases/latest/download/grafana-mongodb-atlas-datasource.zip;mongodb-atlas-datasource" -e "GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=mongodb-atlas-datasource" grafana/grafana
```

For more information about the plugin installation have a look at the [plugin official documentation](https://grafana.com/docs/plugins/installation/).

### Dev setup

This plugin requires node > 8.10 and [dep](https://golang.github.io/dep/docs/installation.html)

```sh
npm install # install JavaScript dependencies
dep ensure  # install go dependencies
make        # build JavaScript frontend and Go backend
```

# Usage

## Create datasource

After installing the datasource in Grafana (see Grafana Setup section), you can create a Grafana datasource.

![Select MongoDB Atlas Logs datasource from list](./screenshots/datasource_list.png)

Please enter here your Atlas email address and the Atlas API token in the two input fields and click on enter. If the credentials are valid, you will see a green info box. For more information, have a look at the [MongoDB Atlas documentation](TBD) to create these credentials.

![Enter your MongoDB Atlas credentials to the form](./screenshots/datasource_setup.png)

## Create Panel

After setting up the datasource, you are able to create a query for a Grafana panel. You have to first select here the project you want to monitor and the cluster. After that, you can select one of three different metrics:

1. [Process logs](https://docs.atlas.mongodb.com/reference/api/process-measurements/),
2. [Database logs](https://docs.atlas.mongodb.com/reference/api/process-databases-measurements/) and
3. [Disk logs](https://docs.atlas.mongodb.com/reference/api/process-disks-measurements/)

Next, you are asked different other parameters, such as the database name and then you can select the dimension you want to display in the query. To name the query, please use the `alias` input. You can use `{{name}}` to use metrics or dimensions for the name (see hint field of `alias` for more information).

![Enter parameters for your MongoDB Atlas Query](./screenshots/query_setup.png)

# Limitations

- Annotations are not supported yet

# Contributing

Pull requests for new features, bug fixes, and suggestions are welcome!

# Release

**1. Add Release Notes to Changelog in README.md**

**2. Update version in src/plugin.json**

**3. Update package.json version**

**4. Create Tag with format vx.y.z**
> We use semversion format for tagging the releases.´

**5. Create Relase Zip**

```bash
make
zip --exclude "*node_modules*" --exclude "*vendor*" --exclude "*\.git*" -r grafana-mongodb-atlas-datasource.zip ./
```

**6. Create Release with zip files as attachment**

see https://help.github.com/en/articles/creating-releases for more information

# Changelog

- **1.0.0** - Initial release

  Support for process, database and disk logs

- **1.0.1** - Remove empty data points from atlas logs

  The logs by Atlas contain a lot of datapoints with null values. They were removed with this release.

- **1.0.2** - Rename Email / API Token to Public Key / Private Key

  API keys aren't bound to accounts anymore: MongoDB deprecated the Personal API Keys in favor of the Programmatic API Keys.

- **1.0.3** - Support Other Timezones
  
  https://github.com/valiton/grafana-mongodb-atlas-datasource/commit/8efac61b1d1eb7915373028e2f98986c2c42923a

- **1.0.4** - Fix alerting errors
  
  https://github.com/valiton/grafana-mongodb-atlas-datasource/pull/15

- **1.1.0** - Fix alerting errors
  https://github.com/valiton/grafana-mongodb-atlas-datasource/commit/8efac61b1d1eb7915373028e2f98986c2c42923a

- **1.2.0** Add Metric & Improve Documentation
  - Add LOGICAL_SIZE Metric
  - Add security fixes
  - Update authentication images in README

# License

[MIT](./LICENSE.txt)
