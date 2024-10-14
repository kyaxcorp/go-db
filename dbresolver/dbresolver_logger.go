package dbresolver

import "github.com/rs/zerolog"

/*
We should enable logging throughout the entire client, that will help us do a better debug

1. Add Enable/Disable logging
2. Add different log levels
3. Add a name of the log instance... usually it should take from the config or somehow?!... or we can setup also manually a name
4.
*/

// LDebug -> 0
func (dr *DBResolver) LDebug() *zerolog.Event {
	return dr.Logger.Debug()
}

// LInfo -> 1
func (dr *DBResolver) LInfo() *zerolog.Event {
	return dr.Logger.Info()
}

// LWarn -> 2
func (dr *DBResolver) LWarn() *zerolog.Event {
	return dr.Logger.Warn()
}

// LError -> 3
func (dr *DBResolver) LError() *zerolog.Event {
	return dr.Logger.Error()
}

// LFatal -> 4
func (dr *DBResolver) LFatal() *zerolog.Event {
	return dr.Logger.Fatal()
}

// LPanic -> 5
func (dr *DBResolver) LPanic() *zerolog.Event {
	return dr.Logger.Panic()
}

//

//-------------------------------------\\

//

func (dr *DBResolver) LEvent(eventType string, eventName string, beforeMsg func(event *zerolog.Event)) {
	dr.Logger.InfoEvent(eventType, eventName, beforeMsg)
}

func (dr *DBResolver) LEventCustom(eventType string, eventName string) *zerolog.Event {
	return dr.Logger.InfoEventCustom(eventType, eventName)
}

func (dr *DBResolver) LEventF(eventType string, eventName string, functionName string) *zerolog.Event {
	return dr.Logger.InfoEventF(eventType, eventName, functionName)
}

//

//-------------------------------------\\

//

// LWarnF -> when you need specifically to indicate in what function the logging is happening
func (dr *DBResolver) LWarnF(functionName string) *zerolog.Event {
	return dr.Logger.WarnF(functionName)
}

// LInfoF -> when you need specifically to indicate in what function the logging is happening
func (dr *DBResolver) LInfoF(functionName string) *zerolog.Event {
	return dr.Logger.InfoF(functionName)
}

// LDebugF -> when you need specifically to indicate in what function the logging is happening
func (dr *DBResolver) LDebugF(functionName string) *zerolog.Event {
	return dr.Logger.DebugF(functionName)
}

// LErrorF -> when you need specifically to indicate in what function the logging is happening
func (dr *DBResolver) LErrorF(functionName string) *zerolog.Event {
	return dr.Logger.ErrorF(functionName)
}
