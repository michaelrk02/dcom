package dcom

type Structure interface {
	Marshal(m Marshaler) error
	Unmarshal(u Unmarshaler) error
}
