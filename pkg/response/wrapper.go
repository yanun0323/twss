package response

func Message(msg string) map[string]any {
	return map[string]any{"message": msg}
}

func Error(err error) map[string]any {
	return map[string]any{"message": err.Error()}
}
