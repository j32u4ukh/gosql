package cntr

import (
	"bytes"
	"fmt"
)

func Slice2dToString[T Element](slice2d [][]T) string {
	var buffer bytes.Buffer
	length := len(slice2d)
	buffer.WriteString("{")
	if length > 0 {
		buffer.WriteString(SliceToString(slice2d[0]))
		for i := 1; i < length; i++ {
			buffer.WriteString(fmt.Sprintf(", %s", SliceToString(slice2d[0])))
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}

func SliceToString[T Element](slice []T) string {
	var buffer bytes.Buffer
	length := len(slice)
	buffer.WriteString("{")
	if length > 0 {
		buffer.WriteString(fmt.Sprintf("%v", slice[0]))
		for i := 1; i < length; i++ {
			buffer.WriteString(fmt.Sprintf(", %v", slice[i]))
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}
