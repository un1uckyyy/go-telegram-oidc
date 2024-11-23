// here you can use any persistent storage for temporary tickets saving instead of default map
// like mongodb, postgres, redis etc.
// token is intended for one-time use, so you should delete it after it has been used
// it would also be nice to make it expiring

package ticket

type KeyValue struct {
	data map[string]int64
}

func NewKeyValue() *KeyValue {
	return &KeyValue{
		data: make(map[string]int64),
	}
}

func (k *KeyValue) GetTmpUserInfo(state string) (userId int64, ok bool) {
	val, ok := k.data[state]
	return val, ok
}

func (k *KeyValue) SetTmpUserInfo(state string, userId int64) error {
	k.data[state] = userId
	return nil
}

func (k *KeyValue) PopTmpUserInfo(state string) (userId int64, ok bool) {
	val, ok := k.data[state]
	delete(k.data, state)
	return val, ok
}
