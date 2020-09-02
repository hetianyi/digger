///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"bytes"
	"digger/models"
	"github.com/hetianyi/gox/logger"
	"github.com/jinzhu/gorm"
	"math"
)

type resultServiceImp struct {
}

// 查询任务的结果列表
func (resultServiceImp) SelectResults(params models.ResultQueryVO) (int64, []*models.Result, error) {
	var baseQuery = func(query *gorm.DB) *gorm.DB {
		query = query.Where("task_id = ?", params.TaskId)
		return query
	}

	// 查询总数
	var countQuery = baseQuery(dbConn.Table("t_result"))
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	if total == 0 {
		return 0, nil, nil
	}

	// 查询数据
	var dataQuery = baseQuery(dbConn.Table("t_result"))
	dataQuery = dataQuery.Order("id asc")
	rows, err := dataQuery.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Rows()
	if transformNotFoundErr(err) != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var retArray []*models.Result
	for rows.Next() {
		var item models.Result
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

// 插入结果
func (resultServiceImp) InsertResults(taskId int, results []models.Result) error {
	err := DoTransaction(func(tx *gorm.DB) error {
		for _, r := range results {
			// TODO batch insert
			r.Id = 0
			r.TaskId = taskId
			if err := tx.Save(&r).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// 保存一个check内的数据
// TODO batch insert
func (resultServiceImp) SaveCheckData(checkData *models.QueueCallbackRequestVO, exceedMaxRetryQueueIds []int64) error {
	if checkData == nil {
		return nil
	}
	return DoTransaction(func(tx *gorm.DB) error {
		if len(checkData.Results) > 0 {
			for _, r := range checkData.Results {
				if err := tx.Save(r).Error; err != nil {
					return err
				}
			}
		}
		if len(checkData.NewQueues) > 0 {
			for _, r := range checkData.NewQueues {
				if r.MiddleData == "" {
					r.MiddleData = "{}"
				}
				if err := tx.Save(r).Error; err != nil {
					return err
				}
			}
		}
		if len(checkData.SuccessQueueIds) > 0 {
			if err := tx.Table("t_queue").Where("id in(?)", checkData.SuccessQueueIds).Update("status", 2).Error; err != nil {
				return err
			}
		}
		if len(checkData.ErrorQueueIds) > 0 {
			if err := tx.Table("t_queue").Where("id in(?)", checkData.ErrorQueueIds).Update("status", 0).Error; err != nil {
				return err
			}
		}
		if len(exceedMaxRetryQueueIds) > 0 {
			if err := tx.Table("t_queue").Where("id in(?)", exceedMaxRetryQueueIds).Update("status", 3).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// 查询任务成功结果数量
func (resultServiceImp) ResultCount(taskId ...int) ([]*models.ResultCountCO, error) {
	rows, err := dbConn.Raw(`
SELECT t.id as task_id,
	COUNT ( r.ID ) AS COUNT 
FROM
	t_task t LEFT JOIN
	t_result r ON t.ID = r.task_id 
WHERE
	t.id IN (?) 
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

// 保存处理结果的数据 v2
func (resultServiceImp) SaveProcessResultData(result *models.QueueProcessResult, exceedMaxRetry bool) error {
	if result == nil {
		return nil
	}
	return DoTransaction(func(tx *gorm.DB) error {
		if exceedMaxRetry {
			if err := tx.Table("t_queue").Where("id = ? and status = 0 and expire = ?", result.QueueId, result.Expire).
				Updates(map[string]interface{}{"status": 2, "expire": 0}).Error; err != nil {
				return err
			}
			return nil
		}

		// queue状态改为成功
		temp := tx.Table("t_queue").Where("id = ? and status = 0 and expire = ?", result.QueueId, result.Expire).
			Updates(map[string]interface{}{"status": 1, "expire": 0})
		if err := temp.Error; err != nil {
			return err
		}
		if temp.RowsAffected == 0 {
			// 状态更新失败
			logger.Error("状态更新失败")
			return nil
		}

		if len(result.Results) > 0 {
			for _, r := range result.Results {
				if err := tx.Save(r).Error; err != nil {
					return err
				}
			}
		}
		if len(result.NewQueues) > 0 {
			for _, r := range result.NewQueues {
				if r.MiddleData == "" {
					r.MiddleData = "{}"
				}
			}
			// 如果queue太多，则分批插入
			index := 0
			for {
				if index >= len(result.NewQueues) {
					break
				}
				sql, values := buildBatchInsertQueue(result.NewQueues[index:int64(math.Min(float64(index+5000), float64(len(result.NewQueues))))])
				index += 5000
				r := tx.Exec(sql, values...)
				if err := r.Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}


// 查询从某个id起后续结果总数
func (resultServiceImp) ResultCountSince(id int64) (int, int64, error) {
	// 查询数据
	type CountResult struct {
		Count int `gorm:"column:count"`
	}
	type MaxId struct {
		Id int64 `gorm:"column:max"`
	}
	max := &MaxId{}
	if err := dbConn.Raw(`SELECT max(id) as max FROM t_result`).
		Scan(max).Error; err != nil {
		return 0, id, err
	}
	if max.Id <= id {
		return 0, id, nil
	}
	ret := &CountResult{}
	if err := dbConn.Raw(`SELECT count(*) FROM t_result r WHERE r.id > ? and r.id < ?`, id, max.Id).
		Scan(ret).Error; err != nil {
		return 0, id, err
	}
	return ret.Count, max.Id, nil
}

func buildBatchInsertQueue(queues []*models.Queue) (string, []interface{}) {
	var buff bytes.Buffer
	buff.WriteString("insert into t_queue (task_id, stage_name, url, middle_data, expire) values")
	var values []interface{}
	for i, q := range queues {
		if i == len(queues)-1 {
			buff.WriteString("(?, ?, ?, ?, ?)")
		} else {
			buff.WriteString("(?, ?, ?, ?, ?),")
		}
		values = append(values, q.TaskId, q.StageName, q.Url, q.MiddleData, 0)
	}
	return buff.String(), values
}
