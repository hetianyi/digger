///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"digger/models"
	"github.com/jinzhu/gorm"
	"time"
)

type proxyServiceImp struct {
}

func (p proxyServiceImp) Save(proxy models.Proxy) error {
	old, err := p.selectByAddress(proxy.Address)
	if err != nil {
		return err
	}

	if old != nil {
		proxy.Id = old.Id
		proxy.CreateTime = old.CreateTime
	} else {
		proxy.CreateTime = time.Now()
	}
	err = DoTransaction(func(tx *gorm.DB) error {
		return tx.Save(&proxy).Error
	})
	return err
}

func (proxyServiceImp) selectByAddress(address string) (*models.Proxy, error) {
	var ret models.Proxy
	temp := dbConn.Table("t_proxy").Where("address = ?", address).First(&ret)
	i := temp.RowsAffected
	err := transformNotFoundErr(temp.Error)
	if i == 0 {
		return nil, err
	}
	return &ret, nil
}

func (proxyServiceImp) List(params *models.ProxyQueryVO) (int64, []*models.Proxy, error) {
	var baseQuery = func(query *gorm.DB) *gorm.DB {
		if params.Key != "" {
			query = query.Where("address like '%' || ? || '%' or remark like '%' || ? || '%'", params.Key, params.Key)
		}
		return query
	}

	var total int64 = 0
	// 查询总数
	var countQuery = baseQuery(dbConn.Table("t_proxy"))
	if err := countQuery.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	if total == 0 {
		return 0, nil, nil
	}

	// 查询数据
	var dataQuery = baseQuery(dbConn.Table("t_proxy"))
	dataQuery = dataQuery.Order("id desc")
	rows, err := dataQuery.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Rows()
	if transformNotFoundErr(err) != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var retArray []*models.Proxy
	for rows.Next() {
		var item models.Proxy
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

func (proxyServiceImp) Delete(idList []int) error {
	return DoTransaction(func(tx *gorm.DB) error {
		return tx.Delete(models.Proxy{}, "id in (?)", idList).Error
	})
}

func (proxyServiceImp) SelectByProject(projectId int) ([]*models.Proxy, error) {
	rows, err := dbConn.Raw(`
select b.* from t_project_proxy a 
left join t_proxy b on a.proxy_id = b.id
where a.project_id = ?
`, projectId).Rows()
	if transformNotFoundErr(err) != nil {
		return nil, err
	}
	defer rows.Close()

	var retArray []*models.Proxy
	for rows.Next() {
		var item models.Proxy
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return nil, err
		}
		retArray = append(retArray, &item)
	}
	return retArray, nil
}
