///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"digger/models"
	"errors"
	"github.com/hetianyi/gox/logger"
	"github.com/jinzhu/gorm"
	"time"
)

type projectConfigServiceImp struct {
}

// 根据项目id查询阶段列表
func (projectConfigServiceImp) SelectStages(projectId int) ([]models.Stage, error) {
	rows, err := dbConn.Table("t_stage").Where("project_id = ?", projectId).Order("id asc").Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []models.Stage
	for rows.Next() {
		var item models.Stage
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		retArray = append(retArray, item)
	}
	return retArray, nil
}

// 根据项目id查询参数列表
func (projectConfigServiceImp) SelectFields(projectId int) ([]models.Field, error) {
	rows, err := dbConn.Table("t_field").Where("project_id = ?", projectId).Order("id asc").Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []models.Field
	for rows.Next() {
		var item models.Field
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		retArray = append(retArray, item)
	}
	return retArray, nil
}

// 保存项目字段数据
func (projectConfigServiceImp) SaveStagesAndFields(projectId int, stages []models.Stage) error {
	err := DoTransaction(func(tx *gorm.DB) error {
		// delete stages and fields
		if err := tx.Delete(models.Stage{}, "project_id = ?", projectId).Error; err != nil {
			return err
		}
		if err := tx.Delete(models.Field{}, "project_id = ?", projectId).Error; err != nil {
			return err
		}
		for _, s := range stages {
			s.Id = 0
			s.ProjectId = projectId
			if err := tx.Save(&s).Error; err != nil {
				return err
			}

			if len(s.Fields) == 0 {
				continue
			}

			for _, f := range s.Fields {
				f.Id = 0
				f.ProjectId = projectId
				f.StageId = s.Id
				if err := tx.Save(&f).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

// 从配置文件解析的数据保存
func (projectConfigServiceImp) SaveProjectConfig(project *models.Project, stages []models.Stage) error {
	err := toDB(project)
	if err != nil {
		return err
	}

	err = DoTransaction(func(tx *gorm.DB) error {
		project.UpdateTime = time.Now()
		temp := tx.Save(&project)
		if temp.Error != nil {
			return temp.Error
		}
		if temp.RowsAffected == 0 {
			return errors.New("error updating project")
		}

		// delete stages and fields
		if err := tx.Delete(models.Stage{}, "project_id = ?", project.Id).Error; err != nil {
			return err
		}
		if err := tx.Delete(models.Field{}, "project_id = ?", project.Id).Error; err != nil {
			return err
		}
		for _, s := range stages {
			s.Id = 0
			s.ProjectId = project.Id
			if err := tx.Save(&s).Error; err != nil {
				return err
			}

			if len(s.Fields) == 0 {
				continue
			}

			for _, f := range s.Fields {
				f.Id = 0
				f.ProjectId = project.Id
				f.StageId = s.Id
				if err := tx.Save(&f).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

// 导入项目配置
func (projectConfigServiceImp) ImportProjectConfig(project *models.Project) error {
	if project == nil {
		return nil
	}

	return DoTransaction(func(tx *gorm.DB) error {

		// 删除plugins
		if err := tx.Delete(models.Plugin{}, "project_id = ?", project.Id).Error; err != nil {
			return err
		}
		// 删除stages
		if err := tx.Delete(models.Stage{}, "project_id = ?", project.Id).Error; err != nil {
			return err
		}
		// 删除fields
		if err := tx.Delete(models.Field{}, "project_id = ?", project.Id).Error; err != nil {
			return err
		}

		if err := toDB(project); err != nil {
			return err
		}
		if err := tx.Table("t_project").Where("id = ?", project.Id).
			Updates(map[string]interface{}{
				"settings":      project.SettingsDB,
				"headers":       project.HeadersDB,
				"node_affinity": project.NodeAffinityDB,
				"start_url":     project.StartUrl,
				"start_stage":   project.StartStage,
			}).Error; err != nil {
			return err
		}

		importedPluginNames := make(map[string]byte)
		for _, s := range project.Stages {
			s.Id = 0
			s.ProjectId = project.Id
			logger.Info("导入stage：" + s.Name)
			if err := tx.Save(&s).Error; err != nil {
				return err
			}

			for _, p := range s.Plugins {
				if importedPluginNames[p.Name] == 1 {
					continue
				}
				p.Id = 0
				p.ProjectId = project.Id
				logger.Info("导入plugin：" + p.Name)
				if err := tx.Save(&p).Error; err != nil {
					return err
				}
				importedPluginNames[p.Name] = 1
			}

			if len(s.Fields) == 0 {
				continue
			}

			for _, f := range s.Fields {
				f.Id = 0
				f.ProjectId = project.Id
				f.StageId = s.Id
				logger.Info("导入field：" + f.Name)
				if err := tx.Save(&f).Error; err != nil {
					return err
				}
				if f.Plugin != nil {
					if importedPluginNames[f.Plugin.Name] == 1 {
						continue
					}
					f.Plugin.Id = 0
					f.Plugin.ProjectId = project.Id
					logger.Info("导入plugin：" + f.Plugin.Name)
					if err := tx.Save(&f.Plugin).Error; err != nil {
						return err
					}
					importedPluginNames[f.Plugin.Name] = 1
				}
			}
		}
		return nil
	})
}
