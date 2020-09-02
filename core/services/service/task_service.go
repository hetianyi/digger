///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"digger/common"
	"digger/models"
	"digger/utils"
	"errors"
	"fmt"
	"github.com/hetianyi/gox/logger"
	"github.com/jinzhu/gorm"
	json "github.com/json-iterator/go"
)

type taskServiceImp struct {
}

// 创建任务
func (taskServiceImp) CreateTask(task models.Task) (*models.Task, error) {
	project, err := ProjectService().SelectFullProjectInfo(task.ProjectId)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("cannot create task: project not exists")
	}

	config, err := json.MarshalToString(project)
	snapshot := &models.ConfigSnapshot{
		Id:        0,
		ProjectId: task.ProjectId,
		Config:    config,
	}

	err = DoTransaction(func(tx *gorm.DB) error {
		active := 0
		if err := tx.Table("t_task").Where("project_id = ? and status in (?)", task.ProjectId, []int{0, 1}).Count(&active).Error; err != nil {
			return err
		}
		// 一个项目只允许同时激活一个任务
		if active > 0 {
			return errors.New("only one active task allowed")
		}

		// 保存配置快照
		if err := tx.Save(&snapshot).Error; err != nil {
			return err
		}

		// 保存task
		task.ConfigSnapShotId = snapshot.Id
		if err := tx.Save(&task).Error; err != nil {
			return err
		}

		// 保存第一个queue
		firstQueue := &models.Queue{
			Id:         0,
			TaskId:     task.Id,
			StageName:  project.StartStage,
			Url:        project.StartUrl,
			MiddleData: "{}",
		}
		if err := tx.Save(&firstQueue).Error; err != nil {
			return err
		}
		return nil
	})
	return &task, err
}

// 查询任务详情
func (taskServiceImp) SelectTask(id int) (*models.Task, error) {
	var ret = &models.Task{}
	temp := dbConn.Where("id = ?", id).First(ret)
	i := temp.RowsAffected
	err := transformNotFoundErr(temp.Error)
	if i == 0 {
		return nil, err
	}
	return ret, err
}

// 查询任务列表
func (taskServiceImp) SelectTaskList(params models.TaskQueryVO) (int64, []*models.Task, error) {
	var baseQuery = func(query *gorm.DB) *gorm.DB {
		if params.ProjectId > 0 {
			query = query.Where("project_id = ?", params.ProjectId)
		}
		if params.Status > -1 {
			query = query.Where("status = ?", params.Status)
		}
		return query
	}

	// 查询总数
	var countQuery = baseQuery(dbConn.Table("t_task"))
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// 查询数据
	var dataQuery = baseQuery(dbConn.Table("t_task"))
	dataQuery = dataQuery.Order("create_time desc")

	rows, err := dataQuery.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Rows()
	if transformNotFoundErr(err) != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var retArray []*models.Task
	for rows.Next() {
		var item models.Task
		if err := dataQuery.ScanRows(rows, &item); err != nil {
			return 0, nil, err
		}
		retArray = append(retArray, &item)
	}
	return total, retArray, nil
}

// 查询所有已激活的task
func (taskServiceImp) SelectActiveTasks() ([]*models.Task, error) {
	rows, err := dbConn.Table("t_task").Where("status = ?", 1).Order("id asc").Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.Task
	for rows.Next() {
		var item models.Task
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		retArray = append(retArray, &item)
	}
	return retArray, nil
}

// 加载任务的配置快照
func (taskServiceImp) LoadConfigSnapshot(snapshotId int) (*models.Project, error) {
	var ret = &models.ConfigSnapshot{}
	temp := dbConn.Where("id = ?", snapshotId).First(ret)
	i := temp.RowsAffected
	err := transformNotFoundErr(temp.Error)
	if i == 0 {
		return nil, err
	}
	p := &models.Project{}
	err = json.UnmarshalFromString(ret.Config, p)
	if p.Settings == nil {
		p.Settings = make(map[string]string)
	}
	if p.Headers == nil {
		p.Headers = make(map[string]string)
	}
	if p.NodeAffinity == nil {
		p.NodeAffinity = []string{}
	}
	if p.NodeAffinityParsed == nil {
		p.NodeAffinityParsed = []models.KV{}
	}
	return p, err
}

func (taskServiceImp) updateStatus(id, status int) error {
	if status == 2 || status == 3 {
		if err := QueueService().StatisticFinal(id); err != nil {
			return err
		}
	}
	err := DoTransaction(func(tx *gorm.DB) error {
		temp := tx.Model(models.Task{}).Where("id = ?", id).Update("status", status)
		if temp.Error != nil {
			return temp.Error
		}
		i := temp.RowsAffected
		if i == 0 {
			return errors.New("update status failed")
		}
		return nil
	})
	// 当任务停止或完成的时候，需要删除该任务产生的queue
	if status == 2 || status == 3 {
		QueueService().DeleteQueues(id)
	}
	return err
}

