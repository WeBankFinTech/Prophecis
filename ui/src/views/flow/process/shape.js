/*
 * Copyright 2019 WeBank
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import { NODETYPE, NODEICON } from '../../../services/nodeType'
export default [
  {
    title: 'flow.algorithmTraining',
    children: [{
      type: NODETYPE.LOGISTICREGRESSION,
      title: 'flow.logisticRegression',
      image: NODEICON[NODETYPE.LOGISTICREGRESSION].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    },
    {
      type: NODETYPE.DECISIONTREE,
      title: 'flow.decisionTree',
      image: NODEICON[NODETYPE.DECISIONTREE].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    },
    {
      type: NODETYPE.RANDOMFOREST,
      title: 'flow.randomForest',
      image: NODEICON[NODETYPE.RANDOMFOREST].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    },
    {
      type: NODETYPE.XGBOOST,
      title: 'flow.XGBoost',
      image: NODEICON[NODETYPE.XGBOOST].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    },
    {
      type: NODETYPE.LIGHTGBM,
      title: 'flow.LightGBM',
      image: NODEICON[NODETYPE.LIGHTGBM].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    }]
  },
  {
    title: 'flow.deepLearningGpu',
    children: [{
      type: NODETYPE.GPU,
      title: 'GPU',
      image: NODEICON[NODETYPE.GPU].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    }]
  },
  {
    title: 'flow.modelFactory',
    children: [{
      type: NODETYPE.MODEL,
      title: 'flow.modelRegistration',
      image: NODEICON[NODETYPE.MODEL].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    },
    {
      type: NODETYPE.IMAGE,
      title: 'flow.imageBuild',
      image: NODEICON[NODETYPE.IMAGE].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    },
    {
      type: NODETYPE.SERVICE,
      title: 'flow.modelDeployment',
      image: NODEICON[NODETYPE.SERVICE].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    },
    {
      type: NODETYPE.REPORTPUSH,
      title: 'flow.reportPush',
      image: NODEICON[NODETYPE.REPORTPUSH].icon,
      editParam: false,
      editBaseInfo: false,
      submitToScheduler: true,
      enableCopy: true,
      shouldCreationBeforeNode: false,
      supportJump: true
    }]
  }
]
