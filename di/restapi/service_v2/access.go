package service_v2

import (
	repo "webank/DI/pkg/v2/repo"
	thirdRepo "webank/DI/pkg/v2/repo/third"
)

func CanAccessExp(userName string, clusterType string, expId string) (bool, error) {
	exp, err := repo.ExperimentRepo.GetLatestVersion(expId)
	if err != nil {
		log.Errorf("ExperimentRepo.GetLatestVersion/GetByVersion failed: +%v", err)
		return false, err
	}
	// (1) 获取用户权限相关的数据
	userDao, err := thirdRepo.UserRepo.GetUser(userName, clusterType)
	if err != nil {
		log.Errorf("thirdRepo.UserRepo.GetUser(%s) failed: %+v", userName, err)
		return false, err
	}
	isSA := thirdRepo.UserRoleRepo.IsSA(userDao.UserID)
	var groupIDs []string
	if !isSA {
		groupIDs, err = thirdRepo.UserRoleRepo.GetGroupIdsByUserId(userDao.UserID)
		if err != nil {
			log.Errorf("thirdRepo.UserRoleRepo.GetGroupIdsByUserId(%s) failed: %+v", userName, err)
			return false, err
		}
	}

	// 接口权限处理逻辑
	var canAccess bool = false
	if !isSA {
		for _, item := range groupIDs {
			if exp.GroupID == item {
				canAccess = true
				break
			}
			if exp.CreateUser == userName {
				canAccess = true
				break
			}
		}
	} else {
		canAccess = true
	}

	if !canAccess {
		return false, nil
	} else {
		return true, nil
	}
}

func CanAccessExpRun(userName string, clusterType string, expRunId string) (bool, error) {
	expRun, err := repo.ExperimentRunRepo.Get(expRunId)
	if err != nil {
		log.Errorf("ExperimentRepo.ExperimentRunRepo.Get failed: +%v", err)
		return false, err
	}
	if expRun.CreateUser == userName {
		return true, nil
	}
	return CanAccessExp(userName, clusterType, expRun.ExpID)
}
