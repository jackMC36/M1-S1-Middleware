package config

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/adrg/frontmatter"
)

type MailMatter struct {
	Subject string `yaml:"subject" toml:"subject" json:"subject"`
}

//go:embed *.tmpl
var embeddedTemplates embed.FS

func GetStringFromEmbeddedTemplate(templatePath string, data any) (content string, matter MailMatter, err error) {
	temp, err := template.ParseFS(embeddedTemplates, templatePath)
	if err != nil {
		return "", matter, err
	}

	var tpl bytes.Buffer
	if err := temp.Execute(&tpl, data); err != nil {
		return "", matter, err
	}

	mailContent, parseErr := frontmatter.Parse(strings.NewReader(tpl.String()), &matter)
	if parseErr != nil {
		return tpl.String(), matter, nil
	}

	return string(mailContent), matter, nil
}
