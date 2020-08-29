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

var (
	configsCache = make(map[string]string)
)

type configServiceImpl struct {
}

func (configServiceImpl) ListConfigs() (map[string]string, error) {
	if len(configsCache) > 0 {
		return configsCache, nil
	}

	// 查询数据
	rows, err := dbConn.Table("t_config").Rows()
	if transformNotFoundErr(err) != nil {
		return map[string]string{}, err
	}
	defer rows.Close()

	var cache = make(map[string]string)
	for rows.Next() {
		var item models.Config
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return map[string]string{}, err
		}
		if err != nil {
			return map[string]string{}, err
		}
		cache[item.Key] = item.Value
	}
	configsCache = cache
	return cache, nil
}

func (configServiceImpl) UpdateConfig(key, value string) error {
	err := DoTransaction(func(tx *gorm.DB) error {
		// 保存配置快照
		if err := tx.Save(&models.Config{
			Key:   key,
			Value: value,
		}).Error; err != nil {
			return err
		}
		configsCache[key] = value
		return nil
	})
	return err
}