func (t taskServiceImp) beforeUpdate(id int) error {
	task, err := t.SelectTask(id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not exists")
	}
	if task.Status == 2 {
		return errors.New("task already shutdown")
	}
	return nil
}

// 完成任务
func (t taskServiceImp) FinishTask(id int) error {
	if err := t.beforeUpdate(id); err != nil {
		return err
	}

	task, err := TaskService().SelectTask(id)
	if err != nil {
		return err
	}

	if err := t.updateStatus(id, 3); err != nil {
		return err
	}

	configs, _ := ConfigService().ListConfigs()
	if configs[common.EMAIL_CONFIG] == "" {
		logger.Info("no email configured, skip email notification")
	} else {
		emailConfig := utils.ParseEmailNotifierStr(configs[common.EMAIL_CONFIG])
		if emailConfig == nil {
			logger.Info("no email configured, skip email notification")
		} else {
			go func() {
				if err := utils.EmailNotify(task, emailConfig); err != nil {
					logger.Error(err)
				}
			}()
		}
	}

	return nil
}

// 关闭任务
func (t taskServiceImp) ShutdownTask(id int) error {
	if err := t.beforeUpdate(id); err != nil {
		return err
	}
	return t.updateStatus(id, 2)
}

// 暂停任务
func (t taskServiceImp) PauseTask(id int) error {
	if err := t.beforeUpdate(id); err != nil {
		return err
	}
	return t.updateStatus(id, 0)
}

// 开启任务
func (t taskServiceImp) StartTask(id int) error {
	if err := t.beforeUpdate(id); err != nil {
		return err
	}
	return t.updateStatus(id, 1)
}

// 查询项目结果数量
func (t taskServiceImp) TaskCount(projectIds ...int) ([]*models.TaskCountCO, error) {
	rows, err := dbConn.Raw(`
SELECT
	( SELECT COUNT ( * ) AS pause_count FROM t_task A WHERE A.project_id = P.ID AND A.status = 0 ),
	( SELECT COUNT ( * ) AS active_count FROM t_task A WHERE A.project_id = P.ID AND A.status = 1 ),
	( SELECT COUNT ( * ) AS stop_count FROM t_task A WHERE A.project_id = P.ID AND A.status = 2 ),
	( SELECT COUNT ( * ) AS finish_count FROM t_task A WHERE A.project_id = P.ID AND A.status = 3 ),
	P.ID AS project_id 
FROM
	t_project
	P LEFT JOIN t_task T ON P.ID = T.project_id 
WHERE
	P.ID IN (?)
GROUP BY
	P.ID
`, projectIds).Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.TaskCountCO
	for rows.Next() {
		var item models.TaskCountCO
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		retArray = append(retArray, &item)
	}
	return retArray, nil
}

// 检查task是否完成
func (t taskServiceImp) CheckTaskFinish(taskId int) (bool, error) {
	// 查询总数
	var countQuery = dbConn.Table("t_queue").Where("task_id = ? and status = ?", taskId, 0)
	var total int64 = 1
	if err := countQuery.Count(&total).Error; err != nil {
		return false, err
	}
	return total == 0, nil
}

// 删除task
func (t taskServiceImp) DeleteTask(taskId int) error {
	logger.Warn(fmt.Sprintf("正在删除任务: %d", taskId))
	return DoTransaction(func(tx *gorm.DB) error {
		// 删除项目
		temp := tx.Delete(models.Task{}, "id = ? and status in(2,3)", taskId)
		if temp.RowsAffected == 0 {
			if temp.Error != nil {
				return temp.Error
			}
			return errors.New("task delete failed")
		}
		if temp.Error != nil {
			return temp.Error
		}

		// 删除queue
		logger.Info("删除queue...")
		if err := tx.Delete(models.Queue{}, "task_id = ?", taskId).Error; err != nil {
			return err
		}
		// 删除result
		logger.Info("删除result...")
		if err := tx.Delete(models.Result{}, "task_id = ?", taskId).Error; err != nil {
			return err
		}
		return nil
	})
}

// 查询所有任务数量
func (t taskServiceImp) AllTaskCount() (int, error) {
	// 查询数据
	type CountResult struct {
		Count int `gorm:"column:count"`
	}
	ret := &CountResult{}
	if err := dbConn.Raw(`SELECT count(*) FROM t_task`).
		Scan(ret).Error; err != nil {
		return 0, err
	}
	return ret.Count, nil
}