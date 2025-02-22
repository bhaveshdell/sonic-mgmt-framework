/*
Copyright 2019 Broadcom. All rights reserved.
The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
*/

/*
Package db implements a wrapper over the go-redis/redis.
*/
package db

import (
	// "fmt"
	// "strconv"

	//	"reflect"
	// "errors"
	// "strings"

	// "github.com/go-redis/redis"
	"github.com/golang/glog"
	// "cvl"
	"translib/tlerr"
)

func init() {
}




func (d *DB) GetMap(ts *TableSpec, mapKey string) (string, error) {

	if glog.V(3) {
		glog.Info("GetMap: Begin: ", "ts: ", ts, " mapKey: ", mapKey)
	}

	v, e := d.client.HGet(ts.Name, mapKey).Result()

	if glog.V(3) {
		glog.Info("GetMap: End: ", "v: ", v, " e: ", e)
	}

	return v, e
}

func (d *DB) GetMapAll(ts *TableSpec) (Value, error) {

	if glog.V(3) {
		glog.Info("GetMapAll: Begin: ", "ts: ", ts)
	}

	var value Value

	v, e := d.client.HGetAll(ts.Name).Result()

	if len(v) != 0 {
		value = Value{Field: v}
	} else {
		if glog.V(1) {
			glog.Info("GetMapAll: HGetAll(): empty map")
		}
		e = tlerr.TranslibRedisClientEntryNotExist { Entry: ts.Name }
	}

	if glog.V(3) {
		glog.Info("GetMapAll: End: ", "value: ", value, " e: ", e)
	}

	return value, e
}

// For Testing only. Do Not Use!!! ==============================
// There is no transaction support on these.
func (d *DB) SetMap(ts *TableSpec, mapKey string, mapValue string) error {

	if glog.V(3) {
		glog.Info("SetMap: Begin: ", "ts: ", ts, " ", mapKey,
			":", mapValue)
	}

	b, e := d.client.HSet(ts.Name, mapKey, mapValue).Result()

	if glog.V(3) {
		glog.Info("GetMap: End: ", "b: ", b, " e: ", e)
	}

	return e
}
// For Testing only. Do Not Use!!! ==============================

// For Testing only. Do Not Use!!!
// There is no transaction support on these.
func (d *DB) DeleteMapAll(ts *TableSpec) error {

	if glog.V(3) {
		glog.Info("DeleteMapAll: Begin: ", "ts: ", ts)
	}

	e := d.client.Del(ts.Name).Err()

	if glog.V(3) {
		glog.Info("DeleteMapAll: End: ", " e: ", e)
	}

	return e
}
// For Testing only. Do Not Use!!! ==============================


