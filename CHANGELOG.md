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

- **2.0.0** Add Metric & Improve Documentation
  - Renamed plugin from `grafana-mongodb-atlas-datasource` to `mongodb-atlas-datasource`: Required to sign the plugin
  - Add LOGICAL_SIZE Metric [#32](https://github.com/valiton/grafana-mongodb-atlas-datasource/issues/32) 
  - Add security fixes
  - Update authentication images in README

- **3.0.0** Upgraded to Grafana Plugin API v2
  - Renamed plugin from `mongodb-atlas-datasource` to `valiton-mongodb-atlas-datasource`: Required for new API
  - Rewrote complete plugin
