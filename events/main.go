package events

// Could not make type system work with unknown ponter in the func argument
var EventHandlers = []any{
	ReadyEventHandler,
	VoiceStateUpdateHandler,
	// MessageCreateHandler,
}
