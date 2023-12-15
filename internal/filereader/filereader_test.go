package filereader

import "testing"

func TestTotalLines(t *testing.T) {
	type Test struct {
		FileName      string
		ExpectedLines int64
	}

	tests := []Test{
		{
			"testdata/20-lines.txt.gz",
			20,
		},
		{
			"testdata/0-lines.txt.gz",
			0,
		},
		{
			"testdata/5-lines-with-missing-line.txt.gz",
			5,
		},
		{
			"testdata/5-lines-with-empty-line.txt.gz",
			5,
		},
		{
			"testdata/10-lines-with-long-line.txt.gz",
			10,
		},
		{
			"testdata/10m-lines.txt.gz",
			10000000,
		},
	}

	for _, test := range tests {
		t.Run(test.FileName, func(t *testing.T) {
			fr, err := New(test.FileName)
			if err != nil {
				t.Fatalf("Failed to open file: %v", err)
			}

			total, err := fr.TotalLines()
			if err != nil {
				t.Fatalf("Failed to get total lines: %v", err)
			}

			if total != test.ExpectedLines {
				t.Errorf("Expected %d lines, got %d", test.ExpectedLines, total)
			}
		})
	}

}
