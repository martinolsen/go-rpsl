package rpsl

import (
	"io"
)

func Lookup(reader *Reader, query string) []*Object {
	var objects []*Object

	for {
		object, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		var ok bool
		for _, vs := range object.Values {
			for _, v := range vs {
				if v == query {
					ok = true
				}
			}
		}
		if ok {
			objects = append(objects, object)
		}
	}

	return objects
}
