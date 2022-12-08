package go_template

import (
	"strings"
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	tm := time.Now()
	s := FormatTime(tm, 5, time.RFC822)
	if !strings.HasSuffix(s, "PKT") {
		t.Error("not PKT")
	}
}
