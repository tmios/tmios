package sql

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

func FoundRecord(err error) (bool, error) {
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err == nil {
		return true, nil
	}

	return false, err
}

func ExistCheck[T interface{}](
	tx *gorm.DB,
	whereFunc WhereFunc,
	retErr error,
) error {
	var (
		model T
		count int64
	)

	if err := whereFunc(tx.Model(&model)).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return retErr
	}

	return nil
}

func NotExistCheck[T interface{}](
	tx *gorm.DB,
	whereFunc WhereFunc,
	retErr error,
) error {
	var (
		model T
		count int64
	)

	if err := whereFunc(tx.Model(&model)).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return retErr
	}

	return nil
}

func MysqlUrl(dbName, username, password, ip string, port int) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, ip, port, dbName)
}

func CreateDatabase(dbName, username, password, ip string, port int) {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, ip, port))
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		if db != nil {
			err := db.Close()
			if err != nil {
				logrus.Warn(err)
			}
		}
	}(db)

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + " DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci")
	if err != nil {
		panic(err)
	}
	logrus.Info("success")
}
