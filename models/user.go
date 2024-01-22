package models

import (
	"errors"
	"html"
	"jwtgin/utils/token"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}

func GetUserById(uid uint) (User, error) {
	var u User
	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found")
	}
	u.PrepareGive()
	return u, nil
}
func (u *User) PrepareGive() {
	u.Password = ""
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, error) {
	var err error
	u := User{}
	err = DB.Model(User{}).Where("username=?", username).Take(&u).Error
	if err != nil {
		return "", err
	}
	err = VerifyPassword(password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	token, err := token.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *User) SaveUser() (*User, error) {
	var err error
	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	//remove spaces in username
	usrnm := html.EscapeString(strings.TrimSpace(u.Username))
	u.Username = usrnm
	return nil
}
