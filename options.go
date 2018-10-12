package csv

// Option は、CSV読み込み時のオプションを表す。
type Option interface {
	Apply(r *Reader) error
}

type indexOption struct {
	index []string
}

// Index は、CSVの1行目に列のタイトルがない場合、それを表すオプション。
func Index(index []string) Option {
	return &indexOption{index}
}

func (o *indexOption) Apply(r *Reader) error {
	m := make(map[string]int)
	for i, v := range o.index {
		m[v] = i
	}
	r.index = m
	return nil
}

func (o *indexOption) String() string {
	return "Index"
}
