package domain

type Profile struct {
	ID        string `gorm:"column ID" json:"id"`
	Firstname string `gorm:"column FIRSTNAME" json:"firstname"`
	Lastname  string `gorm:"column LASTNAME" json:"lastname"`
	Username  string `gorm:"column USERNAME" json:"username"`
	Password  string `gorm:"column PASSWORD" json:"-"`
	Email     string `gorm:"column EMAIL" json:"email"`
	Skills    string `gorm:"column SKILLS" json:"skills"`
	Timezone  string `gorm:"column TIMEZONE" json:"timezone"`
}
