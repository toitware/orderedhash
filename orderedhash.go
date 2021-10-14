// Copyright (C) 2021 Toitware ApS.  All rights reserved.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

// Insertion-ordered hash set and hash map with an implementation inspired by
// the non-hateful maps described in
// https://blog.toit.io/hash-maps-that-dont-hate-you-1a96150b492a
// The collections preserve insertion order and allow custom equality
// functions.  The collections are not thread safe and require external
// synchronization.
package orderedhash

// To use OrderedMap or OrderedSet you must first define an EqualityRelation
// for your set elements or map keys.  The two methods in this interface must
// be coherent in that if two objects are equal, they must have the same hash
// code.  Hash collisions are allowed but will reduce efficiency if they are
// frequent.  The hash code returned should not change on subsequent calls,
// and the return value from Equals should also be stable.
type EqualityRelation interface {
	Equals(object_in_collection interface{}, new_object interface{}) bool
	Hash(object interface{}) int
}

type orderedHash struct {
	hashToIndex map[int][]int
	backing     []interface{}
	equality    EqualityRelation
	len         int
}

// An insertion-ordered hash set with customizable equality function.
type OrderedSet struct {
	orderedHash
}

// An insertion-ordered hash map with customizable equality function.
type OrderedMap struct {
	orderedHash
	valueBacking []interface{}
}

// Create an empty insertion-ordered set with customized equality function.
func NewSet(relation EqualityRelation) *OrderedSet {
	new := OrderedSet{
		orderedHash{
			hashToIndex: make(map[int][]int),
			backing:     []interface{}{},
			equality:    relation,
		},
	}
	return &new
}

// Create an empty insertion-ordered map with customized equality function.
func NewMap(relation EqualityRelation) *OrderedMap {
	new := OrderedMap{
		orderedHash: orderedHash{
			hashToIndex: make(map[int][]int),
			backing:     []interface{}{},
			equality:    relation,
		},
		valueBacking: []interface{}{},
	}
	return &new
}

// Get the number of elements in the set or the number of key-value pairs in
// the map.
func (o *orderedHash) Len() int {
	return o.len
}

// Add the element to the set if it does not already contain an equal element.
func (o *OrderedSet) Add(element interface{}) {
	if element == nil {
		panic("Can't add nil to a set")
	}
	hash := o.equality.Hash(element)
	indices := o.hashToIndex[hash]
	if indices == nil {
		// No entries with this hash code.  Create a new one.
		index := len(o.backing)
		o.backing = append(o.backing, element)
		o.hashToIndex[hash] = []int{index}
		o.len++
		return
	}
	deleted_space := -1
	for i, index := range indices {
		candidate := o.backing[index]
		if candidate != nil {
			if o.equality.Equals(candidate, element) {
				// Already present in set.
				return
			}
		} else if deleted_space == -1 {
			deleted_space = i
		}
	}
	// Not found.  Add an index to the entry for this hash code.
	index := len(o.backing)
	o.backing = append(o.backing, element)
	if deleted_space == -1 {
		o.hashToIndex[hash] = append(indices, index)
	} else {
		indices[deleted_space] = index
	}
	o.len++
}

// Add the key and value to the map.  If the map already contains an
// equal key then the value is overwritten, but the key is unchanged.
func (o *OrderedMap) Set(key interface{}, value interface{}) {
	if key == nil || value == nil {
		panic("Can't add nil to a map")
	}
	hash := o.equality.Hash(key)
	indices := o.hashToIndex[hash]
	if indices == nil {
		// No entries with this hash code.  Create a new one.
		index := len(o.backing)
		o.backing = append(o.backing, key)
		o.valueBacking = append(o.valueBacking, value)
		o.hashToIndex[hash] = []int{index}
		o.len++
		return
	}
	deleted_space := -1
	for i, index := range indices {
		candidate := o.backing[index]
		if candidate != nil {
			if o.equality.Equals(candidate, key) {
				// Already present in map.  Overwrite value.
				o.valueBacking[index] = value
				return
			}
		} else if deleted_space == -1 {
			deleted_space = i
		}
	}
	// Not found.  Add an index to the entry for this hash code.
	index := len(o.backing)
	o.backing = append(o.backing, key)
	o.valueBacking = append(o.valueBacking, value)
	if deleted_space == -1 {
		o.hashToIndex[hash] = append(indices, index)
	} else {
		indices[deleted_space] = index
	}
	o.len++
}

// Whether the set contains an equal element, or whether the map
// contains an equal key.
func (o *orderedHash) Contains(element interface{}) bool {
	if element == nil {
		panic("Can't use nil in a set or use nil as a map key")
	}
	hash := o.equality.Hash(element)
	indices := o.hashToIndex[hash]
	if indices == nil {
		return false
	}
	for _, index := range indices {
		candidate := o.backing[index]
		if candidate != nil {
			if o.equality.Equals(candidate, element) {
				// Already present in set.
				return true
			}
		}
	}
	return false
}

