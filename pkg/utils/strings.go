package utils

import "strings"

func ParseFilter(s string) []string {
	return strings.Split(s, ",")
}

func StringSliceToInterfaceSlice(ss []string) []interface{} {
	is := make([]interface{}, len(ss))
	for i, v := range ss {
		is[i] = v
	}
	return is
}
