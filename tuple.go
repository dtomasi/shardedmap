package shardedmap

func NewTuple(key string, value interface{}) Tuple {
	return Tuple{
		key,
		value,
	}
}

type Tuple struct {
	key   string
	value interface{}
}

func (t Tuple) GetKey() string {
	return t.key
}

func (t Tuple) GetValue() interface{} {
	return t.value
}
