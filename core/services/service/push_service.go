///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"digger/models"
	"github.com/jinzhu/gorm"
)

type pushServiceImp struct {
}

func (p pushServiceImp) Save(source models.PushSource) error {
	err := DoTransaction(func(tx *gorm.DB) error {
		return tx.Save(&source).Error
	})
	return err
}

func (pushServiceImp) List(params *models.PushQueryVO) (int64, []*models.PushSource, error) {

	var total int64 = 0
	// 查询总数
	var countQuery = dbConn.Table("t_push_source")
	if err := countQuery.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	if total == 0 {
		return 0, nil, nil
	}

	// 查询数据
	var dataQuery = dbConn.Table("t_push_source")
	dataQuery = dataQuery.Order("id desc")
	rows, err := dataQuery.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Rows()
	if transformNotFoundErr(err) != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var retArray []*models.PushSource
	for rows.Next() {
		var item models.PushSource
		if err := dataQuery.ScanRows(rows, &item); err != nil {
			return 0, nil, err
		}
		if err != nil {
			return 0, nil, err
		}
		retArray = append(retArray, &item)
	}
	return total, retArray, nil
}

func (pushServiceImp) Delete(idList []int) error {
	return DoTransaction(func(tx *gorm.DB) error {
		return tx.Delete(models.PushSource{}, "id in (?)", idList).Error
	})
}

func (pushServiceImp) SelectByProject(projectId int) ([]*models.PushSource, error) {
	rows, err := dbConn.Raw(`
select b.* from t_project_push a 
left join t_push_source b on a.push_id = b.id
where a.project_id = ?
`, projectId).Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.PushSource
	for rows.Next() {
		var item models.PushSource
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		retArray = append(retArray, &item)
	}
	return retArray, nil
}

func (pushServiceImp) UpdatePushTaskResultId(taskId int, lastResultId int64) error {
	err := DoTransaction(func(tx *gorm.DB) error {
		return tx.Save(&models.PushTask{
			TaskId:   taskId,
			ResultId: lastResultId,
			Finished: false,
		}).Error
	})
	return err
}

func (pushServiceImp) FinishPushTask(taskId int) error {
	err := DoTransaction(func(tx *gorm.DB) error {
		return tx.Save(&models.PushTask{
			TaskId:   taskId,
			Finished: true,
		}).Error
	})
	return err
}

func (pushServiceImp) SelectPushTasks() ([]*models.PushTask, error) {
	// 查询数据
	var dataQuery = dbConn.Table("t_push_task").Where("finished = ?", false)
	rows, err := dataQuery.Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.PushTask
	for rows.Next() {
		var item models.PushTask
		if err := dataQuery.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		retArray = append(retArray, &item)
	}
	return retArray, nil
}

func (pushServiceImp) SelectPushResults(taskId, size int, lastResultId int64) ([]*models.Result, error) {

	// 查询数据
	var dataQuery = dbConn.Table("t_result").
		Where("task_id = ? and id > ?", taskId, lastResultId).
		Order("id asc").
		Limit(size)
	rows, err := dataQuery.Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.Result
	for rows.Next() {
		var item models.Result
		if err := dataQuery.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		retArray = append(retArray, &item)
	}
	return retArray, nil
}
