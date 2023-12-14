package filereader

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

type FileReader struct {
	r   *bufio.Reader
	f   *os.File
	gzr *gzip.Reader
}

func New(name string) (*FileReader, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}

	r := bufio.NewReader(gzr)
	if err != nil {
		return nil, fmt.Errorf("failed to create new reader: %w", err)
	}

	return &FileReader{
		r,
		file,
		gzr,
	}, nil
}

func (fr *FileReader) TotalLines() (int64, error) {
	file, err := os.Open(fr.f.Name())
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return 0, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	r := bufio.NewReader(gzr)

	buf := make([]byte, 32*1024)
	var count int64
	lineSep := []byte{'\n'}
	var lastByte byte

	for {
		c, err := r.Read(buf)
		count += int64(bytes.Count(buf[:c], lineSep))

		if c > 0 {
			lastByte = buf[c-1]
		}

		switch {
		case err == io.EOF:
			if lastByte != '\n' && c > 0 {
				count += 1 // Add one if the file does not end with a newline character
			}
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func (fr *FileReader) ReadLines(n int) ([][]byte, error) {
	var lines [][]byte

	for {
		line, err := fr.r.ReadBytes('\n')
		if err == io.EOF {
			if len(line) > 0 {
				lines = append(lines, line)
			}
			break
		}
		if err != nil {
			return lines, err
		}
		lines = append(lines, line)
		if len(lines) >= n {
			break
		}
	}

	return lines, nil
}

func (fs *FileReader) Close() error {
	if err := fs.f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	if err := fs.gzr.Close(); err != nil {
		return fmt.Errorf("failed to close gzip reader: %w", err)
	}

	return nil
}
