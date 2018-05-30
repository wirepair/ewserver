package converter

// ActionsRegex takes a list of actions for permissions and converts them to a regex
// returns an empty string if invalid
func ActionsRegex(inputs []string) string {
	if inputs == nil || len(inputs) == 0 {
		return ""
	}

	regex := "("
	for _, value := range inputs {
		switch value {
		case "ALL":
			return ".*"
		case "GET":
			regex += "GET|"
		case "HEAD":
			regex += "HEAD|"
		case "OPTIONS":
			regex += "OPTIONS|"
		case "PATCH":
			regex += "PATCH|"
		case "POST":
			regex += "POST|"
		case "PUT":
			regex += "PUT|"
		case "DELETE":
			regex += "DELETE|"
		}
	}
	// replace final | with )
	regex = regex[:len(regex)-1] + ")"
	return regex
}
