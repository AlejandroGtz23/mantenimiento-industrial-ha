package repository

import (

	//"strings"

	"fiber3-app/backend/database"
	"fiber3-app/backend/domain"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{db: database.DB}
}

func (r *ProfileRepository) GetAll() []domain.Profile {
	var profile []domain.Profile

	r.db.Raw(`
		SELECT UUID_TO_CHAR(ID) ID, "FIRSTNAME", "LASTNAME", USERNAME, EMAIL, SKILLS, TIMEZONE
		FROM PROFILE
	`).Scan(&profile)

	return profile
}

func (r *ProfileRepository) GetByID(id string) *domain.Profile {
	var profile domain.Profile

	r.db.Raw(`
		SELECT UUID_TO_CHAR(ID) ID, "FIRSTNAME", "LASTNAME", USERNAME, EMAIL, SKILLS, TIMEZONE
		FROM PROFILE
		WHERE ID = CHAR_TO_UUID(?)`, id).Scan(&profile)

	return &profile
}

func (r *ProfileRepository) Create(p *domain.Profile) error {
	qry := `
		INSERT INTO PROFILE (ID, "FIRSTNAME", "LASTNAME", USERNAME, "PASSWORD", EMAIL,  SKILLS, TIMEZONE)
		VALUES (CHAR_TO_UUID(?), ?, ?, ?, ?, ?, ?, ?)
	`
	return r.db.Exec(qry, p.ID, p.Firstname, p.Lastname, p.Username, p.Password, p.Email, p.Skills, p.Timezone).Error
}

func (r *ProfileRepository) Update(p *domain.Profile) error {
	qry := `
		UPDATE PROFILE SET
			"FIRSTNAME" = ?,
			"LASTNAME" = ?,
			EMAIL = ?,
			USERNAME = ?,
			SKILLS = ?,
			TIMEZONE = ?
		WHERE ID = CHAR_TO_UUID(?)
	`
	return r.db.Exec(qry, p.Firstname, p.Lastname, p.Email, p.Username, p.Skills, p.Timezone).Error
}

func (r *ProfileRepository) Delete(id string) error {
	qry := `DELETE FROM PROFILE WHERE ID = CHAR_TO_UUID(?)`

	return r.db.Exec(qry, id).Error
}
