package file_db

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"sync"
)

type FileDB struct {
	path string
	list sync.Map
	t    reflect.Type
	m    sync.Mutex
}

var (
	ErrInvalidType = errors.New("invalid type")
)

func Open(fileName string, item interface{}) *FileDB {
	// Ensure item is jsonable (not definitive)
	_, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}
	db := FileDB{
		path: fileName,
		t:    reflect.TypeOf(item),
	}
	if err := db.load(); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Data file %q not found, starting with empty database", db.path)
		} else {
			log.Printf("Error loading data from %q: %s", db.path, err)
		}
	}
	return &db
}

// load data from filesystem
func (db *FileDB) load() error {
	db.m.Lock()
	defer db.m.Unlock()
	f, err := os.Open(db.path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	data := make(map[int]json.RawMessage)
	if err = json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}
	for i := range data {
		item := reflect.New(db.t).Interface()
		if err := json.Unmarshal(data[i], item); err != nil {
			return err
		}
		db.list.Store(i, reflect.ValueOf(item).Elem().Interface())
	}
	return nil
}

// store data to filesystem
func (db *FileDB) store() error {
	db.m.Lock()
	defer db.m.Unlock()
	f, err := os.OpenFile(db.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	out := make(map[int]interface{})
	db.list.Range(func(key, value interface{}) bool {
		out[key.(int)] = value
		return true
	})
	e := json.NewEncoder(f)
	e.SetIndent("", "\t")
	return e.Encode(out)
}

// Set adds or replaces `item` with the given `id`
func (db *FileDB) Set(id int, item interface{}) error {
	if reflect.TypeOf(item) != db.t {
		return ErrInvalidType
	}
	db.list.Store(id, item)
	return db.store()
}

// Append adds another `item` and returns its `id` or error
func (db *FileDB) Append(item interface{}) (int, error) {
	if reflect.TypeOf(item) != db.t {
		return 0, ErrInvalidType
	}
	var id int
	db.list.Range(func(key, value interface{}) bool {
		if kid := key.(int); kid >= id {
			id = kid
		}
		return true
	})
	return id + 1, db.Set(id+1, item)
}

// Get returns item by id
func (db *FileDB) Get(id int) interface{} {
	if item, ok := db.list.Load(id); ok {
		return item
	}
	return nil
}

func (db *FileDB) Delete(id int) error {
	db.list.Delete(id)
	return db.store()
}

// All returns slice of all items
func (db *FileDB) All() interface{} {
	items := reflect.MakeSlice(reflect.SliceOf(db.t), 0, 0)
	db.list.Range(func(key, value interface{}) bool {
		items = reflect.Append(items, reflect.ValueOf(value))
		return true
	})
	return items.Interface()
}

// Map applies filter function and returns resulting slice of items
func (db *FileDB) Map(filter func(interface{}) bool) interface{} {
	items := reflect.MakeSlice(reflect.SliceOf(db.t), 0, 0)
	db.list.Range(func(key, value interface{}) bool {
		if filter(value) {
			items = reflect.Append(items, reflect.ValueOf(value))
		}
		return true
	})
	return items.Interface()
}
