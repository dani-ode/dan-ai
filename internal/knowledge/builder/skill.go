package builder

import (
	"fmt"
	"strings"

	"dan-ai/internal/skill/entity"
)

func BuildSkillDocument(skill entity.Skill) (title string, content string) {
	title = fmt.Sprintf("Skill: %s", skill.Technology.Name)
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name: %s\n", skill.Technology.Name))
	
	if skill.Technology.Category != "" {
		sb.WriteString(fmt.Sprintf("Category: %s\n", skill.Technology.Category))
	}
	
	sb.WriteString(fmt.Sprintf("Level: %s\n", skill.Level))
	sb.WriteString(fmt.Sprintf("Years of experience: %.1f\n", skill.Years))
	
	return title, sb.String()
}
