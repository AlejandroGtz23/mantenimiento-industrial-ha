package repository

import (
	"fiber3-app/backend/database"
	"fmt"

	"gorm.io/gorm"
)

type userAccount struct {
	ID        string `gorm:"ID" json:"id"`
	Firstname string `gorm:"FIRSTNAME" json:"firstname"`
	Lastname  string `gorm:"LASTNAME"  json:"lastname"`
	ImageSrc  string `json:"imageSrc"`
	Password  string `gorm:"PASSWORD" json:"-"`
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{db: database.DB}
}

func (r *AuthRepository) UserAccount(username string) *userAccount {
	var userAcc userAccount

	r.db.Raw(`
		SELECT UUID_TO_CHAR(ID) ID, "FIRSTNAME", "LASTNAME", PASSWORD
		FROM PROFILE
		WHERE USERNAME = ?`, username).Scan(&userAcc)

	userAcc.ImageSrc = fmt.Sprintf("/images/users/%s.jpg", userAcc.ID)
	return &userAcc
}

func (r *AuthRepository) UpdatePasswd(userName string, hashPasswd string) error {
	qry := `
		UPDATE PROFILE SET
    		"PASSWORD" = ?
		WHERE USERNAME = ?	
	`
	return r.db.Exec(qry, hashPasswd, userName).Error
}
