package log

import (
	"os"
	"testing"
)

func TestLevel(t *testing.T) {
	SetLevel(ErrorLevel)
	if infoLog.Writer() == os.Stdout || errorLog.Writer() != os.Stdout {
		t.Fatal("failed to set log level")
	}

	SetLevel(InfoLevel)
	if infoLog.Writer() != os.Stdout || errorLog.Writer() != os.Stdout {
		t.Fatal("failed to set log level")
	}
}
