package article

type HandleType int

const (
	DOI HandleType = iota
	ISBN
)

type Handle struct {
	Value string
	Type  HandleType
}

func (h Handle) String() string {
	return h.Value
}
