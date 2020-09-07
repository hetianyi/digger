///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"digger/models"
	"digger/utils"
	"errors"
	"fmt"
	"github.com/hetianyi/gox/logger"
	"github.com/jinzhu/gorm"
	json "github.com/json-iterator/go"
	"strings"
	"time"
)

type projectServiceImp struct {
}

// 根据项目id查询项目信息
func (projectServiceImp) SelectProjectById(projectId int) (*models.Project, error) {
	var ret = &models.Project{}
	temp := dbConn.Where("id = ?", projectId).First(ret)
	i := temp.RowsAffected
	err := transformNotFoundErr(temp.Error)
	if i == 0 {
		return nil, err
	}

	err = fromDB(ret)
	if err != nil {
		return nil, err
	}

	return ret, err
}

// 根据项目名称查询项目信息
func (projectServiceImp) SelectProjectByName(name string) (*models.Project, error) {
	var ret = &models.Project{}
	temp := dbConn.Where("name = ?", name).First(ret)
	i := temp.RowsAffected
	err := transformNotFoundErr(temp.Error)
	if i == 0 {
		return nil, err
	}

	err = fromDB(ret)
	if err != nil {
		return nil, err
	}

	return ret, err
}

// 根据项目id查询项目完整信息
func (p projectServiceImp) SelectFullProjectInfo(projectId int) (*models.Project, error) {
	project, err := p.SelectProjectById(projectId)
	if err != nil {
		logger.Debug(fmt.Sprintf("error get project data: %s", err.Error()))
		return nil, err
	}
	if project == nil {
		return nil, nil
	}
	stages, err := ProjectConfigService().SelectStages(projectId)
	if err != nil {
		logger.Debug(fmt.Sprintf("error get stage data: %s", err.Error()))
		return nil, err
	}

	fields, err := ProjectConfigService().SelectFields(projectId)
	if err != nil {
		logger.Debug(fmt.Sprintf("error get field data: %s", err.Error()))
		return nil, err
	}

	plugins, err := PluginService().SelectPluginsByProject(projectId)
	if err != nil {
		logger.Debug(fmt.Sprintf("error get plugin data: %s", err.Error()))
		return nil, err
	}

	for i, s := range stages {
		var fs []models.Field
		hasNextStage := false
		for _, f := range fields {
			if f.StageId == s.Id {
				fs = append(fs, f)
				if f.NextStage != "" {
					hasNextStage = true
				}
			}
		}
		stages[i].Fields = fs
		stages[i].HasNextStage = hasNextStage
	}
	project.Stages = stages

	AssembleProjectPlugin(project, plugins)

	return project, nil
}

// 根据条件查询项目列表
func (projectServiceImp) SelectProjectList(params models.ProjectQueryVO) (int64, []*models.Project, error) {

	var baseQuery = func(query *gorm.DB) *gorm.DB {
		if params.Name != "" {
			query = query.Where("name = ?", params.Name)
		}
		if params.DisplayName != "" {
			query = query.Where("display_name = ?", params.DisplayName)
		}
		if params.Tags != nil && len(params.Tags) > 0 {
			query = query.Where("tags in (?)", params.Tags)
		}
		return query
	}

	// 查询总数
	var countQuery = baseQuery(dbConn.Table("t_project"))
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	if total == 0 {
		return 0, nil, nil
	}

	// 查询数据
	var dataQuery = baseQuery(dbConn.Table("t_project"))

	if params.Order == 0 {
		dataQuery = dataQuery.Order("create_time")
	} else {
		dataQuery = dataQuery.Order("create_time desc")
	}
	rows, err := dataQuery.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Rows()
	if transformNotFoundErr(err) != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var retArray []*models.Project
	for rows.Next() {
		var item models.Project
		if err := dataQuery.ScanRows(rows, &item); err != nil {
			return 0, nil, err
		}

		err := fromDB(&item)
		if err != nil {
			return 0, nil, err
		}

		retArray = append(retArray, &item)
	}
	return total, retArray, nil
}

// 新增项目
func (projectServiceImp) CreateProject(project models.Project) (*models.Project, error) {
	project.Id = 0

	err := toDB(&project)
	if err != nil {
		return nil, err
	}
	project.CreateTime = time.Now()
	project.UpdateTime = time.Now()

	err = DoTransaction(func(tx *gorm.DB) error {
		return tx.Save(&project).Error
	})
	return &project, err
}

