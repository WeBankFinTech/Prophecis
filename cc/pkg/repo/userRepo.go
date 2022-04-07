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
	"mlss-controlcenter-go/pkg/models"
)

func GetUserByUserId(userId int64) (models.User, error) {
	var user models.User
	err := datasource.GetDB().Find(&user, "id = ? AND enable_flag = ?", userId, 1).Error // find product with id 1
	return user, err
}

func DeleteUserById(db *gorm.DB, userId int64) error {
	return db.Table("t_user").Where("id = ?", userId).Update("enable_flag", 0).Error
}

func GetDeleteUserById(userId int64) (*models.User, error) {
	var user *models.User
	err := datasource.GetDB().Find(&user, "id = ? AND enable_flag = ?", userId, 0).Error // find product with id 1
	return user, err
}

func GetUserByUsername(name string, db *gorm.DB) (*models.User, error) {
	var user models.User
	err := db.Find(&user, "name = ? AND enable_flag = ?", name, 1).Error // find product with id 1
	return &user, err
}

func GetUserByName(name string, db *gorm.DB) (*models.User, error) {
	var user models.User
	err := db.Find(&user, "name = ? AND enable_flag = ?", name, 1).Error // find product with id 1
	return &user, err
}

func IsUserInDB(name string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.Table("t_user").Where("name = ? AND enable_flag = ?", name, 1).Count(&count).Error
	return count, err
}

func GetDeleteUserByName(name string) (models.User, error) {
	var user models.User
	err := datasource.GetDB().Find(&user, "name = ? AND enable_flag = ?", name, 0).Error // find product with id 1
	return user, err
}

func GetAllUsersByOffset(offSet int64, size int64) ([]*models.User, error) {
	var users []*models.User
	err := datasource.GetDB().Offset(offSet).Limit(size).Table("t_user").Find(&users, "enable_flag = ?", 1).Error
	return users, err
}

func GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	err := datasource.GetDB().Table("t_user").Find(&users, "enable_flag = ?", 1).Error
	return users, err
}

func GetUserByIds(ids []int) ([]*models.User, error) {
	var users []*models.User
	err := datasource.GetDB().Table("t_user").Find(&users, "enable_flag = ? AND id IN (?)", 1, ids).Error
	return users, err
}

func CountUsers() (int64, error) {
	var count int64
	err := datasource.GetDB().Table("t_user").Where("enable_flag = ?", 1).Count(&count).Error
	return count, err
}

func AddUser(user models.User, db *gorm.DB) error {
	err := db.Create(&user).Error
	return err
}

func UpdateUser(user models.User, db *gorm.DB) error {
	err := db.Model(&user).Omit("name").Updates(&user).Error
	return err
}

func GetUserByUID(uid int64) (models.User, error) {
	var user models.User
	err := datasource.GetDB().Find(&user, "uid = ?", uid).Error // find product with id 1
	return user, err
}

func CountUserByUID(uid int64) (int, error) {
	var count int
	err := datasource.GetDB().Table("t_user").Where("uid = ?", uid).Count(&count).Error // find product with id 1
	return count, err
}

func GetSAByName(username string) (models.Superadmin, error) {
	var sa models.Superadmin
	err := datasource.GetDB().Find(&sa, "name = ? AND enable_flag = ?", username, 1).Error
	return sa, err
}
