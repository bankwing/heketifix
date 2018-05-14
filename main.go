//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package main

import (
	"os"
	"flag"
	"time"

	"github.com/boltdb/bolt"
)

func main() {

        dbPtr := flag.String("db", "", "Heketi db file to fix. (Required)")
        forcePtr := flag.Bool("f", false, "For take action")

        flag.Parse()

        if *dbPtr == "" {
          flag.PrintDefaults()
          os.Exit(1)
        }


        dbfile := *dbPtr
        
        if *forcePtr == false {
          println("***********************************************")
          println("* Running in DRYRUN mode. NO action was taken *")
          println("***********************************************")
        }

	// Setup BoltDB database
	db, err := bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		println(err.Error())
	}

	//defer os.Remove(dbfile)

	// Try to view
	// err = db.View(func(tx *bolt.Tx) error {

	// 	brick_entries, err := BrickList(tx)
	// 	if err != nil {
	// 		println(err.Error())
	// 	}

	// 	for _, brick := range brick_entries {
	// 		println(brick)
	// 	}

	// 	return nil
	// })

	db.Update(func(tx *bolt.Tx) error {

		//Get all bricklist
		brick_entries, err := BrickList(tx)
		if err != nil {
			println(err.Error())
		}

		for _, brickid := range brick_entries {

			//Access brick
			//println(brick.Info.Id)
			brick, err := NewBrickEntryFromId(tx, brickid)
			if err != nil {
				println(err.Error())
				return err
			}

			//Check brick without Path
			//println(brick.Info.Id+" "+brick.Info.Path)
			if brick.Info.Path == "" {
				//println(brick.Info.Id + " No Path -> Fixing")

				// Access device
				device, err := NewDeviceEntryFromId(tx, brick.Info.DeviceId)
				if err != nil {
					println(err.Error())
					return err
				}

                                if *forcePtr == true {
				  v := &VolumeEntry{}
				  v.removeBrickFromDb(tx, brick)
                                }

				if err != nil {
					println(err.Error())
				}

				println("brickId: " + brick.Info.Id + " -> Deleted path on device: " + device.Info.Id )
			} else {
				println("brickId: " +brick.Info.Id + " -> Path Found no action")
			}

		}

		// //Access brick
		// brick, err := NewBrickEntryFromId(tx, "5e74f6fe2c2e522e81e041e3d09ea7aa")
		// if err != nil {
		// 	println(err.Error())
		// 	return err
		// }

		// // Access device
		// device, err := NewDeviceEntryFromId(tx, brick.Info.DeviceId)
		// if err != nil {
		// 	println(err.Error())
		// 	return err
		// }

		// // Deallocate space on device
		// device.StorageFree(brick.TotalSize())

		// // Delete brick from device
		// device.BrickDelete(brick.Info.Id)

		// // Save device
		// err = device.Save(tx)
		// if err != nil {
		// 	println(err.Error())
		// 	return err
		// }

		// // Delete brick entryfrom db
		// err = brick.Delete(tx)
		// if err != nil {
		// 	println(err.Error())
		// 	return err
		// }

		// v := &VolumeEntry{}
		// v.Bricks = make(sort.StringSlice, 0)

		// gob.Register(&NoneDurability{})
		// gob.Register(&VolumeReplicaDurability{})
		// gob.Register(&VolumeDisperseDurability{})

		// // Delete brick from volume db
		// v.BrickDelete(brick.Info.Id)
		// if err != nil {
		// 	println(err.Error())
		// 	return err
		// }

		//d.Bricks = utils.SortedStringsDelete(d.Bricks, id)

		// NewBrickEntry, err := NewBrickEntryFromId(tx, "5e74f6fe2c2e522e81e041e3d09ea7aa")
		// if err != nil {
		// 	println(err.Error())
		// }
		//var brickname = "2dd18ca3c56d86a608edc852e887c4c6"
		//NewBrickEntry := &BrickEntry{}

		// err = EntryDelete(tx, NewBrickEntry, brickname)
		// if err != nil {
		// 	println(err.Error())
		// }
		// println("Delete brick id [%s]:", brickname)

		return nil
	})

	// Verify to view
	err = db.View(func(tx *bolt.Tx) error {

		return nil
	})

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
