package main_test

import (
	"os/exec"
	"testing"
	"time"
)

func TestInstructions(t *testing.T) {
	if err := exec.Command("go", "build", ".").Run(); err != nil {
		t.Fatal(err)
	}

	server := exec.Command("./emojis")
	server.Env = append(server.Environ(), "SERVICEWEAVER_CONFIG=config.toml")
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := server.Process.Kill(); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(time.Second)

	for _, test := range []struct{ url, want string }{
		{"localhost:9000/search?q=sushi", `["🍣"]` + "\n"},
		{"localhost:9000/search?q=curry", `["🍛"]` + "\n"},
		{"localhost:9000/search?q=shrimp", `["🍤","🦐"]` + "\n"},
		{"localhost:9000/search?q=lobster", `["🦞"]` + "\n"},
	} {
		out, err := exec.Command("curl", test.url).Output()
		if err != nil {
			t.Fatalf("curl %s: %v", test.url, err)
		}
		if string(out) != test.want {
			t.Fatalf("curl %s: got %q, want %q", test.url, string(out), test.want)
		}
	}
}
