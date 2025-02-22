package dcom

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/google/uuid"
)

type defaultMarshaler struct {
	wr *bufio.Writer
}

func NewDefaultMarshaler(w io.Writer) Marshaler {
	return &defaultMarshaler{
		wr: bufio.NewWriter(w),
	}
}

func (m *defaultMarshaler) WriteUUID(v uuid.UUID) error {
	_, err := m.wr.Write(v[:])
	if err != nil {
		return errors.Join(ErrMarshalerWrite, err)
	}

	err = m.wr.Flush()
	if err != nil {
		return errors.Join(ErrMarshalerWrite, err)
	}

	return nil
}

func (m *defaultMarshaler) WriteBool(v bool) error {
	n := byte(0)
	if v {
		n = byte(1)
	}

	err := m.wr.WriteByte(n)
	if err != nil {
		return errors.Join(ErrMarshalerWrite, err)
	}

	err = m.wr.Flush()
	if err != nil {
		return errors.Join(ErrMarshalerWrite, err)
	}

	return nil
}

func (m *defaultMarshaler) WriteBoolOptional(v *bool) error {
	return m.writeOptional(
		v != nil,
		func() error {
			return m.WriteBool(*v)
		},
	)
}

func (m *defaultMarshaler) WriteBoolArray(v []bool) error {
	return m.writeArray(
		len(v),
		func(i int) error {
			return m.WriteBool(v[i])
		},
	)
}

func (m *defaultMarshaler) WriteInt(v int) error {
	vx := *(*uint64)(unsafe.Pointer(&v))
	for i := 0; i < 8; i++ {
		b := byte((vx >> (8 * i)) & 0xFF)
		err := m.wr.WriteByte(b)
		if err != nil {
			return errors.Join(ErrMarshalerWrite, err)
		}
	}

	err := m.wr.Flush()
	if err != nil {
		return errors.Join(ErrMarshalerWrite, err)
	}

	return nil
}

func (m *defaultMarshaler) WriteIntOptional(v *int) error {
	return m.writeOptional(
		v != nil,
		func() error {
			return m.WriteInt(*v)
		},
	)
}

func (m *defaultMarshaler) WriteIntArray(v []int) error {
	return m.writeArray(
		len(v),
		func(i int) error {
			return m.WriteInt(v[i])
		},
	)
}

func (m *defaultMarshaler) WriteFloat(v float64) error {
	vx := *(*uint64)(unsafe.Pointer(&v))
	for i := 0; i < 8; i++ {
		b := byte((vx >> (8 * i)) & 0xFF)
		err := m.wr.WriteByte(b)
		if err != nil {
			return errors.Join(ErrMarshalerWrite, err)
		}
	}

	err := m.wr.Flush()
	if err != nil {
		return errors.Join(ErrMarshalerWrite, err)
	}

	return nil
}

func (m *defaultMarshaler) WriteFloatOptional(v *float64) error {
	return m.writeOptional(
		v != nil,
		func() error {
			return m.WriteFloat(*v)
		},
	)
}

func (m *defaultMarshaler) WriteFloatArray(v []float64) error {
	return m.writeArray(
		len(v),
		func(i int) error {
			return m.WriteFloat(v[i])
		},
	)
}

func (m *defaultMarshaler) WriteString(v string) error {
	cap := len(v)
	data := []byte(v)

	err := m.WriteInt(cap)
	if err != nil {
		return err
	}

	for _, c := range data {
		err := m.wr.WriteByte(c)
		if err != nil {
			return errors.Join(ErrMarshalerWrite, err)
		}
	}

	err = m.wr.Flush()
	if err != nil {
		return errors.Join(ErrMarshalerWrite, err)
	}

	return nil
}

func (m *defaultMarshaler) WriteStringOptional(v *string) error {
	return m.writeOptional(
		v != nil,
		func() error {
			return m.WriteString(*v)
		},
	)
}

func (m *defaultMarshaler) WriteStringArray(v []string) error {
	return m.writeArray(
		len(v),
		func(i int) error {
			return m.WriteString(v[i])
		},
	)
}

func (m *defaultMarshaler) WriteObject(obj Object) error {
	err := m.WriteUUID(obj.GetCLSID())
	if err != nil {
		return err
	}

	err = m.WriteUUID(obj.GetInstanceID())
	if err != nil {
		return err
	}

	return nil
}

