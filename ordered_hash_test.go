// Copyright (C) 2021 Toitware ApS.  All rights reserved.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package orderedhash

import "testing"

// An unreasonable equality function for strings, where all
// strings have the same hash.  This stress tests the collection.
type AllSameHashEquality struct {
}

func (_ AllSameHashEquality) Equals(a interface{}, b interface{}) bool {
	a_string := a.(string)
	b_string := b.(string)
	return a_string == b_string
}

func (_ AllSameHashEquality) Hash(a interface{}) int {
	return 0
}

func newStressSet() *OrderedSet {
	new := NewSet(AllSameHashEquality{})
	return new
}

func newStressMap() *OrderedMap {
	new := NewMap(AllSameHashEquality{})
	return new
}

func checkSet(t *testing.T, s *OrderedSet, expect []string) {
	if s.Contains("") {
		t.Fatal("We did not expect the empty string")
	}
	if s.Contains("foo") {
		t.Fatal("We did not expect the foo string")
	}
	if s.Len() != len(expect) {
		t.Fatal("Len")
	}
	for _, str := range expect {
		if !s.Contains(str) {
			t.Fatal("Doesn't contain")
		}
		if s.GetKey(str) != str {
			t.Fatal("Didn't get element back")
		}
	}
	count := 0
	for i, str := range s.Entries() {
		if str != expect[i] {
			t.Fatal(str)
		}
		count++
	}
	if count != len(expect) {
		t.Fatal("Iterations")
	}
}

func checkMap(t *testing.T, m *OrderedMap, expect []string) {
	if m.Contains("") {
		t.Fatal("We did not expect the empty string")
	}
	if m.Contains("foo") {
		t.Fatal("We did not expect the foo string")
	}
	if m.Len() != len(expect)/2 {
		t.Fatal("Len")
	}
	for i := 0; i < len(expect); i += 2 {
		key := expect[i]
		if !m.Contains(key) {
			t.Fatal("Doesn't contain")
		}
		if m.GetKey(key) != key {
			t.Fatal("Didn't get key back")
		}
		if m.Get(key) != expect[i+1] {
			t.Fatal("Value not right")
		}
	}
	count := 0
	for i, str := range m.Entries() {
		if str != expect[i*2] {
			t.Fatal(str)
		}
		count++
	}
	if count != len(expect)/2 {
		t.Fatal("Iterations")
	}
	count = 0
	for i, str := range m.Values() {
		if str != expect[i*2+1] {
			t.Fatal(str)
		}
		count++
	}
	if count != len(expect)/2 {
		t.Fatal("Iterations")
	}
}

func TestAll(t *testing.T) {
	s := NewStringSet()
	testSet(t, s)
	s2 := newStressSet()
	testSet(t, s2)
	m := NewStringMap()
	testMap(t, m)
	m2 := newStressMap()
	testMap(t, m2)
}

func testSet(t *testing.T, s *OrderedSet) {
	checkSet(t, s, []string{})
	s.Add("Foo")
	checkSet(t, s, []string{"Foo"})
	s.Add("Bar")
	checkSet(t, s, []string{"Foo", "Bar"})
	s.Add("Bar")
	checkSet(t, s, []string{"Foo", "Bar"})
	s.Add("Foo")
	checkSet(t, s, []string{"Foo", "Bar"})
	s.Remove("Fizz")
	checkSet(t, s, []string{"Foo", "Bar"})
	s.Add("Fizz")
	checkSet(t, s, []string{"Foo", "Bar", "Fizz"})
	s.Remove("Fizz")
	checkSet(t, s, []string{"Foo", "Bar"})
	s.Add("Now")
	s.Add("Is")
	s.Add("The")
	s.Add("Time")
	checkSet(t, s, []string{"Foo", "Bar", "Now", "Is", "The", "Time"})
	s.ReplaceWith("Time")
	checkSet(t, s, []string{"Foo", "Bar", "Now", "Is", "The", "Time"})
	s.Remove("Now")
	s.Remove("Is")
	s.Remove("The")
	s.Remove("Time")
	s.Remove("Foo")
	s.Remove("Foo")
	s.Remove("Bar")
	checkSet(t, s, []string{})
	s.Remove("Time")
	checkSet(t, s, []string{})
}

func testMap(t *testing.T, m *OrderedMap) {
	checkMap(t, m, []string{})
	m.Set("Foo", "Bar")
	checkMap(t, m, []string{"Foo", "Bar"})
	m.Set("Fizz", "Buzz")
	checkMap(t, m, []string{"Foo", "Bar", "Fizz", "Buzz"})
	m.Remove("Foo")
	checkMap(t, m, []string{"Fizz", "Buzz"})
	m.ReplaceWith("Foo")
	checkMap(t, m, []string{"Fizz", "Buzz"})
	m.ReplaceWith("Fizz")
	checkMap(t, m, []string{"Fizz", "Buzz"})
	m.Set("Foo", "Bar")
	checkMap(t, m, []string{"Fizz", "Buzz", "Foo", "Bar"})
	m.Set("Foo", "Bar2")
	checkMap(t, m, []string{"Fizz", "Buzz", "Foo", "Bar2"})
}
