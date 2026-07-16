package builder

import (
	"fmt"
	"strings"

	"dan-ai/internal/profile/entity"
)

func BuildProfileDocument(prof entity.Profile) (title string, content string) {
	title = fmt.Sprintf("Profile: %s", prof.FullName)
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name: %s\n", prof.FullName))
	sb.WriteString(fmt.Sprintf("Email: %s\n", prof.Email))
	sb.WriteString(fmt.Sprintf("Headline: %s\n", prof.Headline))
	sb.WriteString(fmt.Sprintf("Bio: %s\n", prof.Bio))
	
	if prof.Phone != "" {
		sb.WriteString(fmt.Sprintf("Phone: %s\n", prof.Phone))
	}
	if prof.Location != "" {
		sb.WriteString(fmt.Sprintf("Location: %s\n", prof.Location))
	}
	if prof.Github != "" {
		sb.WriteString(fmt.Sprintf("Github: %s\n", prof.Github))
	}
	if prof.Linkedin != "" {
		sb.WriteString(fmt.Sprintf("LinkedIn: %s\n", prof.Linkedin))
	}
	if prof.Website != "" {
		sb.WriteString(fmt.Sprintf("Website: %s\n", prof.Website))
	}
	
	return title, sb.String()
}
