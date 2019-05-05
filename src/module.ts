import './query_parameter_ctrl';
import {GenericDatasource} from './datasource';
import {GenericDatasourceQueryCtrl} from './query_ctrl';

class GenericConfigCtrl {
  static templateUrl = 'partials/config.html';
}

class GenericAnnotationsQueryCtrl {
  static templateUrl = 'partials/annotations.editor.html'
}

export {
  GenericDatasource as Datasource,
  GenericDatasourceQueryCtrl as QueryCtrl,
  GenericConfigCtrl as ConfigCtrl,
  GenericAnnotationsQueryCtrl as AnnotationsQueryCtrl
};
