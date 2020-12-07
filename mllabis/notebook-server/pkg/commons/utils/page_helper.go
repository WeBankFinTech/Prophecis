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
package utils

import (
	"errors"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/models"
)

const DEFAULT_SIZE = 10
const DEFAULT_PAGE = 1

func GetSubListByPageAndSize(list []*models.Notebook, page int, size int) (pages []*models.Notebook, err error) {
	logger.Logger().Debugf("origin list size: %v, page: %v, size: %v", len(list), page, size)
	if size <= 0 {
		//size = DEFAULT_SIZE
		err := errors.New("page must lt 0")
		return nil, err
	}
	if page < 1 {
		//page = 1
		err := errors.New("page must le 1")
		return nil, err
	}

	offset := GetOffSet(page, size)
	logger.Logger().Debugf("offset: %v,", offset)

	var subList []*models.Notebook

	last := offset + size
	logger.Logger().Debugf("offset + size: %v,", last)

	if offset >= len(list) {
		return subList, nil
	} else if last > len(list) {
		subList = list[offset:]
	} else {
		subList = list[offset:last]
	}
	logger.Logger().Debugf("len(subList): %v,", len(subList))

	return subList, nil
}

func GetOffSet(page int, size int) int {
	if page == 0 || size == 0 {
		return 0
	}
	return (page - 1) * size
}

func GetPages(total int, size int) int {
	pages := total / size
	if total%size > 0 {
		pages++
	}
	return pages
}
