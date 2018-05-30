package converter

import "testing"

func TestActionsRegex(t *testing.T) {
	actions := []string{"GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "DELETE"}
	regex := ActionsRegex(actions)
	if regex != "(GET|HEAD|OPTIONS|PATCH|POST|PUT|DELETE)" {
		t.Fatalf("expected (GET|HEAD|OPTIONS|PATCH|POST|PUT|DELETE) got: %s\n", regex)
	}

	regex = ActionsRegex([]string{"GET", "POST"})
	if regex != ("(GET|POST)") {
		t.Fatalf("expected (GET|POST) got: %s\n", regex)
	}

	regex = ActionsRegex([]string{})
	if regex != "" {
		t.Fatalf("expected empty string for invalid actions list, got: %s\n", regex)
	}

	regex = ActionsRegex([]string{"GET", "POST", "ALL"})
	if regex != ".*" {
		t.Fatalf("expected .* for actions containing ALL, got: %s\n", regex)
	}

	regex = ActionsRegex([]string{"GET", "asdf"})
	if regex != "(GET)" {
		t.Fatalf("expected (GET) for actions containing invalid action, got: %s\n", regex)
	}

}
