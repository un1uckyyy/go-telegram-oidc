// here you should implement methods to interact with your internal database with user's data

package storage

type Db struct {
	data map[int64]string
}

func NewDb() *Db {
	return &Db{
		data: make(map[int64]string),
	}
}

func (d *Db) GetUser(telegramId int64) (string, bool) {
	userId, ok := d.data[telegramId]
	if !ok {
		return "", false
	}
	return userId, true
}

func (d *Db) SetUser(userId string, telegramId int64) error {
	d.data[telegramId] = userId
	return nil
}
