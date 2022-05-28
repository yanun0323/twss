package delivery

func MessageResponse(msg string) map[string]any {
	return map[string]any{"message": msg}
}

func ErrorResponse(err error) map[string]any {
	return map[string]any{"message": err.Error()}
}
