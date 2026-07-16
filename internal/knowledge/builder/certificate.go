package builder

import (
	"fmt"
	"strings"

	"dan-ai/internal/certificate/entity"
)

func BuildCertificateDocument(cert entity.Certificate) (title string, content string) {
	title = fmt.Sprintf("Certificate: %s by %s", cert.Title, cert.Issuer)
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Title: %s\n", cert.Title))
	sb.WriteString(fmt.Sprintf("Issuer: %s\n", cert.Issuer))
	
	if cert.IssueDate != nil {
		sb.WriteString(fmt.Sprintf("Issue Date: %s\n", cert.IssueDate.Format("2006-01-02")))
	}
	if cert.ExpirationDate != nil {
		sb.WriteString(fmt.Sprintf("Expiration Date: %s\n", cert.ExpirationDate.Format("2006-01-02")))
	}
	if cert.CredentialURL != "" {
		sb.WriteString(fmt.Sprintf("Credential URL: %s\n", cert.CredentialURL))
	}
	if len(cert.Skills) > 0 {
		sb.WriteString(fmt.Sprintf("Skills validated: %s\n", strings.Join(cert.Skills, ", ")))
	}
	
	return title, sb.String()
}
