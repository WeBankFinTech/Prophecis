/*
 * Copyright 2020 WeBank
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
 */

package repo

import (
	"github.com/jinzhu/gorm"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
)

var repo HDFSPrivsRepo

type Token struct {
	User        string `gorm:"column:USER"`
	Token       string `gorm:"column:TOKEN"`
	ClientToken string `gorm:"column:CLIENT_TOKEN"`
}

func GetHDFSPrivsRepo() *HDFSPrivsRepo {
	return &repo
}

func init() {
	repo = HDFSPrivsRepo{}
}

type HDFSPrivsRepo struct {
}

type HDFSPrivsRepoIF interface {
	GetTokenByUser(user string, ctx ...*gorm.DB) (*Token, error)
}

func (r *HDFSPrivsRepo) GetTokenByUser(user string, tx *gorm.DB) (*Token, error) {
	log := logger.Logger()
	//db := datasource.TokenDS.GetDBByDBS(txs)
	var token Token
	err := tx.Table(datasource.WbUserToken).Find(&token, "USER = ?", user).Error
	if err != nil {
		log.Errorf("db.Find failed, err: %v", err.Error())
		return nil, err
	}
	return &token, nil
}
