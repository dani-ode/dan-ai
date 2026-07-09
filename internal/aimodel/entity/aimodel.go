// internal/aimodel/entity/aimodel.go
package entity

type AIModel struct {
	ID             string  `gorm:"type:char(26);primaryKey"`
	Name           string  `gorm:"type:text;uniqueIndex;not null"`
	Provider       string  `gorm:"type:text;not null"`
	Temperature    float64 `gorm:"type:numeric(3,2);default:0.7"`
	MaxTokens      int32   `gorm:"type:int"`
	ContextWindow  int32   `gorm:"type:int"`
	SupportsTools  bool    `gorm:"type:boolean;default:false"`
	SupportsStream bool    `gorm:"type:boolean;default:false"`
	Enabled        bool    `gorm:"type:boolean;default:true"`
}

// TableName overrides the table name used by AIModel to ai_models.
func (AIModel) TableName() string {
	return "ai_models"
}
