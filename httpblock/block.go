package httpblock

import (
	"fmt"
	"os"
	"strings"
)

const hostsPath = `C:\Windows\System32\drivers\etc\hosts`

type Blocker struct {
	Token string
	Sites []string
}

func (b *Blocker) AddBlockEntries() error {
	input, err := os.ReadFile(hostsPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	newline := "\n"
	if strings.Contains(string(input), "\r\n") {
		newline = "\r\n"
	}

	marker := "# " + b.Token // add space

	existing := make(map[string]bool)
	for _, line := range strings.Split(string(input), "\n") {
		l := strings.TrimRight(line, "\r")
		if strings.Contains(l, marker) {
			existing[l] = true
		}
	}

	var toAppend []string
	for _, site := range b.Sites {
		entry := fmt.Sprintf("127.0.0.1\t%s\t%s", site, marker)
		if !existing[entry] {
			toAppend = append(toAppend, entry)
		}
	}

	if len(toAppend) == 0 {
		return nil
	}

	f, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(input) > 0 && !strings.HasSuffix(string(input), "\n") && !strings.HasSuffix(string(input), "\r\n") {
		if _, err := f.WriteString(newline); err != nil {
			return err
		}
	}

	for _, line := range toAppend {
		if _, err := f.WriteString(line + newline); err != nil {
			return err
		}
	}
	return nil
}

func (b *Blocker) RemoveBlockEntries() error {
	input, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	marker := "# " + b.Token
	lines := strings.Split(string(input), "\n")
	out := make([]string, 0, len(lines))

	for _, line := range lines {
		l := strings.TrimRight(line, "\r")
		if strings.Contains(l, marker) {
			continue
		}
		out = append(out, l)
	}

	output := strings.Join(out, "\r\n")
	if !strings.HasSuffix(output, "\r\n") {
		output += "\r\n"
	}
	return os.WriteFile(hostsPath, []byte(output), 0644)
}
