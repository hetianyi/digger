///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"digger/models"
	"github.com/jinzhu/gorm"
	json "github.com/json-iterator/go"
	"time"
)

type statisticServiceImp struct {
}

func (statisticServiceImp) Save(data map[string]interface{}) error {
	d, err := json.MarshalToString(data)
	if err != nil {
		return err
	}
	statistic := &models.Statistic{
		Data:       d,
		CreateTime: time.Now(), //data["time"].(time.Time),
	}
	err = DoTransaction(func(tx *gorm.DB) error {
		// 保存配置快照
		if err := tx.Save(statistic).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (statisticServiceImp) List(start time.Time, end time.Time) ([]*models.StatisticVO, error) {
	rows, err := dbConn.Table("t_statistic").
		Where("create_time BETWEEN ? and ?", start, end).
		Order("id asc").
		Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.StatisticVO
	for rows.Next() {
		var item models.Statistic
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		retArray = append(retArray, models.StatisticVO{}.From(&item))
	}
	return retArray, nil
}
