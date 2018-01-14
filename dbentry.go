//
// Copyright (c) 2016 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package main

import (
	"github.com/boltdb/bolt"
	//	"github.com/heketi/heketi/apps/glusterfs"

	"errors"

	"github.com/heketi/heketi/pkg/glusterfs/api"
	"github.com/lpabon/godbc"
)

const (
	ASYNC_ROUTE                    = "/queue"
	BOLTDB_BUCKET_CLUSTER          = "CLUSTER"
	BOLTDB_BUCKET_NODE             = "NODE"
	BOLTDB_BUCKET_VOLUME           = "VOLUME"
	BOLTDB_BUCKET_DEVICE           = "DEVICE"
	BOLTDB_BUCKET_BRICK            = "BRICK"
	BOLTDB_BUCKET_BLOCKVOLUME      = "BLOCKVOLUME"
	BOLTDB_BUCKET_DBATTRIBUTE      = "DBATTRIBUTE"
	DB_CLUSTER_HAS_FILE_BLOCK_FLAG = "DB_CLUSTER_HAS_FILE_BLOCK_FLAG"
)

var (
	NoSpace          = errors.New("No space")
	Found            = errors.New("Id already exists")
	NotFound         = errors.New("Id not found")
	Conflict         = errors.New("The target exists, contains other items, or is in use.")
	MaxBricks        = errors.New("Maximum number of bricks reached.")
	MinimumBrickSize = errors.New("Minimum brick size limit reached.  Out of space.")
	DbAccess         = errors.New("Unable to access db")
	AccessList       = errors.New("Unable to access list")
	KeyExists        = errors.New("Key already exists in the database")
	NoReplacement    = errors.New("No Replacement was found for resource requested to be removed")
)

type testDbEntry struct {
}

func (t *testDbEntry) BucketName() string {
	return "TestBucket"
}

func (t *testDbEntry) Marshal() ([]byte, error) {
	return nil, nil
}

func (t *testDbEntry) Unmarshal(data []byte) error {
	return nil
}

type Entry struct {
	State api.EntryState
}

func (e *Entry) isOnline() bool {
	return e.State == api.EntryStateOnline
}

func (e *Entry) SetOnline() {
	e.State = api.EntryStateOnline
}

type DbEntry interface {
	BucketName() string
	Marshal() ([]byte, error)
	Unmarshal(buffer []byte) error
}

// type DeviceEntry struct {
// 	Entry

// 	Info       api.DeviceInfo
// 	Bricks     sort.StringSlice
// 	NodeId     string
// 	ExtentSize uint64
// }

// func DeviceList(tx *bolt.Tx) ([]string, error) {

// 	list := EntryKeys(tx, BOLTDB_BUCKET_DEVICE)
// 	if list == nil {
// 		return nil, AccessList
// 	}
// 	return list, nil
// }

// func NewDeviceEntry() *DeviceEntry {
// 	entry := &DeviceEntry{}
// 	entry.Bricks = make(sort.StringSlice, 0)
// 	entry.SetOnline()

// 	// Default to 4096KB
// 	entry.ExtentSize = 4096

// 	return entry
// }

// func (d *DeviceEntry) BucketName() string {
// 	return BOLTDB_BUCKET_DEVICE
// }

// func (d *DeviceEntry) Marshal() ([]byte, error) {
// 	var buffer bytes.Buffer
// 	enc := gob.NewEncoder(&buffer)
// 	err := enc.Encode(*d)

// 	return buffer.Bytes(), err
// }

// func (d *DeviceEntry) Unmarshal(buffer []byte) error {
// 	dec := gob.NewDecoder(bytes.NewReader(buffer))
// 	err := dec.Decode(d)
// 	if err != nil {
// 		return err
// 	}

// 	// Make sure to setup arrays if nil
// 	if d.Bricks == nil {
// 		d.Bricks = make(sort.StringSlice, 0)
// 	}

// 	return nil
// }

// Checks if the key already exists in the database.  If it does not exist,
// then it will save the key value pair in the database bucket.
func EntryRegister(tx *bolt.Tx, entry DbEntry, key string, value []byte) ([]byte, error) {
	godbc.Require(tx != nil)
	godbc.Require(len(key) > 0)

	// Access bucket
	b := tx.Bucket([]byte(entry.BucketName()))
	if b == nil {
		err := DbAccess
		println(err)
		return nil, err
	}

	// Check if key exists already
	val := b.Get([]byte(key))
	if val != nil {
		return val, KeyExists
	}

	// Key does not exist.  We can save it
	err := b.Put([]byte(key), value)
	if err != nil {
		println(err)
		return nil, err
	}

	return nil, nil
}

