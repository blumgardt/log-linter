package config

import "sync"

var (
	mu sync.RWMutex

	defaultLogMethodNames = []string{
		"Debug", "Info", "Warn", "Error", "Fatal", "Panic",
		"Debugw", "Infow", "Warnw", "Errorw", "Fatalw", "Panicw",
		"Debugf", "Infof", "Warnf", "Errorf", "Fatalf", "Panicf",
	}

	defaultSensitiveKeywords = []string{
		"password", "passwd", "psswrd", "pswrd", "pwd",
		"token", "jwt", "bearer",
		"secret",
		"api_key", "apikey", "api-key",
		"authorization", "auth",
		"cookie", "session",
		"private_key", "privatekey", "ssh_key",
	}

	LogMethodNames    = make(map[string]struct{})
	SensitiveKeywords []string
	ExtraAllowedChars string
	EnableRule1       = true
	EnableRule2       = true
	EnableRule3       = true
	EnableRule4       = true
)

func init() {
	resetToDefaults()
}

func resetToDefaults() {
	LogMethodNames = make(map[string]struct{}, len(defaultLogMethodNames))
	for _, n := range defaultLogMethodNames {
		LogMethodNames[n] = struct{}{}
	}
	SensitiveKeywords = append([]string(nil), defaultSensitiveKeywords...)
	ExtraAllowedChars = ""
	EnableRule1 = true
	EnableRule2 = true
	EnableRule3 = true
	EnableRule4 = true
}

type PluginSettings struct {
	EnableRule1LowercaseStart *bool
	EnableRule2ASCIIOnly      *bool
	EnableRule3NoSpecialChars *bool
	EnableRule4Sensitive      *bool

	ExtraAllowedChars string
	SensitiveKeywords []string
	LogMethodNames    []string
}

func ApplyPluginSettings(s PluginSettings) {
	mu.Lock()
	defer mu.Unlock()

	if s.EnableRule1LowercaseStart != nil {
		EnableRule1 = *s.EnableRule1LowercaseStart
	}
	if s.EnableRule2ASCIIOnly != nil {
		EnableRule2 = *s.EnableRule2ASCIIOnly
	}
	if s.EnableRule3NoSpecialChars != nil {
		EnableRule3 = *s.EnableRule3NoSpecialChars
	}
	if s.EnableRule4Sensitive != nil {
		EnableRule4 = *s.EnableRule4Sensitive
	}

	if s.ExtraAllowedChars != "" {
		ExtraAllowedChars = s.ExtraAllowedChars
	}

	if len(s.SensitiveKeywords) > 0 {
		SensitiveKeywords = append([]string(nil), s.SensitiveKeywords...)
	}

	if len(s.LogMethodNames) > 0 {
		LogMethodNames = make(map[string]struct{}, len(s.LogMethodNames))
		for _, n := range s.LogMethodNames {
			LogMethodNames[n] = struct{}{}
		}
	}
}

func IsRule1Enabled() bool { mu.RLock(); defer mu.RUnlock(); return EnableRule1 }
func IsRule2Enabled() bool { mu.RLock(); defer mu.RUnlock(); return EnableRule2 }
func IsRule3Enabled() bool { mu.RLock(); defer mu.RUnlock(); return EnableRule3 }
func IsRule4Enabled() bool { mu.RLock(); defer mu.RUnlock(); return EnableRule4 }

func GetExtraAllowedChars() string {
	mu.RLock()
	defer mu.RUnlock()
	return ExtraAllowedChars
}
