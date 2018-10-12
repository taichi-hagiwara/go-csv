package csv

import (
	"encoding/csv"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// Reader はCSVを読み込む。
type Reader struct {
	reader         *csv.Reader
	index          map[string]int
	sliceDelimiter string
	closer         io.Closer
}

// ReadLine は、CSVを1行読み込む。
func (r *Reader) ReadLine(v interface{}) error {
	e := reflect.ValueOf(v).Elem()
	record, err := r.reader.Read()
	if err == io.EOF {
		return io.EOF
	} else if err != nil {
		return errors.Wrap(err, "failed to read a record")
	}

	t := e.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		csv := field.Tag.Get("csv")
		if csv != "" {
			index, ok := r.index[csv]
			if ok && index < len(record) {
				r.setValue(e.FieldByName(field.Name), record[r.index[csv]])
			}
		}
	}

	return nil
}

func (r *Reader) setValue(f reflect.Value, value string) {
	switch f.Type().Kind() {
	case reflect.String:
		f.SetString(value)
		// TODO: not string fields
	case reflect.Ptr:
		v := reflect.New(f.Type().Elem())
		r.setValue(v, value)
		f.Set(v)
	case reflect.Slice:
		if r.sliceDelimiter != "" {
			s := strings.Split(value, r.sliceDelimiter)
			slice := reflect.MakeSlice(reflect.SliceOf(f.Type().Elem()), len(s), len(s))
			for i, v := range s {
				r.setValue(slice.Index(i), v)
			}
			f.Set(slice)
		}
	}
}

// Close は、Readerを閉じる。
func (r *Reader) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}

// NewReader は、Reader を作成する。
func NewReader(r *csv.Reader, options ...Option) (*Reader, error) {
	reader := &Reader{
		reader: r,
	}
	for _, o := range options {
		if err := o.Apply(reader); err != nil {
			return nil, errors.Wrapf(err, "failed to apply option: %s", o)
		}
	}

	if reader.index == nil {
		record, err := r.Read()
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read first line")
		}
		Index(record).Apply(reader)
	}
	return reader, nil
}

// FromIOReader は、io.Reader から Reader を初期化する。
func FromIOReader(r io.Reader, options ...Option) (*Reader, error) {
	return NewReader(csv.NewReader(r), options...)
}

// FromFile は、ファイルから Reader を初期化する。
func FromFile(name string, options ...Option) (*Reader, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", name)
	}

	reader, err := FromIOReader(file, options...)
	if err != nil {
		file.Close()
		return nil, errors.Wrapf(err, "failed to create Reader")
	}

	reader.closer = file
	return reader, nil
}
