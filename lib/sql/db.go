package sql

import (
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

type Config struct {
	URL   string
	Debug bool
}

func openDB(debug bool, dialect, url string) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	switch dialect {
	case "mysql":
		db, err = gorm.Open(mysql.Open(url), &gorm.Config{})
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(url), &gorm.Config{})
	default:
		db, err = nil, errors.New("unsupported dialect")
	}

	if err != nil {
		return nil, err
	}

	if debug {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else {
		db.Logger = db.Logger.LogMode(logger.Silent)
	}

	return db, nil
}

type DBOptions struct {
	MaxIdleConns int
	MaxOpenConns int
	Debug        bool
}

type DBOption func(o *DBOptions)

func Open(dialect, dburl string, opts ...DBOption) (*DB, error) {
	var o = DBOptions{
		MaxIdleConns: 10,
		MaxOpenConns: 128,
		Debug:        false,
	}

	for _, opt := range opts {
		opt(&o)
	}

	db, err := openDB(o.Debug, dialect, dburl)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	DB := &DB{
		DB: db,
	}

	sqlDB.SetMaxIdleConns(o.MaxIdleConns)
	sqlDB.SetMaxOpenConns(o.MaxOpenConns)

	return DB, nil
}

type WhereFunc func(q *gorm.DB) *gorm.DB

func UpdateModelInTx[T interface{}](
	tx *gorm.DB,
	whereFunc WhereFunc,
	updateFunc func(model *T) error,
) error {
	model, err := GetModelInTx[T](tx, whereFunc, true)
	if err != nil {
		return err
	}

	if model == nil {
		return gorm.ErrRecordNotFound
	}

	if err := updateFunc(model); err != nil {
		return err
	}

	if err := tx.Save(model).Error; err != nil {
		return err
	}

	return nil
}

func UpdateModel[T interface{}](db *gorm.DB,
	whereFunc WhereFunc,
	updateFunc func(model *T) error,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return UpdateModelInTx(tx, whereFunc, updateFunc)
	})
}

func UpdateModelsInTx[T interface{}](
	tx *gorm.DB,
	whereFunc WhereFunc,
	updateFunc func(model []*T) error,
) error {
	models, err := GetModelsInTx[T](tx, whereFunc, true)
	if err != nil {
		return err
	}

	if err := updateFunc(models); err != nil {
		return nil
	}

	return tx.Save(models).Error
}

func UpdateModels[T interface{}](
	db *gorm.DB,
	whereFunc WhereFunc,
	updateFunc func(model []*T) error,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return UpdateModelsInTx(tx, whereFunc, updateFunc)
	})
}

func CreateModel[T interface{}](db *gorm.DB, model *T) error {
	return db.Save(model).Error
}

func GetModel[T any](
	tx *gorm.DB,
	whereFunc WhereFunc,
) (*T, error) {
	return GetModelInTx[T](tx, whereFunc, false)
}

func GetModels[T any](
	tx *gorm.DB,
	whereFunc WhereFunc,
) ([]*T, error) {
	return GetModelsInTx[T](tx, whereFunc, false)
}

func GetModelInTx[T interface{}](
	tx *gorm.DB,
	whereFunc WhereFunc,
	forUpdate bool,
) (*T, error) {
	var (
		model T
	)

	query := whereFunc(tx.Model(&model))
	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	err := query.First(&model).Error
	found, err := FoundRecord(err)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &model, nil
}

func GetModelsInTx[T interface{}](
	tx *gorm.DB,
	whereFunc WhereFunc,
	forUpdate bool,
) ([]*T, error) {
	var (
		model  T
		models []*T
	)

	query := whereFunc(tx.Model(&model))
	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

func PageModel[T interface{}](
	tx *gorm.DB, whereFunc WhereFunc, pageIndex, pageSize int,
) (models []*T, count int64, err error) {
	var (
		model T
	)

	query := whereFunc(tx.Model(&model))
	if err = query.Count(&count).Error; err != nil {
		return
	}

	if err = query.Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&models).Error; err != nil {
		return
	}

	return
}
func CountModel[T interface{}](
	tx *gorm.DB, whereFunc WhereFunc,
) (count int64, err error) {
	var (
		model T
	)

	query := whereFunc(tx.Model(&model))
	if err = query.Count(&count).Error; err != nil {
		return
	}
	return
}

func AutoMigrate[T interface{}](db *gorm.DB, model ...*T) error {
	return db.AutoMigrate(model)
}
