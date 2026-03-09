package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// useTemp redirects the cache directory to t.TempDir() for the duration of the test.
func useTemp(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("PODWISE_CACHE_DIR", dir)
	return dir
}

// ----- Dir ----------------------------------------------------------------

func TestDir_Default(t *testing.T) {
	t.Setenv("PODWISE_CACHE_DIR", "")
	got, err := Dir()
	if err != nil {
		t.Fatalf("Dir() error: %v", err)
	}
	if got == "" {
		t.Fatal("Dir() returned empty string")
	}
}

func TestDir_Override(t *testing.T) {
	want := t.TempDir()
	t.Setenv("PODWISE_CACHE_DIR", want)
	got, err := Dir()
	if err != nil {
		t.Fatalf("Dir() error: %v", err)
	}
	if got != want {
		t.Errorf("Dir() = %q, want %q", got, want)
	}
}

// ----- filePath -----------------------------------------------------------

func TestFilePath(t *testing.T) {
	dir := useTemp(t)
	got, err := filePath(123, "transcript")
	if err != nil {
		t.Fatalf("filePath() error: %v", err)
	}
	want := filepath.Join(dir, "123_transcript.json")
	if got != want {
		t.Errorf("filePath() = %q, want %q", got, want)
	}
}

// ----- Read (miss) --------------------------------------------------------

func TestRead_Miss(t *testing.T) {
	useTemp(t)
	var out any
	hit, err := Read(9999, "transcript", &out)
	if err != nil {
		t.Fatalf("unexpected error on cache miss: %v", err)
	}
	if hit {
		t.Fatal("expected cache miss, got hit")
	}
}

// ----- Write + Read (round-trip) ------------------------------------------

type testPayload struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TestWriteRead_RoundTrip(t *testing.T) {
	useTemp(t)

	in := testPayload{Name: "hello", Count: 42}
	if err := Write(1, "summary", in); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	var out testPayload
	hit, err := Read(1, "summary", &out)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if !hit {
		t.Fatal("expected cache hit, got miss")
	}
	if out != in {
		t.Errorf("Read() = %+v, want %+v", out, in)
	}
}

func TestWriteRead_SliceRoundTrip(t *testing.T) {
	useTemp(t)

	type item struct{ V int }
	in := []item{{1}, {2}, {3}}
	if err := Write(2, "transcript", in); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	var out []item
	hit, err := Read(2, "transcript", &out)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if !hit {
		t.Fatal("expected cache hit")
	}
	if len(out) != len(in) {
		t.Fatalf("len = %d, want %d", len(out), len(in))
	}
	for i := range in {
		if out[i] != in[i] {
			t.Errorf("out[%d] = %v, want %v", i, out[i], in[i])
		}
	}
}

// ----- Write creates missing directories ----------------------------------

func TestWrite_CreatesDirs(t *testing.T) {
	dir := useTemp(t)
	subdir := filepath.Join(dir, "nested")
	t.Setenv("PODWISE_CACHE_DIR", subdir) // deeper than what os.MkdirAll sees

	if err := Write(3, "outline", map[string]int{"x": 1}); err != nil {
		t.Fatalf("Write() should create dirs: %v", err)
	}

	if _, err := os.Stat(subdir); err != nil {
		t.Fatalf("directory was not created: %v", err)
	}
}

// ----- Read: corrupt JSON -------------------------------------------------

func TestRead_CorruptJSON(t *testing.T) {
	dir := useTemp(t)
	p := filepath.Join(dir, "5_mindmap.json")
	if err := os.WriteFile(p, []byte("{not valid json"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out any
	hit, err := Read(5, "mindmap", &out)
	if hit {
		t.Error("expected hit=false for corrupt JSON")
	}
	if err == nil {
		t.Error("expected error for corrupt JSON, got nil")
	}
}

// ----- Stat ---------------------------------------------------------------

func TestStat_Miss(t *testing.T) {
	useTemp(t)
	_, exists, err := Stat(7, "qa")
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if exists {
		t.Fatal("expected exists=false for missing file")
	}
}

func TestStat_Hit(t *testing.T) {
	useTemp(t)

	before := time.Now().Truncate(time.Second)
	if err := Write(8, "qa", "data"); err != nil {
		t.Fatal(err)
	}
	after := time.Now().Add(time.Second)

	modTime, exists, err := Stat(8, "qa")
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if !exists {
		t.Fatal("expected exists=true after Write()")
	}
	if modTime.Before(before) || modTime.After(after) {
		t.Errorf("modTime %v not in expected range [%v, %v]", modTime, before, after)
	}
}

// ----- Key isolation: different (seq, type) pairs don't collide -----------

func TestIsolation(t *testing.T) {
	useTemp(t)

	type val struct{ X int }
	if err := Write(10, "summary", val{X: 1}); err != nil {
		t.Fatal(err)
	}
	if err := Write(10, "transcript", val{X: 2}); err != nil {
		t.Fatal(err)
	}
	if err := Write(11, "summary", val{X: 3}); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		seq         int
		contentType string
		wantX       int
	}{
		{10, "summary", 1},
		{10, "transcript", 2},
		{11, "summary", 3},
	}
	for _, tc := range cases {
		var out val
		hit, err := Read(tc.seq, tc.contentType, &out)
		if err != nil || !hit {
			t.Errorf("Read(%d, %q): hit=%v err=%v", tc.seq, tc.contentType, hit, err)
			continue
		}
		if out.X != tc.wantX {
			t.Errorf("Read(%d, %q).X = %d, want %d", tc.seq, tc.contentType, out.X, tc.wantX)
		}
	}
}
