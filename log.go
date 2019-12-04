package virtualbox

// LogFunc is the signature to log traces.
type LogFunc func(string, ...interface{})

func noLog(string, ...interface{}) {}

// Debug is the Logger currently in use.
var Debug LogFunc = noLog