func (m *defaultMarshaler) WriteObjectOptional(obj Object) error {
	return m.writeOptional(
		!reflect.ValueOf(obj).IsNil(),
		func() error {
			return m.WriteObject(obj)
		},
	)
}

func (m *defaultMarshaler) WriteObjectArray(objs []Object) error {
	return m.writeArray(
		len(objs),
		func(i int) error {
			return m.WriteObject(objs[i])
		},
	)
}

func (m *defaultMarshaler) WriteStructure(v Structure) error {
	return v.Marshal(m)
}

func (m *defaultMarshaler) WriteStructureOptional(v Structure) error {
	return m.writeOptional(
		!reflect.ValueOf(v).IsNil(),
		func() error {
			return m.WriteStructure(v)
		},
	)
}

func (m *defaultMarshaler) WriteStructureArray(v []Structure) error {
	return m.writeArray(
		len(v),
		func(i int) error {
			return m.WriteStructure(v[i])
		},
	)
}

func (m *defaultMarshaler) WriteError(e error) error {
	if e == nil {
		err := m.WriteBool(true)
		if err != nil {
			return err
		}
		return nil
	}

	err := m.WriteBool(false)
	if err != nil {
		return err
	}

	err = m.WriteString(e.Error())
	if err != nil {
		return err
	}

	return nil
}

