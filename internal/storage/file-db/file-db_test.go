package file_db_test

import (
	"os"
	"sort"
	"testing"

	fdb "github.com/Sergii-Kirichok/pr/internal/storage/file-db"
	"github.com/stretchr/testify/assert"
)

func TestFileDB(t *testing.T) {
	type Test struct {
		ID   int
		Name string
		Flag bool
	}
	item := Test{
		ID:   1,
		Name: "name",
		Flag: true,
	}
	assert.Panics(t, func() {
		fdb.Open("test.json", func() {})
	})
	db := fdb.Open("test.json", &Test{})
	if assert.NoError(t, db.Set(item.ID, &item)) {
		if it, ok := db.Get(1).(*Test); ok {
			assert.Equal(t, &item, it)
		}
	}
	if id, err := db.Append(&Test{
		ID:   6,
		Name: "another",
		Flag: false,
	}); assert.NoError(t, err) {
		assert.Equal(t, 2, id)
	}
	expected := []*Test{
		{
			ID:   1,
			Name: "name",
			Flag: true,
		},
		{
			ID:   6,
			Name: "another",
			Flag: false,
		},
	}
	all := db.All()
	sort.SliceStable(all, func(i, j int) bool {
		return all.([]*Test)[i].ID < all.([]*Test)[j].ID
	})
	assert.Equal(t, expected, all)
	names := db.Map(func(i interface{}) bool {
		return i.(*Test).Name == "name"
	})
	assert.Equal(t, []*Test{
		{
			ID:   1,
			Name: "name",
			Flag: true,
		},
	}, names)
	assert.Empty(t, db.Get(99))
	assert.Error(t, db.Set(7, true))
	_, err := db.Append(5)
	assert.Error(t, err)
	db2 := fdb.Open("test.json", &Test{})
	items := db2.All()
	sort.SliceStable(items, func(i, j int) bool {
		return items.([]*Test)[i].ID < items.([]*Test)[j].ID
	})
	assert.Equal(t, expected, items)
	assert.NoError(t, db2.Delete(1))
	items = db2.All()
	sort.SliceStable(items, func(i, j int) bool {
		return items.([]*Test)[i].ID < items.([]*Test)[j].ID
	})
	assert.Equal(t, []*Test{
		{
			ID:   6,
			Name: "another",
			Flag: false,
		},
	}, items)
	_ = os.Remove("test.json")
}
