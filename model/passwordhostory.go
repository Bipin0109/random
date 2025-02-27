
package model 


type PasswordHistory struct {
    ID               int     `gorm:"primaryKey" json:"id"`
    Username         string  `gorm:"unique;not null" json:"username"`
    PreviousPassword1 *string `gorm:"column:previous_password1" json:"previous_password_1"`
    PreviousPassword2 *string `gorm:"column:previous_password2" json:"previous_password_2"`
    PreviousPassword3 *string `gorm:"column:previous_password3" json:"previous_password_3"`
}

func NewPasswordHistory(username string, PreviousPassword1 , PreviousPassword2 , PreviousPassword3 string) PasswordHistory {
	return PasswordHistory{
		Username: username,
		PreviousPassword1: &PreviousPassword1,
		PreviousPassword2: &PreviousPassword2,
		PreviousPassword3: &PreviousPassword3,
	}
}
func (PasswordHistory) TableName() string {
	return "password_history"
}
