///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"digger/models"
	"errors"
	"github.com/jinzhu/gorm"
)

type pluginServiceImp struct {
}

// 创建任务

// 查询队列爬虫任务
func (pluginServiceImp) SelectPluginByName(name string) (*models.Plugin, error) {
	var ret = &models.Plugin{}
	temp := dbConn.Where("name = ?", name).First(ret)
	i := temp.RowsAffected
	err := transformNotFoundErr(temp.Error)
	if i == 0 {
		return nil, err
	}
	return ret, err
}

// 插入队列
func (pluginServiceImp) InsertPlugin(plugin models.Plugin) error {
	plugin.Id = 0
	return DoTransaction(func(tx *gorm.DB) error {
		return tx.Save(&plugin).Error
	})
}

// 根据项目查询插件列表
func (pluginServiceImp) SelectPluginsByProject(projectId int) ([]models.Plugin, error) {
	rows, err := dbConn.Table("t_plugin").
		Where("project_id = ?", projectId).
		Order("id asc").
		Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []models.Plugin
	for rows.Next() {
		var item models.Plugin
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		retArray = append(retArray, item)
	}
	return retArray, nil
}

// 批量保存插件
func (pluginServiceImp) SavePlugins(projectId int, plugins []*models.Plugin) error {
	return DoTransaction(func(tx *gorm.DB) error {

		// 全部删除
		if len(plugins) == 0 {
			// 删除多余的插件
			if err := tx.Delete(models.Plugin{}, "project_id = ?", projectId).Error; err != nil {
				return err
			}
			return nil
		}

		var existIds []int
		for _, p := range plugins {
			if p == nil {
				continue
			}
			r := &models.Plugin{}
			if err := tx.Table("t_plugin").Select("id").Where("project_id = ? and name = ?", projectId, p.Name).First(r).Error; err != nil {
				if transformNotFoundErr(err) != nil {
					return err
				}
			}

			// exists, update
			if r.Id > 0 {
				p.Id = r.Id
			} else {
				p.Id = 0
			}
			p.ProjectId = projectId

			temp := tx.Save(p)
			if temp.RowsAffected == 0 {
				if temp.Error != nil {
					return temp.Error
				}
				return errors.New("update failed")
			}
			if temp.Error != nil {
				return temp.Error
			}
			existIds = append(existIds, p.Id)
		}

		// 删除多余的插件
		if err := tx.Delete(models.Plugin{}, "project_id = ? and id not in (?)", projectId, existIds).Error; err != nil {
			return err
		}
		return nil
	})
}
