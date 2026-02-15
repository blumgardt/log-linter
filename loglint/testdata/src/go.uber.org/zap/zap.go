package zap

type Field struct{}

type Logger struct{}
type SugaredLogger struct{}

func NewProduction() (*Logger, error) { return &Logger{}, nil }

func (l *Logger) Info(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)  {}
func (l *Logger) Error(msg string, fields ...Field) {}

func (l *Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

func (s *SugaredLogger) Infow(msg string, keysAndValues ...any)  {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...any)  {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...any) {}

func (s *SugaredLogger) Infof(template string, args ...any)  {}
func (s *SugaredLogger) Warnf(template string, args ...any)  {}
func (s *SugaredLogger) Errorf(template string, args ...any) {}

func String(key, val string) Field  { return Field{} }
func Any(key string, val any) Field { return Field{} }
