///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"digger/models"
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"github.com/jinzhu/gorm"
	"time"
)

type queueServiceImpl struct {
}

// 查询任务的结果列表
func (queueServiceImpl) SelectQueues(params models.QueueQueryVO) ([]*models.Queue, error) {
	rows, err := dbConn.Table("t_queue").
		Where("task_id = ? and status = ? and expire = 0", params.TaskId, params.Status).
		Order("id asc").
		Limit(params.Limit).
		Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.Queue
	var retIds []int64
	for rows.Next() {
		var item models.Queue
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		//item.Expire = now
		retIds = append(retIds, item.Id)
		retArray = append(retArray, &item)
	}

	if len(retArray) == 0 {
		return nil, nil
	}

	now := gox.GetTimestamp(time.Now().Add(time.Second * 10 * time.Duration(len(retArray))))
	for _, v := range retArray {
		v.Expire = now
	}

	// 更新状态为处理中
	le := len(retIds)
	if params.LockStatus && le > 0 {
		if err := dbConn.Table("t_queue").Where("id in(?) and expire = 0", retIds).
			Update("expire", now).Error; err != nil {
			return nil, err
		}
	}
	return retArray, nil
}

// 插入结果
func (queueServiceImpl) InsertQueue(queue models.Queue) error {
	err := DoTransaction(func(tx *gorm.DB) error {
		// 保存配置快照
		if err := tx.Save(&queue).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// 重置处理中的queue状态为未处理
func (queueServiceImpl) ResetQueuesStatus() error {
	now := gox.GetTimestamp(time.Now())
	err := DoTransaction(func(tx *gorm.DB) error {
		temp := tx.Table("t_queue").Where("expire > ? and expire < ?", 0, now).Update("expire", 0)
		if err := temp.Error; err != nil {
			return err
		}
		if temp.RowsAffected > 0 {
			logger.Info(fmt.Sprintf("已重置%d条queue记录", temp.RowsAffected))
		}
		return nil
	})
	return err
}

// 查询任务未完成的queue数量
func (queueServiceImpl) GetUnFinishedCount(taskId int) (int, error) {
	// 查询数据
	type CountResult struct {
		Count int `gorm:"column:count"`
	}
	ret := &CountResult{}
	if err := dbConn.Raw(`select count(*) as count from t_queue a where a.task_id = ? and a.status = 0`, taskId, taskId).
		Scan(ret).Error; err != nil {
		return 0, err
	}
	return ret.Count, nil
}

// 查询任务失败queue数量
func (queueServiceImpl) ErrorCount(taskId ...int) ([]*models.ResultCountCO, error) {
	rows, err := dbConn.Raw(`
SELECT t.id as task_id,
	COUNT ( q.ID ) AS COUNT 
FROM
	t_task t LEFT JOIN
	t_queue q ON t.ID = q.task_id 
WHERE
	t.id IN (?) and q.status = 2 
GROUP BY
	t.id
`, taskId).Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.ResultCountCO
	for rows.Next() {
		var item models.ResultCountCO
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		retArray = append(retArray, &item)
	}
	return retArray, nil
}

// 根据任务id删除queue
func (queueServiceImpl) DeleteQueues(taskId int) error {
	// delete queues
	logger.Info("deleting queues of task ", taskId)
	if err := dbConn.Delete(models.Queue{}, "task_id = ?", taskId).Error; err != nil {
		return err
	}
	logger.Info("finish deleting queues of task ", taskId)
	return nil
}

// 结束时统计最终错误数
func (queueServiceImpl) StatisticFinal(taskId int) error {
	// delete queues
	logger.Info("doing final statistic work for queue, this may take a while...", taskId)
	// 查询数据
	type CountResult struct {
		Count int `gorm:"column:count"`
	}
	ret := &CountResult{}
	if err := dbConn.Raw(`select count(*) as count from t_queue a where a.task_id = ?`, taskId).
		Scan(ret).Error; err != nil {
		return err
	}
	total := ret.Count
	if err := dbConn.Raw(`select count(*) as count from t_queue a where a.task_id = ? and a.status = ?`, taskId, 2).
		Scan(ret).Error; err != nil {
		return err
	}
	errorCount := ret.Count
	if err := dbConn.Raw(`select count(*) as count from t_result a where a.task_id = ?`, taskId).
		Scan(ret).Error; err != nil {
		return err
	}
	resultCount := ret.Count

	err := DoTransaction(func(tx *gorm.DB) error {
		temp := tx.Table("t_task").Where("id = ?", taskId).
			Update(map[string]interface{}{"result_count": resultCount, "success_request": total - errorCount, "error_request": errorCount})
		if err := temp.Error; err != nil {
			return err
		}
		return nil
	})
	return err
}