func EntryKeys(tx *bolt.Tx, bucket string) []string {
	list := make([]string, 0)

	// Get all the cluster ids from the DB
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return nil
	}

	err := b.ForEach(func(k, v []byte) error {
		list = append(list, string(k))
		return nil
	})
	if err != nil {
		return nil
	}

	return list
}

func EntrySave(tx *bolt.Tx, entry DbEntry, key string) error {
	godbc.Require(tx != nil)
	godbc.Require(len(key) > 0)

	// Access bucket
	b := tx.Bucket([]byte(entry.BucketName()))
	if b == nil {
		err := DbAccess
		println(err)
		return err
	}

	// Save device entry to db
	buffer, err := entry.Marshal()
	if err != nil {
		println(err)
		return err
	}

	// Save data using the id as the key
	err = b.Put([]byte(key), buffer)
	if err != nil {
		println(err)
		return err
	}

	return nil
}

func EntryDelete(tx *bolt.Tx, entry DbEntry, key string) error {
	godbc.Require(tx != nil)
	godbc.Require(len(key) > 0)

	// Access bucket
	b := tx.Bucket([]byte(entry.BucketName()))
	if b == nil {
		err := DbAccess
		println(err)
		return err
	}

	// Delete key
	err := b.Delete([]byte(key))
	if err != nil {
		println("Unable to delete key [%v] in db: ", key)
		return err
	}

	return nil
}

func EntryLoad(tx *bolt.Tx, entry DbEntry, key string) error {
	godbc.Require(tx != nil)
	godbc.Require(len(key) > 0)

	b := tx.Bucket([]byte(entry.BucketName()))
	if b == nil {
		err := DbAccess
		println(err)
		return err
	}

	val := b.Get([]byte(key))
	if val == nil {
		return NotFound
	}

	err := entry.Unmarshal(val)
	if err != nil {
		println(err)
		return err
	}

	return nil
}

// func NewDeviceEntryFromId(tx *bolt.Tx, id string) (*DeviceEntry, error) {
// 	godbc.Require(tx != nil)

// 	entry := NewDeviceEntry()
// 	err := EntryLoad(tx, entry, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return entry, nil
// }

func xxmain() {
	//tmpfile := "heketi.db"

	// // Setup BoltDB database
	// db, err := bolt.Open(tmpfile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	// if err != nil {
	// 	println(err.Error())
	// }

	// //defer os.Remove(tmpfile)

	// // Try to view
	// err = db.View(func(tx *bolt.Tx) error {

	// 	newDeviceEntry, err := NewDeviceEntryFromId(tx, "0b67d12089a50012862ee95778a5185d")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	//println(string(newDeviceEntry.Bricks))
	// 	fmt.Printf("%s\n", newDeviceEntry.Bricks)

	// 	return nil
	// })

	// // Create a bucket
	// entry := &testDbEntry{}
	// err = db.Update(func(tx *bolt.Tx) error {

	// 	// Create Cluster Bucket
	// 	_, err := tx.CreateBucketIfNotExists([]byte(entry.BucketName()))
	//     if err != nil {
	//             println(err.Error())
	//     }

	// 	// Register a value
	// 	_, err = EntryRegister(tx, entry, "mykey", []byte("myvalue"))
	//     if err != nil {
	//             println(err.Error())
	//     }

	// 	return nil
	// })

	// // Try to write key again
	// err = db.Update(func(tx *bolt.Tx) error {

	// 	// Save again, it should not work
	// 	val, err := EntryRegister(tx, entry, "mykey", []byte("myvalue"))
	//     if err != nil {
	//             println(err.Error())
	//     }

	//             println(string(val))

	//             newDeviceEntry, err := NewDeviceEntryFromId(tx, d.Id())
	//             if err != nil {
	//                     return err
	//             }

	// 	// Remove key
	// 	err = EntryDelete(tx, entry, "mykey")

	// 	// Register again
	// 	_, err = EntryRegister(tx, entry, "mykey", []byte("myvalue"))

	// 	return nil
	// })

}
