package db

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"snippet_box/pkg/models"
)

type UserModel struct {
	DB *gorm.DB
}

func (userModel *UserModel) Insert(name, email, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user := models.Users{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	res := userModel.DB.Create(&user)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (userModel *UserModel) Login(email, password string) (*models.Users, error) {

	user := &models.Users{}
	result := userModel.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userModel *UserModel) IsEmailAlreadyTaken(email string) error {
	user := models.Users{}
	result := userModel.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
