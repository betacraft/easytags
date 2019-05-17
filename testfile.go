package main

type TestStruct struct {
	Field1     int `json:"-"`
	TestField2 string
	ExistingTag string `custom:"" json:"etag"`
	Embed
}

type Embed struct {
}
