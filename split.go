package csv

type sliceSplitOption struct {
	delimiter string
}

// SliceSplit は、フィールドがスライスのときにどのように文字列を分割するかを表すオプション。
func SliceSplit(delimiter string) Option {
	return &sliceSplitOption{delimiter}
}

func (o *sliceSplitOption) Apply(r *Reader) error {
	r.sliceDelimiter = o.delimiter
	return nil
}
