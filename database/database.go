package database

import (
	"encoding/json"

	"github.com/arthurbailao/suno-kindle/domain"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type DB struct {
	boltdb *bolt.DB
}

func Open(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open boltdb file")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, e := tx.Cursor().Bucket().CreateBucketIfNotExists([]byte("DEVICES"))
		return e
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create buckets")
	}

	return &DB{db}, nil

}

func (db DB) CreateOrUpdateDevice(k domain.Device) error {
	return db.boltdb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("DEVICES"))

		buf, err := json.Marshal(k)
		if err != nil {
			return errors.Wrap(err, "failed to json marshal")
		}

		err = bucket.Put([]byte(k.Email), buf)
		if err != nil {
			return errors.Wrap(err, "failed to put into boltdb")
		}

		return nil
	})
}

func (db DB) ListDevices() ([]domain.Device, error) {
	var devices []domain.Device
	err := db.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("DEVICES"))

		return bucket.ForEach(func(k, v []byte) error {
			var d domain.Device
			err := json.Unmarshal(v, &d)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal json")
			}

			devices = append(devices, d)
			return nil
		})
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to complete view transaction")
	}
	return devices, nil
}