func (m *defaultMarshaler) writeOptional(present bool, write func() error) error {
	err := m.WriteBool(present)
	if err != nil {
		return err
	}

	if present {
		err = write()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *defaultMarshaler) writeArray(cap int, write func(i int) error) error {
	err := m.WriteInt(cap)
	if err != nil {
		return err
	}

	for i := 0; i < cap; i++ {
		err = write(i)
		if err != nil {
			return err
		}
	}

	return nil
}

type defaultUnmarshaler struct {
	rd *bufio.Reader
}

func NewDefaultUnmarshaler(r io.Reader) Unmarshaler {
	return &defaultUnmarshaler{
		rd: bufio.NewReader(r),
	}
}

func (u *defaultUnmarshaler) ReadUUID() (uuid.UUID, error) {
	b := uuid.UUID{}
	_, err := u.rd.Read(b[:])
	if err != nil {
		return uuid.Nil, err
	}
	return b, nil
}

func (u *defaultUnmarshaler) ReadBool() (bool, error) {
	b, err := u.rd.ReadByte()
	if err != nil {
		return false, errors.Join(ErrUnmarshalerRead, err)
	}

	return b != 0, nil
}

func (u *defaultUnmarshaler) ReadBoolOptional() (*bool, error) {
	var vp *bool
	err := u.readOptional(func() error {
		v, err := u.ReadBool()
		if err != nil {
			return err
		}
		vp = &v
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vp, nil
}

func (u *defaultUnmarshaler) ReadBoolArray() ([]bool, error) {
	var va []bool
	err := u.readArray(
		func(cap int) {
			va = make([]bool, cap)
		},
		func(i int) error {
			var err error
			va[i], err = u.ReadBool()
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return va, nil
}

func (u *defaultUnmarshaler) ReadInt() (int, error) {
	b := make([]byte, 8)
	_, err := u.rd.Read(b)
	if err != nil {
		return 0, errors.Join(ErrUnmarshalerRead, err)
	}

	v := uint64(0)
	for i := 0; i < 8; i++ {
		v = v | (uint64(b[i]) << (8 * i))
	}

	return *(*int)(unsafe.Pointer(&v)), nil
}

func (u *defaultUnmarshaler) ReadIntOptional() (*int, error) {
	var vp *int
	err := u.readOptional(func() error {
		v, err := u.ReadInt()
		if err != nil {
			return err
		}
		vp = &v
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vp, nil
}

func (u *defaultUnmarshaler) ReadIntArray() ([]int, error) {
	var va []int
	err := u.readArray(
		func(cap int) {
			va = make([]int, cap)
		},
		func(i int) error {
			var err error
			va[i], err = u.ReadInt()
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return va, nil
}

func (u *defaultUnmarshaler) ReadFloat() (float64, error) {
	b := make([]byte, 8)
	_, err := u.rd.Read(b)
	if err != nil {
		return 0, errors.Join(ErrUnmarshalerRead, err)
	}

	v := uint64(0)
	for i := 0; i < 8; i++ {
		v = v | (uint64(b[i]) << (8 * i))
	}

	return *(*float64)(unsafe.Pointer(&v)), nil
}

func (u *defaultUnmarshaler) ReadFloatOptional() (*float64, error) {
	var vp *float64
	err := u.readOptional(func() error {
		v, err := u.ReadFloat()
		if err != nil {
			return err
		}
		vp = &v
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vp, nil
}

func (u *defaultUnmarshaler) ReadFloatArray() ([]float64, error) {
	var va []float64
	err := u.readArray(
		func(cap int) {
			va = make([]float64, cap)
		},
		func(i int) error {
			var err error
			va[i], err = u.ReadFloat()
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return va, nil
}

func (u *defaultUnmarshaler) ReadString() (string, error) {
	cap, err := u.ReadInt()
	if err != nil {
		return "", err
	}

	b := make([]byte, cap)
	_, err = u.rd.Read(b)
	if err != nil {
		return "", errors.Join(ErrUnmarshalerRead, err)
	}

	return string(b), nil
}

func (u *defaultUnmarshaler) ReadStringOptional() (*string, error) {
	var vp *string
	err := u.readOptional(func() error {
		v, err := u.ReadString()
		if err != nil {
			return err
		}
		vp = &v
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vp, nil
}

func (u *defaultUnmarshaler) ReadStringArray() ([]string, error) {
	var va []string
	err := u.readArray(
		func(cap int) {
			va = make([]string, cap)
		},
		func(i int) error {
			var err error
			va[i], err = u.ReadString()
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return va, nil
}

func (u *defaultUnmarshaler) ReadObject(f Factory) (Object, error) {
	clsid, err := u.ReadUUID()
	if err != nil {
		return nil, err
	}

	instanceID, err := u.ReadUUID()
	if err != nil {
		return nil, err
	}

	obj, err := f.CreateInstance(clsid, &instanceID)
	if err != nil {
		return nil, errors.Join(ErrUnmarshalerRead, err)
	}

	return obj, nil
}

func (u *defaultUnmarshaler) ReadObjectOptional(f Factory) (Object, error) {
	var vp Object
	err := u.readOptional(func() error {
		v, err := u.ReadObject(f)
		if err != nil {
			return err
		}
		vp = v
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vp, nil
}

func (u *defaultUnmarshaler) ReadObjectArray(f Factory) ([]Object, error) {
	var va []Object
	err := u.readArray(
		func(cap int) {
			va = make([]Object, cap)
		},
		func(i int) error {
			var err error
			va[i], err = u.ReadObject(f)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return va, nil
}

func (u *defaultUnmarshaler) ReadStructure(ref func() Structure) (Structure, error) {
	v := ref()
	err := v.Unmarshal(u)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (u *defaultUnmarshaler) ReadStructureOptional(ref func() Structure) (Structure, error) {
	var vp Structure
	err := u.readOptional(func() error {
		v, err := u.ReadStructure(ref)
		if err != nil {
			return err
		}
		vp = v
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vp, nil
}

func (u *defaultUnmarshaler) ReadStructureArray(ref func() Structure) ([]Structure, error) {
	var va []Structure
	err := u.readArray(
		func(cap int) {
			va = make([]Structure, cap)
		},
		func(i int) error {
			var err error
			va[i], err = u.ReadStructure(ref)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return va, nil
}

func (u *defaultUnmarshaler) ReadError() (error, error) {
	ok, err := u.ReadBool()
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, nil
	}

	msg, err := u.ReadString()
	if err != nil {
		return nil, err
	}

	return fmt.Errorf("remote: %s", msg), nil
}

func (u *defaultUnmarshaler) readOptional(read func() error) error {
	present, err := u.ReadBool()
	if err != nil {
		return err
	}

	if present {
		err := read()
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *defaultUnmarshaler) readArray(init func(cap int), read func(i int) error) error {
	cap, err := u.ReadInt()
	if err != nil {
		return err
	}

	init(cap)
	for i := 0; i < cap; i++ {
		err = read(i)
		if err != nil {
			return err
		}
	}

	return nil
}
