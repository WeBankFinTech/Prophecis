
import flow from '../views/flow/process/images/newIcon/flow.svg'
import GPU from '../views/flow/process/images/GPU.svg'
import model from '../views/flow/process/images/model.svg'
import image from '../views/flow/process/images/image.svg'
import service from '../views/flow/process/images/service.svg'
import LogisticRegression from '../views/flow/process/images/LogisticRegression.svg'
import DecisionTree from '../views/flow/process/images/DecisionTree.svg'
import RandomForest from '../views/flow/process/images/RandomForest.svg'
import XGBoost from '../views/flow/process/images/XGBoost.svg'
import LightGBM from '../views/flow/process/images/LightGBM.svg'
import reportPush from '../views/flow/process/images/reportPush.svg'
const NODETYPE = {
  FLOW: 'workflow.subflow',
  GPU: 'linkis.appconn.mlflow.gpu',
  MODEL: 'linkis.appconn.mlflow.model',
  IMAGE: 'linkis.appconn.mlflow.image',
  SERVICE: 'linkis.appconn.mlflow.service',
  MODELPUSH: 'linkis.appconn.mlflow.modelPush',
  REPORTPUSH: 'linkis.appconn.mlflow.reportPush',
  LOGISTICREGRESSION: 'linkis.appconn.mlflow.LogisticRegression',
  DECISIONTREE: 'linkis.appconn.mlflow.DecisionTree',
  RANDOMFOREST: 'linkis.appconn.mlflow.RandomForest',
  XGBOOST: 'linkis.appconn.mlflow.XGBoost',
  LIGHTGBM: 'linkis.appconn.mlflow.LightGBM'
}
const NODEICON = {
  [NODETYPE.FLOW]: {
    icon: flow,
    class: { 'flow': true }
  },
  [NODETYPE.GPU]: {
    icon: GPU,
    class: { 'gpu': false }
  },
  [NODETYPE.MODEL]: {
    icon: model,
    class: { 'model': false }
  },
  [NODETYPE.IMAGE]: {
    icon: image,
    class: { 'image': false }
  },
  [NODETYPE.SERVICE]: {
    icon: service,
    class: { 'service': false }
  },
  [NODETYPE.REPORTPUSH]: {
    icon: reportPush,
    class: { 'reportPush': false }
  },
  [NODETYPE.LOGISTICREGRESSION]: {
    icon: LogisticRegression,
    class: { 'logistic': false }
  },
  [NODETYPE.DECISIONTREE]: {
    icon: DecisionTree,
    class: { 'decisionTree': false }
  },
  [NODETYPE.RANDOMFOREST]: {
    icon: RandomForest,
    class: { 'randomForest': false }
  },
  [NODETYPE.XGBOOST]: {
    icon: XGBoost,
    class: { 'XGBoost': false }
  },
  [NODETYPE.LIGHTGBM]: {
    icon: LightGBM,
    class: { 'lightGBM': false }
  }
}
export { NODETYPE, NODEICON }
