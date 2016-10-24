package service

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/uber-go/zap"
)

// DB struct
type DB struct {
	dbname     string
	db         *bolt.DB
	bucketName []byte
}

func NewDB(dbname string, bucketname string) *DB {
	db := &DB{
		dbname:     dbname,
		bucketName: []byte(bucketname),
	}
	return db
}

func (d *DB) Init() {
	db, err := bolt.Open(d.dbname, 0600, nil)
	if err != nil {
		VLogger.Fatal("DB", zap.Error(err))
	}
	d.db = db
}

func (d *DB) Insert(gr *GroupRaw) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(d.bucketName)
		if err != nil {
			VLogger.Error("Insert", zap.Error(err))
			return fmt.Errorf("create bucket error: %s", err)
		}

		bv, err := json.Marshal(gr)
		if err != nil {
			VLogger.Error("Marshal", zap.Error(err))
			return err
		}
		err = b.Put([]byte(gr.ID), bv)
		if err != nil {
			VLogger.Error("Put", zap.Error(err))
			return err
		}
		return nil
	})
}

func (d *DB) isExist() (bool, error) {
	err := d.db.View(func(tx *bolt.Tx) error {
		if tx.Bucket(d.bucketName) == nil {
			return fmt.Errorf("%s groups_not_exist", string(d.bucketName))
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
