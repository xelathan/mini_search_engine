package db

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"-"`
	IsAdmin   bool       `json:"isAdmin" gorm:"default:false"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (u *User) CreateAdmin() error {
	user := User{
		Email:    "admin@gmail.com",
		Password: "admin",
		IsAdmin:  true,
	}
	// Hash Password
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(password)

	if err := DBConn.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) LoginAsAdmin(email, password string) (*User, error) {
	if err := DBConn.Where("email = ? AND is_admin = ?", email, true).First(&u).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, err
	}

	return u, nil
}

// type Settings struct {
// 	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
// 	Amount     int
// 	SearchOn   bool
// 	AddNewUrls bool
// 	UserID     string `gorm:"type:uuid"`
// }
