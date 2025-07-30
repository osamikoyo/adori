package models

type YadoriFile struct {
	Path              string
	Content           string
	Include           bool
	ExpireAtTimestamp int64
}