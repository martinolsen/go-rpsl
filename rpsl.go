package rpsl

import (
	"strings"
	"fmt"
)

type Object struct {
	Class  string
	Values map[string][]string
}

func (o *Object) Get(key string) string {
	if vs, ok := o.Values[strings.ToLower(key)]; ok && len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (o *Object) String() string {
	var s string

	s = fmt.Sprintf("%s:\t%s\n", o.Class, o.Get(o.Class))

	for k, vs := range o.Values {
		if k == o.Class {
			continue
		}

		for _, v := range vs {
			s = fmt.Sprintf("%s%s:\t%s\n", s, k, v)
		}
	}

	return fmt.Sprintf("%s\n", s)
}
