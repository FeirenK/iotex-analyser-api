package db

import "gorm.io/gorm"

type BaseModel struct {
	IndexName   string `gorm:"-"`
	IndexHeight uint64 `gorm:"-"`
}

type IndexHeight struct {
	ID     uint32 `gorm:"primary_key;auto_increment"`
	Name   string `gorm:"size:128;not null;unique"`
	Height uint64 `sql:"type:bigint"`
}

func GetIndexHeight(name string) (uint64, error) {
	var m *IndexHeight
	idx, err := m.ByName(name)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, err
		}
		m = &IndexHeight{
			Name: name,
		}
		return 0, db.Create(m).Error
	}
	return idx.Height, nil
}

func (m *IndexHeight) ByName(name string) (*IndexHeight, error) {
	err := db.Model(m).Where("name = ?", name).Take(&m).Error
	if err != nil {
		return nil, err
	}
	return m, err
}
