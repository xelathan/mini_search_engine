package db

import "time"

type SearchSetting struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SearchOn   bool      `json:"searchOn"`
	AddNewUrls bool      `json:"addNewUrls"`
	Amount     int       `json:"amount"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (s *SearchSetting) Get() error {
	if err := DBConn.Where("id = ?", 1).First(&s).Error; err != nil {
		return err
	}
	return nil
}
func (s *SearchSetting) Update() error {
	if err := DBConn.Select("search_on", "add_new", "amount", "updated_at").Where("id = 1").Updates(s).Error; err != nil {
		return err
	}
	return nil
}

func (s *SearchSetting) Create() error {
	if err := DBConn.Create(&s).Error; err != nil {
		return err
	}
	return nil
}
