package saver

type Null struct {
	Dir string
}

func (s Null) Save(path string, data []byte) error {
	return nil
}

func NewNull() *Null {
	return &Null{}
}

func NewNullFromMap(m map[string]interface{}) (*Null, error) {
	return NewNull(), nil
}
