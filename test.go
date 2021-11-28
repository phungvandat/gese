package gese

type tE struct {
	F string
	V interface{}
}

type tD struct {
	E *tE
	M map[string]interface{}
}

type tC struct {
	D tD
	V interface{}
}

type tB struct {
	C *tC
	V interface{}
}

type tA struct {
	B tB
	V interface{}
}
