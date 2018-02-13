package goyesql

import (
	"testing"
)

func TestMustParseFilePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustParseFile should panic if an error occurs, got '%s'", r)
		}
	}()
	MustParseFile("tests/samples/missing.sql")
}

func TestMustParseFileNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustParseFile should not panic if no error occurs, got '%s'", r)
		}
	}()
	MustParseFile("tests/samples/valid.sql")
}

func TestMustParseBytesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustParseBytes should panic if an error occurs, got '%s'", r)
		}
	}()
	MustParseBytes([]byte("I won't work"))
}

func TestMustParseBytesNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustParseBytes should not panic if an error occurs, got '%s'", r)
		}
	}()
	MustParseBytes([]byte("-- name: byte-me\nSELECT * FROM bytes;"))
}

func TestScanToStruct(t *testing.T) {
	type Q struct {
		Ignore   string
		RawQuery string `query:"multiline"`
	}
	type Q2 struct {
		RawQuery string `query:"does-not-exist"`
	}
	var (
		q  Q
		q2 Q2
	)

	queries := MustParseFile("tests/samples/valid.sql")
	err := ScanToStruct(&q, queries, nil)
	if err != nil {
		t.Errorf("failed to scan raw query to struct: %v", err)
	}

	err = ScanToStruct(&q2, queries, nil)
	if err == nil {
		t.Error("expected to fail at non-existent query 'does-not-exist' but didn't")
	}
}

func BenchmarkMustParseFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustParseFile("tests/samples/valid.sql")
	}
}
