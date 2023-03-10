package util

type (
	Set interface {
		Has(interface{}) bool
		Add(interface{})
		Remove(interface{})
		Len() int
		Values() interface{}
	}
	StringSet map[string]struct{}
	IntSet    map[int]struct{}
)

func (ss StringSet) Has(v interface{}) bool {
	_, ok := ss[v.(string)]
	return ok
}

func (ss StringSet) Add(v interface{}) {
	ss[v.(string)] = Empty()
}

func (ss StringSet) Remove(v interface{}) {
	delete(ss, v.(string))
}

func (ss StringSet) Len() int {
	return len(ss)
}

func (ss StringSet) Values() interface{} {
	st := make([]string, 0, len(ss))
	for i := range ss {
		st = append(st, i)
	}
	return st
}

func (ss IntSet) Has(v interface{}) bool {
	_, ok := ss[v.(int)]
	return ok
}

func (ss IntSet) Add(v interface{}) {
	ss[v.(int)] = Empty()
}

func (ss IntSet) Remove(v interface{}) {
	delete(ss, v.(int))
}

func (ss IntSet) Len() int {
	return len(ss)
}

func (ss IntSet) Values() interface{} {
	st := make([]int, 0, len(ss))
	for i := range ss {
		st = append(st, i)
	}
	return st
}

func Empty() struct{} {
	return struct{}{}
}

func NewStringSet() Set {
	return StringSet(make(map[string]struct{}))
}

func NewIntSet() Set {
	return IntSet(make(map[int]struct{}))
}

func SetForStringSlice(s []string) Set {
	res := NewStringSet()
	for _, v := range s {
		res.Add(v)
	}
	return res
}

func SetForIntSlice(s []int) Set {
	res := NewIntSet()
	for _, v := range s {
		res.Add(v)
	}
	return res
}
