package loglint

import (
	"github.com/blumgardt/log-linter/internal/config"
	loglintanalyzer "github.com/blumgardt/log-linter/loglint"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type Settings struct {
	EnableRule1LowercaseStart *bool `json:"enable_rule1_lowercase_start"`
	EnableRule2ASCIIOnly      *bool `json:"enable_rule2_ascii_only"`
	EnableRule3NoSpecialChars *bool `json:"enable_rule3_no_special_chars"`
	EnableRule4Sensitive      *bool `json:"enable_rule4_sensitive"`

	ExtraAllowedChars string   `json:"extra_allowed_chars"`
	SensitiveKeywords []string `json:"sensitive_keywords"`
	LogMethodNames    []string `json:"log_method_names"`
}

type Plugin struct{}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	config.ApplyPluginSettings(config.PluginSettings{
		EnableRule1LowercaseStart: s.EnableRule1LowercaseStart,
		EnableRule2ASCIIOnly:      s.EnableRule2ASCIIOnly,
		EnableRule3NoSpecialChars: s.EnableRule3NoSpecialChars,
		EnableRule4Sensitive:      s.EnableRule4Sensitive,
		ExtraAllowedChars:         s.ExtraAllowedChars,
		SensitiveKeywords:         s.SensitiveKeywords,
		LogMethodNames:            s.LogMethodNames,
	})

	return &Plugin{}, nil
}

func (*Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{loglintanalyzer.Analyzer}, nil
}

func (*Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
