package hostman

import "testing"

func TestGetSection_data(t *testing.T) {
	data := `
<START>
data
<END>
`
	res, err := GetSection(data, "<START>", "<END>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "data" {
		t.Fatalf("unexpected result: %q", res)
	}
}

func TestGetSection_nodata(t *testing.T) {
	data := `
<START>

<END>
`
	res, err := GetSection(data, "<STARTB>", "<END>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "" {
		t.Fatalf("unexpected result: %q", res)
	}
}

func TestInsertOrReplaceSection_NormalizesNewlines(t *testing.T) {
	data := ""
	start := "<START>"
	end := "<END>"
	content := ""
	got, err := InsertOrReplaceSection(data, start, end, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<START>\n\n<END>\n"
	if got != want {
		t.Fatalf("unexpected result.\nGot:%q\nWant:%q", got, want)
	}
}

func TestInsertOrReplaceSection_Append(t *testing.T) {
	data := "line1\nline2\n"
	start := "<START>"
	end := "<END>"
	content := "hello\nworld\n"

	got, err := InsertOrReplaceSection(data, start, end, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := data + start + "\n" + content + "\n" + end + "\n"
	if got != want {
		t.Fatalf("unexpected result.\nGot:\n%q\nWant:\n%q", got, want)
	}
}

func TestInsertOrReplaceSection_Replace(t *testing.T) {
	data := "before\n<START>old content\n<END>\nafter\n"
	start := "<START>"
	end := "<END>"
	content := "new content\nmore\n"

	got, err := InsertOrReplaceSection(data, start, end, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "before\n" + start + "\n" + "new content\nmore\n\n" + end + "\n" + "after\n"
	if got != want {
		t.Fatalf("unexpected result.\nGot:\n%q\nWant:\n%q", got, want)
	}
}

func TestInsertOrReplaceSection_ErrorOnMismatchedMarkers(t *testing.T) {
	cases := []struct{ name, data string }{
		{"only start", "foo\n<START>here\nbar"},
		{"only end", "foo\n<END>here\nbar"},
	}
	for _, tc := range cases {
		_, err := InsertOrReplaceSection(tc.data, "<START>", "<END>", "x")
		if err == nil {
			t.Fatalf("%s: expected error, got nil", tc.name)
		}
	}
}