// Get an equal element that is already in the set or an equal
// key that is already in the map.
func (o *orderedHash) GetKey(element interface{}) interface{} {
	if element == nil {
		panic("Can't use nil as a set element or as a map key")
	}
	hash := o.equality.Hash(element)
	indices := o.hashToIndex[hash]
	if indices == nil {
		return nil
	}
	for _, index := range indices {
		candidate := o.backing[index]
		if candidate != nil {
			if o.equality.Equals(candidate, element) {
				// Found.
				return candidate
			}
		}
	}
	return nil
}

// Get the value corresponding to a key.  Returns nil if an equal
// key is not found in the map.
func (o *OrderedMap) Get(key interface{}) interface{} {
	if key == nil {
		panic("Can't use nil as a map key")
	}
	hash := o.equality.Hash(key)
	indices := o.hashToIndex[hash]
	if indices == nil {
		return nil
	}
	for _, index := range indices {
		candidate := o.backing[index]
		if candidate != nil {
			if o.equality.Equals(candidate, key) {
				// Found.
				return o.valueBacking[index]
			}
		}
	}
	return nil
}

// Remove an equal element from a set.
// If an element is removed and then later re-added, its iteration order
// is moved to the end.
func (o *OrderedSet) Remove(element interface{}) {
	if element == nil {
		panic("Can't use nil as a set element")
	}
	hash := o.equality.Hash(element)
	indices := o.hashToIndex[hash]
	if indices == nil {
		// No entries with this hash code.
		return
	}
	for _, index := range indices {
		candidate := o.backing[index]
		if candidate != nil {
			if o.equality.Equals(candidate, element) {
				// Found.  We are using nil as tombstone.
				o.backing[index] = nil
				o.len--
				// If there was only one entry in the map with this hash code
				// we might as well remove it.  TODO: We could also remove a
				// single entry when there are hash collisions.
				if len(indices) == 1 {
					delete(o.hashToIndex, hash)
				}
			}
		}
	}
}

// Remove an equal key and its associated value from a map.
// If a key is removed and then later re-added, its iteration order
// is moved to the end.
func (o *OrderedMap) Remove(key interface{}) {
	if key == nil {
		panic("Can't use nil as a map key")
	}
	hash := o.equality.Hash(key)
	indices := o.hashToIndex[hash]
	if indices == nil {
		// No entries with this hash code.
		return
	}
	for _, index := range indices {
		candidate := o.backing[index]
		if candidate != nil {
			if o.equality.Equals(candidate, key) {
				// Found.  We are using nil as tombstone.
				o.backing[index] = nil
				o.valueBacking[index] = nil
				o.len--
				// If there was only one entry in the map with this hash code
				// we might as well remove it.  TODO: We could also remove a
				// single entry when there are hash collisions.
				if len(indices) == 1 {
					delete(o.hashToIndex, hash)
				}
			}
		}
	}
}

// If the set already contains an equal element, replace it with the given one.
// If the map already contains an equal key, replace it with the given one.
// The new element or key inherites the insertion order of the element or key
// it replaces.
func (o *orderedHash) ReplaceWith(element interface{}) {
	if element == nil {
		panic("Can't use nil as a set element or a map key")
	}
	hash := o.equality.Hash(element)
	indices := o.hashToIndex[hash]
	if indices == nil {
		// Does not contain element/key.
		return
	}
	for _, index := range indices {
		candidate := o.backing[index]
		if candidate != nil {
			if o.equality.Equals(candidate, element) {
				// Already present in set/map.  Replace with given element/key.
				o.backing[index] = element
				return
			}
		}
	}
}

// Iterable slice of the elements in the set or the keys in a map.
// Iteration is in insertion order.  If the set or map is modified
// during iteration the changes may or may not be reflected in this
// slice.
func (o *orderedHash) Entries() []interface{} {
	if o.len == len(o.backing) {
		return o.backing
	}
	result := make([]interface{}, o.len)
	i := 0
	for _, entry := range o.backing {
		if entry != nil {
			result[i] = entry
			i = i + 1
		}
	}
	return result
}

// Iterable slice of the values in a map.
func (o *OrderedMap) Values() []interface{} {
	if o.len == len(o.valueBacking) {
		return o.valueBacking
	}
	result := make([]interface{}, o.len)
	i := 0
	for _, entry := range o.valueBacking {
		if entry != nil {
			result[i] = entry
			i = i + 1
		}
	}
	return result
}

// A reasonable equality function for strings.
type StringEquality struct {
}

func (_ StringEquality) Equals(a interface{}, b interface{}) bool {
	a_string := a.(string)
	b_string := b.(string)
	return a_string == b_string
}

func (_ StringEquality) Hash(a interface{}) int {
	h := 0
	for _, char := range a.(string) {
		h *= 11
		h = h + int(char)
	}
	return h
}

func NewStringSet() *OrderedSet {
	new := NewSet(StringEquality{})
	return new
}

func NewStringMap() *OrderedMap {
	new := NewMap(StringEquality{})
	return new
}