// 修改项目信息
func (projectServiceImp) UpdateProject(project models.Project) (bool, error) {
	if project.Id <= 0 {
		return false, errors.New("更新项目必须提供Id")
	}
	logger.Debug(fmt.Sprintf("更新项目信息: %d", project.Id))
	success := false

	err := toDB(&project)
	if err != nil {
		return false, err
	}

	project.UpdateTime = time.Now()

	err = DoTransaction(func(tx *gorm.DB) error {
		temp := tx.Save(&project)
		if temp.RowsAffected == 0 {
			success = false
			if temp.Error != nil {
				return temp.Error
			}
			return errors.New("update failed")
		}
		if temp.Error == nil {
			success = true
		}
		return temp.Error
	})
	return success, err
}

// 根据项目名称查询项目信息
func (projectServiceImp) DeleteProject(projectId int) (bool, error) {
	logger.Warn(fmt.Sprintf("正在删除项目: %d", projectId))
	success := false
	err := DoTransaction(func(tx *gorm.DB) error {
		// 删除项目
		temp := tx.Delete(models.Project{}, "id = ?", projectId)
		if temp.RowsAffected == 0 {
			success = false
			if temp.Error != nil {
				return temp.Error
			}
			return errors.New("update failed")
		}
		if temp.Error == nil {
			success = true
		}
		return temp.Error
	})
	return success, err
}

// 查询启用定时任务的项目
func (projectServiceImp) SelectCronProjectList() ([]*models.Project, error) {
	// 查询数据
	rows, err := dbConn.Table("t_project").Where("enable_cron = ? and cron != ''", true).Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.Project
	for rows.Next() {
		var item models.Project
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}

		err := fromDB(&item)
		if err != nil {
			return nil, err
		}

		retArray = append(retArray, &item)
	}
	return retArray, nil
}

func toDB(project *models.Project) error {
	s, err := json.MarshalToString(project.Settings)
	if err != nil {
		return err
	}
	project.SettingsDB = s

	h, err := json.MarshalToString(project.Headers)
	if err != nil {
		return err
	}
	project.HeadersDB = h

	u, err := json.MarshalToString(project.StartUrls)
	if err != nil {
		return err
	}
	project.StartUrlsDB = u

	n, err := json.MarshalToString(project.NodeAffinity)
	if err != nil {
		return err
	}
	project.NodeAffinityDB = n
	return nil
}

func fromDB(project *models.Project) error {
	s := make(map[string]string)
	h := make(map[string]string)
	var n []string
	var us []string
	err := json.UnmarshalFromString(project.SettingsDB, &s)
	if err != nil {
		return err
	}
	project.Settings = s

	err = json.UnmarshalFromString(project.HeadersDB, &h)
	if err != nil {
		return err
	}
	project.Headers = h

	err = json.UnmarshalFromString(project.StartUrlsDB, &us)
	if err != nil {
		return err
	}
	project.StartUrls = us

	err = json.UnmarshalFromString(project.NodeAffinityDB, &n)
	if err != nil {
		return err
	}
	project.NodeAffinity = n
	var kvs []models.KV
	for _, s := range n {
		kv := utils.ParseNodeAffinity(s)
		if kv != nil {
			kvs = append(kvs, *kv)
		}
	}
	project.NodeAffinityParsed = kvs
	return nil
}

func AssembleProjectPlugin(project *models.Project, plugins []models.Plugin) {
	var filterPlugins = func(pluginString string) []models.Plugin {
		ps := strings.Split(pluginString, ",")
		var ret []models.Plugin
		for _, s := range ps {
			plugin := parsePlugin(s)
			if plugin == nil {
				continue
			}
			for _, p := range plugins {
				if plugin.Name == p.Name {
					p.Slot = plugin.Slot
					ret = append(ret, p)
				}
			}
		}
		return ret
	}

	stages := project.Stages

	for i, s := range stages {
		stages[i].Plugins = filterPlugins(stages[i].PluginsDB)
		fields := s.Fields
		for k, f := range fields {
			fps := filterPlugins(f.PluginDB)
			if len(fps) > 0 {
				fields[k].Plugin = &fps[0]
			}
		}
		stages[i].Fields = fields
	}
	project.Stages = stages
}

func parsePlugin(s string) *models.Plugin {
	arr := strings.Split(s, "@")
	if len(arr) != 2 {
		return nil
	}
	name := strings.TrimSpace(arr[0])
	slot := strings.TrimSpace(arr[1])
	if slot != "s1" && slot != "sr" && slot != "s2" && slot != "s3" && slot != "s4" {
		logger.Error("invalid slot ", slot)
		return nil
	}
	return &models.Plugin{
		Name: name,
		Slot: slot,
	}
}

func (t projectServiceImp) AllProjectCount() (int, error) {
	// 查询数据
	type CountResult struct {
		Count int `gorm:"column:count"`
	}
	ret := &CountResult{}
	if err := dbConn.Raw(`SELECT count(*) FROM t_project`).
		Scan(ret).Error; err != nil {
		return 0, err
	}
	return ret.Count, nil
}
