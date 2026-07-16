package builder

import (
	"fmt"
	"strings"

	"dan-ai/internal/experience/entity"
)

func BuildExperienceDocument(exp entity.Experience) (title string, content string) {
	title = fmt.Sprintf("Experience: %s at %s", exp.Position, exp.Company)
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Position: %s\n", exp.Position))
	sb.WriteString(fmt.Sprintf("Company: %s\n", exp.Company))
	sb.WriteString(fmt.Sprintf("Location: %s\n", exp.Location))
	sb.WriteString(fmt.Sprintf("Description: %s\n", exp.Description))
	
	if exp.StartDate != nil {
		sb.WriteString(fmt.Sprintf("Start Date: %s\n", exp.StartDate.Format("2006-01-02")))
	}
	if exp.EndDate != nil {
		sb.WriteString(fmt.Sprintf("End Date: %s\n", exp.EndDate.Format("2006-01-02")))
	}
	if len(exp.Skills) > 0 {
		sb.WriteString(fmt.Sprintf("Skills used: %s\n", strings.Join(exp.Skills, ", ")))
	}
	
	return title, sb.String()
}
