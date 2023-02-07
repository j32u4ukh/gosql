package cntr

type Iterator struct {
	datas  []any
	index  int
	length int
}

func NewIterator(datas []any) *Iterator {
	it := &Iterator{
		datas:  datas,
		index:  -1,
		length: len(datas),
	}
	return it
}

func (it *Iterator) HasNext() bool {
	it.index += 1
	return it.index < it.length
}

func (it *Iterator) Next() any {
	return it.datas[it.index]
}
