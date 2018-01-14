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
	"time"

	"github.com/boltdb/bolt"
)

func main() {
	tmpfile := "heketi.db"

	// Setup BoltDB database
	db, err := bolt.Open(tmpfile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		println(err.Error())
	}

	//defer os.Remove(tmpfile)

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

		// brick_entries := []string{
		// 	"34591efd238d9f7ac21bb13490a7a283",
		// 	"6402662654ed481bfe50fbf8fd37460a",
		// 	"9a234bd84cc8b0e1e6d1aa7efb886170",
		// 	"f9d3eedeb6f19647e345bc3a411c41e0",
		// 	"023ded95d15898e2ef97c765ca4bbfd6",
		// 	"15d5e1a7e780d2903dc4c995477992cb",
		// 	"1d53b681a2035567d6e20013a2fd560c",
		// 	"454430e4786eb2357c9322366aafd139",
		// 	"b92ac58a05ec0d9aa3046b413e52e180",
		// 	"e2a701a6196e1c9e7d5b3b4413133807",
		// 	"e2dbfd2ed83a2f7e23792a4eb08a712e",
		// 	"465e468150ef40754d4f61a1f2258a6c",
		// 	"6a295b4645ad01c14308c572ddce5767",
		// 	"e7f42751851066df7a03513f50869e2a",
		// 	"f65beb18530489094f6566f17d9ec669",
		// 	"3d0844fad24bf22aa675fa5768675d7f",
		// 	"960e6398e4e96dc9a351be32bcee357c",
		// 	"a73c64b1146e6bf59a8f7ef64554c264",
		// 	"d1657bbe50e590f41abaed7e343ed1a6",
		// 	"d90dd853c78330fa8e7552b2a3998aef",
		// 	"ea9dbac7ab7feb44fc615cba9cef547d",
		// 	"fa53015f4967aeb13f5e8f3facd3832c",
		// 	"2dd18ca3c56d86a608edc852e887c4c6",
		// 	"5a0958689482d5e41e2f77ce2ea2b5a9",
		// 	"5f0de6f30413285f0b7d677a400ef8ad",
		// 	"7293eaffff30dfc57c9070472199c0cb",
		// 	"be8332c53324ee35470983f2368e8ddd",
		// 	"c41d9b9f81a522c932a7eec76ae7f05a",
		// 	"c4aaf9e19c23b5790e5da790cf2308e5",
		// 	"5e74f6fe2c2e522e81e041e3d09ea7aa",
		// 	"bcd2fda472eb27a41773ef81d65a2999",
		// 	"d2f24c095543d62ac5de56ccd8a1c4ee",
		// 	"e7372375d9c37c40402e1af6116ef038"}

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
			if brick.Info.Path == "" {
				println(brick.Info.Id + " No Path -> Fixing")

				// Access device
				device, err := NewDeviceEntryFromId(tx, brick.Info.DeviceId)
				if err != nil {
					println(err.Error())
					return err
				}

				v := &VolumeEntry{}
				v.removeBrickFromDb(tx, brick)

				if err != nil {
					println(err.Error())
				}

				println(brick.Info.Id + " Deleted on device: " + device.Info.Id)
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
