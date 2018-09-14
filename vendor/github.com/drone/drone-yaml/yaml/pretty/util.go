package pretty

func isPrimative(v interface{}) bool {
	switch v.(type) {
	case bool, string, int, float64:
		return true
	default:
		return false
	}
}

func isSlice(v interface{}) bool {
	switch v.(type) {
	case []interface{}:
		return true
	case []string:
		return true
	default:
		return false
	}
}

func isZero(v interface{}) bool {
	switch v := v.(type) {
	case bool:
		return v == false
	case string:
		return len(v) == 0
	case int:
		return v == 0
	case float64:
		return v == 0
	case []interface{}:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case map[interface{}]interface{}:
		return len(v) == 0
	case map[string]string:
		return len(v) == 0
	default:
		return false
	}
}

func isQuoted(b rune) bool {
	switch b {
	case '#', ',', '[', ']', '{', '}', '&', '*', '!', '|', '>', '\'', '"', '%', '@', '`':
		return true
	case '\a', '\b', '\f', '\n', '\r', '\t', '\v':
		return true
	default:
		return false
	}
}

func chunk(s string, chunkSize int) []string {
	if len(s) == 0 {
		return []string{s}
	}
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		nn := i + chunkSize
		if nn > len(s) {
			nn = len(s)
		}
		chunks = append(chunks, s[i:nn])
	}
	return chunks
}
