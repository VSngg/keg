package mark

import "fmt"

type TagTooLong struct {
	Tag string
}

func (e TagTooLong) Error() string {
	return fmt.Sprintf("tag too long: %v", e.Tag)
}

type Expected struct {
	This string
}

func (e Expected) Error() string {
	return fmt.Sprintf("expected: %v", e.This)
}
