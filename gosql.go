package gosql

type FuncQuery func(columns []string, datas [][]string, generator func() any) (objs []any, err error)
