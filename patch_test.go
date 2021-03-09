package diff

import "testing"

type bar struct {
	Data string `diff:"data"`
}
type container struct {
	Name        string
	Number      int64
	StringSlice []string
	MiscSlice   []interface{}
	Slice       []bar
	Struct      bar
	Map         map[string]bar
	SimpleMap   map[string]string
	NestedMap   map[string]interface{}
}

func initValue() container {
	return container{
		Name:        "foo",
		Number:      1,
		Slice:       []bar{bar{"bar"}, bar{"bar2"}},
		StringSlice: []string{"bar", "bar2"},
		MiscSlice:   []interface{}{},
		Struct:      bar{"bar"},
		Map:         map[string]bar{"bar": bar{"bar"}, "bar2": bar{"bar2"}},
		NestedMap:   map[string]interface{}{"bar": bar{"bar"}, "bar2": bar{"bar2"}, "bar3": map[string]string{"foo": "foo2"}},
		SimpleMap:   map[string]string{"bar": "bar", "bar2": "bar2"},
	}
}

func TestPatch(t *testing.T) {
	t.Run("UpdateSimple", func(t *testing.T) {
		source := initValue()
		update := initValue()
		update.Name = "foo Updated"
		update.Number = 2
		changelog, err := Diff(source, update)
		if err != nil {
			t.Log("Failed to calc diff")
			t.Log(err)
			t.FailNow()
		}
		patchLog := Patch(changelog, &source)
		if len(patchLog) != 2 {
			t.Log("Incorrect number of patches - should be 2")
			t.Log(changelog)
			t.Log(patchLog)
			t.FailNow()
		}
		if source.Name != update.Name {
			t.Log("Name not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.FailNow()
		}
		if source.Number != update.Number {
			t.Log("Number not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.Fail()
		}
	})
	t.Run("UpdateStringSlice", func(t *testing.T) {
		source := initValue()
		update := initValue()
		update.StringSlice[0] = "bar update"
		changelog, err := Diff(source, update)
		if err != nil {
			t.Log("Failed to calc diff")
			t.Log(err)
			t.FailNow()
		}
		patchLog := Patch(changelog, &source)
		if len(patchLog) != 1 {
			t.Log("Incorrect number of patches - should be 1")
			t.Log(changelog)
			t.Log(patchLog)
			t.FailNow()
		}
		if source.StringSlice[0] != update.StringSlice[0] {
			t.Log("Slice not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.Fail()
		}
	})
	t.Run("UpdateSlice", func(t *testing.T) {
		source := initValue()
		update := initValue()
		update.Slice[0].Data = "bar update"
		changelog, err := Diff(source, update)
		if err != nil {
			t.Log("Failed to calc diff")
			t.Log(err)
			t.FailNow()
		}
		patchLog := Patch(changelog, &source)
		if len(patchLog) != 1 {
			t.Log("Incorrect number of patches - should be 1")
			t.Log(changelog)
			t.Log(patchLog)
			t.FailNow()
		}
		if source.Slice[0].Data != update.Slice[0].Data {
			t.Log("Slice not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.Fail()
		}
	})
	t.Run("UpdateStruct", func(t *testing.T) {
		source := initValue()
		update := initValue()
		update.Struct.Data = "bar update"
		changelog, err := Diff(source, update)
		if err != nil {
			t.Log("Failed to calc diff")
			t.Log(err)
			t.FailNow()
		}
		patchLog := Patch(changelog, &source)
		if len(patchLog) != 1 {
			t.Log("Incorrect number of patches - should be 1")
			t.Log(changelog)
			t.Log(patchLog)
			t.FailNow()
		}
		if source.Struct.Data != update.Struct.Data {
			t.Log("Struct not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.Fail()
		}
	})
	t.Run("UpdateMap", func(t *testing.T) {
		source := initValue()
		update := initValue()
		b := update.Map["bar"]
		b.Data = "bar Update"
		update.Map["bar"] = b
		changelog, err := Diff(source, update)
		if err != nil {
			t.Log("Failed to calc diff")
			t.Log(err)
			t.FailNow()
		}
		patchLog := Patch(changelog, &source)
		if len(patchLog) != 1 {
			t.Log("Incorrect number of patches - should be 1")
			t.Log(changelog)
			t.Log(patchLog)
			t.FailNow()
		}
		if source.Map["bar"].Data != update.Map["bar"].Data {
			t.Log("Map not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.Fail()
		}
	})
	t.Run("UpdateSimpleMap", func(t *testing.T) {
		source := initValue()
		update := initValue()
		update.SimpleMap["bar"] = "bar32"
		changelog, err := Diff(source, update)
		if err != nil {
			t.Log("Failed to calc diff")
			t.Log(err)
			t.FailNow()
		}
		patchLog := Patch(changelog, &source)
		if len(patchLog) != 1 {
			t.Log("Incorrect number of patches - should be 1")
			t.Log(changelog)
			t.Log(patchLog)
			t.FailNow()
		}
		if source.SimpleMap["bar"] != update.SimpleMap["bar"] {
			t.Log("Map not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.Fail()
		}
	})

	t.Run("UpdateNestedMap", func(t *testing.T) {
		source := initValue()
		update := initValue()
		update.NestedMap["bar"] = bar{"bar32"}
		update.NestedMap["bar3"] = map[string]string{"foo2": "foo2", "foo3": "f003"}
		delete(update.NestedMap, "bar2")
		changelog, err := Diff(source, update)
		if err != nil {
			t.Log("Failed to calc diff")
			t.Log(err)
			t.FailNow()
		}
		patchLog := Patch(changelog, &source)
		if len(patchLog) != 5 {
			t.Log("Incorrect number of patches - should be 5")
			t.Log(changelog)
			t.Log(patchLog)
			t.FailNow()
		}
		transformed, _ := update.NestedMap["bar3"].(map[string]string)
		if transformed["foo2"] != "foo2" {
			t.Log("Map not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.Fail()
		}
		if transformed["foo3"] != "f003" {
			t.Log("Map not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.Fail()
		}
		if transformed["foo"] != "" {
			t.Log("Map not updated")
			t.Log(changelog)
			t.Log(patchLog)
			t.Fail()
		}
	})
}
