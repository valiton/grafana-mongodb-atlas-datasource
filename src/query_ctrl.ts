import {QueryCtrl} from 'grafana/app/plugins/sdk';
import './css/query-editor.css';

export class GenericDatasourceQueryCtrl extends QueryCtrl {
  static templateUrl = 'partials/query.editor.html';

  scope: any;

  /** @ngInject */
  constructor($scope, $injector) {
    super($scope, $injector);
  }
}
