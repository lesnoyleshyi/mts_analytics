package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const appName = `analytics`
const hostIP = `localhost`

// DynamicLogLevel is a zap.AtomicLevel, that can change logging level
// in runtime. We can change it via http-request.
var DynamicLogLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

// NewStdEncoder returns standard EncoderConfig.
// All configurations should be created by this func to sustain uniformity of logs.
// Minor changes (e.g. time representation) could be done afterwards\
// on particular config variable
func NewStdEncoder() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey: "message",
		LevelKey:   "level",
		TimeKey:    "timestamp",
		//NameKey:             "logger_name",
		//CallerKey:           zapcore.OmitKey,
		NameKey:             zapcore.OmitKey,
		CallerKey:           "logger_name",
		FunctionKey:         zapcore.OmitKey,
		StacktraceKey:       "stack_trace",
		SkipLineEnding:      false,
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.RFC3339TimeEncoder,
		EncodeDuration:      zapcore.SecondsDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		EncodeName:          zapcore.FullNameEncoder,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "\t",
	}
}

func NewLogger() *zap.Logger {
	// Configure encoders from std prototype
	//
	//machineEncoderConfig := NewStdEncoder()

	humanEncoderConfig := NewStdEncoder()
	humanEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Create different encoders to generate human-readable (NewConsoleEncoder)
	// and machine-readable (NewJSONEncoder) logs
	//
	//sentryEncoder := zapcore.NewJSONEncoder(machineEncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(humanEncoderConfig)

	// Create all our WriteSyncer-s == source to which we write logs
	//
	//sentryOut := zapcore.AddSync(sentry_wrapper_implementing_ioWriter)
	stdout := zapcore.Lock(os.Stdout)
	//stderr := zapcore.Lock(os.Stderr)

	// Define functions to decide whether to log message or not.
	// errLvl - write to stderr and sentry
	// lowPriorityLvl - write to stdout
	//errLvl := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
	//	return level >= zapcore.ErrorLevel
	//})
	//lowPriorityLvl := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
	//	return level < zapcore.ErrorLevel
	//})

	// Chain all separate Cores into one Core
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, DynamicLogLevel),
		//zapcore.NewCore(sentryEncoder, stdout, DynamicLogLevel), //test
		//we can route log messages to stderr and stdout separately
		// depending on their level
		//zapcore.NewCore(consoleEncoder, stderr, errLvl),
		//zapcore.NewCore(consoleEncoder, stdout, lowPriorityLvl),
		//
		// We can also write logs to Sentry depending on their level
		//zapcore.NewCore(sentryEncoder, sentryOut, errLvl),
	)

	// if we want to implement sampling of logs, we can wrap our core in sampledLogger
	//
	//sampledLogger := zapcore.NewSamplerWithOptions(core, time.Second, 10, 5,
	//	// here could be hooks (funcs which will be called when Sampler makes a decision)
	//	//zapcore.SamplerHook(),
	//	)

	// Create logger with required fields that should be equal in each log message
	l := zap.New(core).WithOptions(
		zap.Fields(
			zap.String("app_name", appName),
			zap.String("host_ip", hostIP),
		),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.ErrorOutput(stdout),
		zap.AddCaller(),
		//zap.ErrorOutput(stderr),
	)

	// it's another simpler way of configuring logger
	//
	//zap.Config{
	//	Level:             zap.AtomicLevel{},
	//	Development:       false,
	//	DisableCaller:     false,
	//	DisableStacktrace: false,
	//	Sampling:          nil,
	//	Encoding:          "",
	//	EncoderConfig: zapcore.EncoderConfig{
	//		MessageKey:          "",
	//		LevelKey:            "",
	//		TimeKey:             "",
	//		NameKey:             "",
	//		CallerKey:           "",
	//		FunctionKey:         "",
	//		StacktraceKey:       "",
	//		SkipLineEnding:      false,
	//		LineEnding:          "",
	//		EncodeLevel:         nil,
	//		EncodeTime:          nil,
	//		EncodeDuration:      nil,
	//		EncodeCaller:        nil,
	//		EncodeName:          nil,
	//		NewReflectedEncoder: nil,
	//		ConsoleSeparator:    "",
	//	},
	//	OutputPaths:      nil,
	//	ErrorOutputPaths: nil,
	//	InitialFields:    nil,
	//}
	//var cfg zap.Config
	//_ = yaml.Unmarshal([]byte(""), &cfg)

	return l
}
