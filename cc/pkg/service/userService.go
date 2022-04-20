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
package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"math"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
	"os"
	"path"
)

func GetUserByUserId(userId int64) (models.User, error) {
	user, err := repo.GetUserByUserId(userId)
	if err != nil {
		logger.Logger().Error("Get User by user id error:%v", err.Error())
	}
	return user, err
}

func GetUserByUserName(userName string) (*models.User, error) {
	log := logger.Logger()
	db := datasource.GetDB()
	user, err := repo.GetUserByName(userName, db)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return user, nil
}

func CheckUserPermission(userName string) (bool, error) {
	log := logger.Logger()
	db := datasource.GetDB()
	count, err := repo.IsUserInDB(userName, db)
	if err != nil {
		log.Errorf("CheckUserPermission", err.Error())
		return false, err
	}
	if count == 0 {
		return false, err
	}
	return true, nil
}

func DeleteUserById(userId int64) (*models.User, error) {
	tx := datasource.GetDB().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	userGroupByUserId, err := repo.GetAllUserGroupByUserId(userId)
	if err != nil {
		logger.Logger().Error("Get User by user id error:%v", err.Error())
		return nil, err
	}

	if userGroupByUserId != nil && len(userGroupByUserId) > 0 && userGroupByUserId[0].ID != 0 {
		if err := repo.DeleteUserGroupByUserId(tx, userId); nil != err {
			tx.Rollback()
			return nil, err
		}
	}

	user, err := repo.GetUserByUserId(userId)
	if err != nil {
		logger.Logger().Error("Get User by user id error:%v", err.Error())
		return nil, err
	}
	groupName := "gp-private-" + user.Name

	grByName, err := repo.GetGroupByName(groupName, tx)
	if err != nil {
		logger.Logger().Error("GetGroupByName:%v", err.Error())
		return nil, err
	}

	if groupName == grByName.Name && constants.TypePRIVATE == grByName.GroupType {
		err := DeleteGroupByIdDB(tx, grByName.ID)
		if nil != err {
			logger.Logger().Errorf("failed to delete group from db with msg: %v", err.Error())
			tx.Rollback()
			return nil, err
		}

	}

	if err := repo.DeleteUserById(tx, userId); nil != err {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit().Error
	if err != nil {
		logger.Logger().Errorf("failed to commit transaction: %v", err.Error())
		return nil, err
	}

	deleteUser, err := repo.GetDeleteUserById(userId)
	if err != nil {
		logger.Logger().Errorf("failed to GetDeleteUserById: %v", err.Error())
		return nil, err
	}
	return deleteUser, err
}

func GetAllUser() ([]*models.User, error) {
	return repo.GetAllUsers()
}

func GetAllUsers(page int64, size int64) (*models.PageUserList, error) {
	var allUsers []*models.User
	offSet := common.GetOffSet(page, size)
	total, err := repo.CountUsers()
	if err != nil {
		logger.Logger().Errorf("CountUsers Error:%v", err.Error())
		return nil, err
	}
	var totalPage float64
	if size > 0 {
		totalPage = math.Ceil(float64(total) / float64(size))
	}

	if page > 0 && size > 0 {
		allUsers, err = repo.GetAllUsersByOffset(offSet, size)
	} else {
		allUsers, err = repo.GetAllUsers()
	}
	if err != nil {
		logger.Logger().Errorf("GetAllUsersByOffset or GetAllUsers error:%v", err.Error())
	}

	var pageUser = models.PageUserList{
		Models:     allUsers,
		Total:      total,
		TotalPage:  int64(totalPage),
		PageNumber: page,
		PageSize:   size,
	}

	return &pageUser, err
}

func AddUser(user models.UserRequest) (*models.User, error) {
	log := logger.Logger()
	var uByName *models.User
	var err error
	err = datasource.GetDB().Transaction(func(tx *gorm.DB) error {
		// create user
		if err := createUser(common.FromAddUserRequest(user), tx); nil != err {
			log.Errorf(err.Error())
			return err
		}

		// create group for user
		if err := createGroupForUser(user, tx); nil != err {
			log.Errorf(err.Error())
			return err
		}

		uByName, err = repo.GetUserByUsername(user.Name, tx)
		if nil != err {
			log.Errorf(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		//log.Errorf("AddUser failed, %v", err.Error())
		return nil, err
	}
	return uByName, nil
}

func UpdateUser(user models.User) (*models.User, error) {

	//更新user gid uid时校验权限
	uid := user.UID
	gid := user.Gid
	pathStr := path.Join(constants.PRO_USER_PATH, user.Name)

	if isExist, _ := common.PathExists(pathStr); !isExist {
		err := os.MkdirAll(pathStr, os.ModePerm)
		if err != nil {
			return nil, err
		}
		err = os.Chown(pathStr, int(uid), int(gid))
		if err != nil {
			return nil, err
		}
	}

	if user.Type == constants.TypeSYSTEM {
		err := os.Chmod(pathStr, 0770)
		if err != nil {
			return nil, err
		}
	} else {
		err := os.Chmod(pathStr, 0700)
		if err != nil {
			return nil, err
		}
	}

	log := logger.Logger()
	db := datasource.GetDB()
	err := repo.UpdateUser(user, db)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	result, err := repo.GetUserByName(user.Name, db)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return result, nil
}

func GetUserByUID(uid int64) (*models.User, error) {
	user, err := repo.GetUserByUID(uid)
	if err != nil {
		logger.Logger().Errorf("getUserByUID error:%v", err.Error())
		return nil, err
	}
	return &user, err
}

func IfUserIdExisted(uid int64) (bool, error) {
	count, err := repo.CountUserByUID(uid)
	if err != nil {
		logger.Logger().Errorf("getUserByUID error:%v", err.Error())
		return false, err
	}
	if count == 0 {
		return false, err
	}
	return true, err
}

func GetSAByName(username string) models.Superadmin {
	sa, err := repo.GetSAByName(username)
	if err != nil {
		logger.Logger().Errorf("GetSAByName error:%v", err.Error())
	}
	return sa
}

func createUser(user models.User, tx *gorm.DB) error {
	log := logger.Logger()
	var err error
	_, err = repo.GetUserByName(user.Name, tx)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Errorf(err.Error())
		return err
	}
	if err == nil {
		return errors.New("user is exist in db")
	}

	delUByName, err := repo.GetDeleteUserByName(user.Name)
	if err != nil {
		err = repo.AddUser(user, tx)
		if err != nil {
			log.Errorf(err.Error())
			return err
		}
	} else {
		if user.Name == delUByName.Name {
			user.ID = delUByName.ID
			err = repo.UpdateUser(user, tx)
			if err != nil {
				log.Errorf(err.Error())
				return err
			}
		}
	}
	return nil
}

func createGroupForUser(user models.UserRequest, tx *gorm.DB) error {
	log := logger.Logger()
	// create group for user
	groupName := constants.GROUP_SUB_NAME + user.Name

	dbGroup, _ := repo.GetDeleteGroupByName(groupName)

	if groupName == dbGroup.Name {
		dbGroup.EnableFlag = 1
		err := repo.UpdateGroup(*dbGroup, tx)
		if err != nil {
			log.Errorf(err.Error())
			return err
		}
	} else {
		var group = models.Group{
			Remarks:        constants.USER_GROUP,
			GroupType:      constants.TypePRIVATE,
			Name:           groupName,
			DepartmentName: constants.MLSS,
			DepartmentID:   0,
			EnableFlag:     1,
		}
		err := repo.AddGroup(group, tx)
		if err != nil {
			log.Errorf(err.Error())
			return err
		}
	}

	//add user to private group
	dbUser, err := repo.GetUserByName(user.Name, tx)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	log.Debug("%+v", dbUser)
	dbGroup, err = repo.GetGroupByName(groupName, tx)
	if err != nil {
		logger.Logger().Errorf("GetGroupByName error:%v", err.Error())
	}
	log.Debug("%+v", dbGroup)

	dbUG, err := repo.GetDeleteUserGroupByUserIdAndGroupId(dbUser.ID, dbGroup.ID)
	if err != nil {
		logger.Logger().Errorf("GetGroupByName error:%v", err.Error())
	}
	if dbUser.ID == dbUG.UserID && dbGroup.ID == dbUG.GroupID && 0 != dbGroup.ID && dbUser.ID != 0 {
		dbUG.EnableFlag = 1
		if err := repo.UpdateUserGroupDB(tx, *dbUG); nil != err {
			//tx.Rollback()
			log.Errorf(err.Error())
			return err
		}
	} else {
		var userGroup = models.UserGroup{
			Remarks:    constants.USER_GROUP,
			UserID:     dbUser.ID,
			GroupID:    dbGroup.ID,
			RoleID:     constants.GA_ID,
			EnableFlag: 1,
		}
		if err := repo.AddUserToGroupDB(tx, userGroup); nil != err {
			//tx.Rollback()
			log.Errorf(err.Error())
			return err
		}
	}

	return nil
}

func HasSA(username string) bool {
	_, err := repo.GetSAByName(username)
	if err != nil {
		return false
	}
	return true
}

func GetDeleteUserByUserName(userName string) (models.User, error) {
	user, err := repo.GetDeleteUserByName(userName)
	if err != nil {
		logger.Logger().Errorf("GetSAByName error:%v", err.Error())
	}
	return user, err
}
