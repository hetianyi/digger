///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

type dbServiceImp struct {
}

// 创建任务

// 查询队列爬虫任务
func (dbServiceImp) Check() error {
	/*sql := `select table_name from information_schema.tables a where a.table_catalog='test' and a.table_type='BASE TABLE' and table_schema='public'`
	rows, err := dbConn.Raw(sql).Rows()
	if transformNotFoundErr(err) != nil {
		logger.Fatal(err)
		return errors.New("cannot get table schema: " + err.Error())
	}
	defer rows.Close()

	var retArray []string
	for rows.Next() {
		var item string
		if err := dbConn.ScanRows(rows, &item); err != nil {
			return err
		}
		retArray = append(retArray, item)
	}

	var missTab []string

	if len(retArray) == 0 {

	}
	return retArray, nil*/
	return nil
}
