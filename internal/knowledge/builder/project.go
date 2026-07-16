package builder

import (
	"fmt"
	"strings"

	"dan-ai/internal/project/entity"
)

func BuildProjectDocument(p entity.Project) (title string, content string) {
	title = fmt.Sprintf("Project: %s", p.Title)
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Title: %s\n", p.Title))
	sb.WriteString(fmt.Sprintf("Slug: %s\n", p.Slug))
	sb.WriteString(fmt.Sprintf("Description: %s\n", p.Description))
	
	if p.Summary != "" {
		sb.WriteString(fmt.Sprintf("Summary: %s\n", p.Summary))
	}
	if p.Architecture != "" {
		sb.WriteString(fmt.Sprintf("Architecture: %s\n", p.Architecture))
	}
	if p.RepositoryURL != "" {
		sb.WriteString(fmt.Sprintf("Repository URL: %s\n", p.RepositoryURL))
	}
	if p.DemoURL != "" {
		sb.WriteString(fmt.Sprintf("Demo URL: %s\n", p.DemoURL))
	}
	if p.Status != "" {
		sb.WriteString(fmt.Sprintf("Status: %s\n", p.Status))
	}
	
	return title, sb.String()
}
