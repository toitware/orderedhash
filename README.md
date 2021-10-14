# orderedhash

## Overview

Insertion-ordered hash set and hash map with an implementation inspired by
the non-hateful maps described in
https://blog.toit.io/hash-maps-that-dont-hate-you-1a96150b492a
The collections preserve insertion order and allow custom equality
functions.  The collections are not thread safe and require external
synchronization.

## Usage

```
// Implements orderedhash.EqualityRelation
type MyClassEquality struct {
}

func (_ MyClassEquality) Equals(a interface{}, b interface{}) bool {
	a_object := a.(MyClass)
	b_object := b.(MyClass)
	return a_object.Equals(b_object)
}

func (_ MyClassEquality) Hash(a interface{}) int {
	object := a.(MyClass)
    return object.Hash()
}

main() {
    s := orderedhash.NewSet(MyClassEquality{})
    s.Add(NewMyClass(42))
    s.Add(NewMyClass(103))
    for _, element := range s.Entries() {
        element.Foo()
    }

    m := orderedhash.NewMap(MyClassEquality{})
    m.Set(NewMyClass(42), "FortyTwo")
    m.Set(NewMyClass(103), "OneHundredAndThree")
    for _, key := range m.Entries() {
        Bar(key, m.Get(key))
    }
}
```
