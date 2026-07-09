// internal/profile/entity/profile.go
package entity

import "time"

type Profile struct {
	ID           string    `gorm:"type:char(26);primaryKey"`
	FullName     string    `gorm:"type:text;not null"`
	Headline     string    `gorm:"type:text"`
	Bio          string    `gorm:"type:text"`
	Email        string    `gorm:"type:text"`
	Phone        string    `gorm:"type:text"`
	Location     string    `gorm:"type:text"`
	Github       string    `gorm:"type:text"`
	Linkedin     string    `gorm:"type:text"`
	Website      string    `gorm:"type:text"`
	Avatar       string    `gorm:"type:text"`
	ResumeURL    string    `gorm:"type:text"`
	Availability string    `gorm:"type:text"`
	Timezone     string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

// TableName overrides the table name used by Profile to profiles.
func (Profile) TableName() string {
	return "profiles"
}
