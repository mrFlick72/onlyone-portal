package logging

import (
	"go.uber.org/zap/zapcore"
)

type levelledCore struct {
	zapcore.Core
	level zapcore.LevelEnabler
}

func (c *levelledCore) Enabled(l zapcore.Level) bool {
	return c.level.Enabled(l)
}

func (c *levelledCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !c.Enabled(ent.Level) {
		return ce
	}
	return c.Core.Check(ent, ce)
}

func (c *levelledCore) With(fields []zapcore.Field) zapcore.Core {
	return &levelledCore{Core: c.Core.With(fields), level: c.level}
}
