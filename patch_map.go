package diff

import (
	"reflect"

	"github.com/vmihailenco/msgpack"
)

//renderMap - handle map rendering for patch
func (d *Differ) renderMap(c *ChangeValue) (m, k, v *reflect.Value) {

	//we must tease out the type of the key, we use the msgpack from diff to recreate the key
	kt := c.target.Type().Key()
	field := reflect.New(kt)
	if err := msgpack.Unmarshal([]byte(c.change.Path[c.pos]), field.Interface()); err != nil {
		c.SetFlag(FlagIgnored)
		c.AddError(NewError("Unable to unmarshal path element to target type for key in map", err))
		return
	}
	c.key = field.Elem()
	if c.target.IsNil() && c.target.IsValid() {
		c.target.Set(reflect.MakeMap(c.target.Type()))
	}
	x := c.target.MapIndex(c.key)

	if !x.IsValid() && c.change.Type != DELETE && !c.HasFlag(OptionNoCreate) {
		x = c.NewElement()
	}
	if x.IsValid() { //Map elements come out as read only so we must convert
		nv := reflect.New(x.Type()).Elem()
		nv.Set(x)
		x = nv
	}

	if x.IsValid() && !reflect.DeepEqual(c.change.From, x.Interface()) &&
		c.HasFlag(OptionOmitUnequal) {
		c.SetFlag(FlagIgnored)
		c.AddError(NewError("target change doesn't match original"))
		return
	}
	mp := *c.target //these may change out from underneath us as we recurse
	key := c.key    //so we make copies and pass back pointers to them
	c.swap(&x)
	return &mp, &key, &x

}

//deleteMapEntry - deletes are special, they are handled differently based on options
//            container type etc. We have to have special handling for each
//            type. Set values are more generic even if they must be instanced
func (d *Differ) deleteMapEntry(c *ChangeValue, m, k, v *reflect.Value) {
	if m != nil && m.CanSet() && v.IsValid() {
		for _, x := range m.MapKeys() {
			if !m.MapIndex(x).IsZero() {
				m.SetMapIndex(*k, *v)
				return
			}
		} //if all the fields are zero, remove from map
		m.SetMapIndex(*k, reflect.Value{})
		c.SetFlag(FlagDeleted)
	}
}
