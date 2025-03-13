package bsearchd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	ErrNotFound  = errors.New("value not found")
	ErrEmptyFile = errors.New("empty input file")
)

type Store struct {
	values       []int
	inputFile    string
	conformation int
}

func NewStore(inputFile string, conformation int) *Store {
	return &Store{inputFile: inputFile, conformation: conformation}
}

func (s *Store) Load() error {
	f, err := os.Open(s.inputFile)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("convert line: %w", err)
		}

		s.values = append(s.values, v)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("check scanner error: %w", err)
	}

	if len(s.values) == 0 {
		return ErrEmptyFile
	}

	return nil
}

type Entry struct {
	Index int
	Value int
}

func (s *Store) GetIndex(searchValue int) (Entry, error) {
	l, mid, r := 0, 0, len(s.values)-1

	for l <= r {
		mid = l + (r-l)/2

		if s.values[mid] == searchValue {
			return Entry{Index: mid, Value: searchValue}, nil
		} else if s.values[mid] > searchValue {
			r = mid - 1
		} else if s.values[mid] < searchValue {
			l = mid + 1
		}
	}

	lower := mid

	if 0 < mid && s.values[mid] > searchValue {
		lower = mid - 1
	}

	delta := searchValue * 1 / s.conformation

	if delta == 0 {
		return Entry{}, ErrNotFound
	}

	if 0 <= lower && searchValue-delta <= s.values[lower] {
		return Entry{Index: lower, Value: s.values[lower]}, nil
	}

	if lower < len(s.values)-1 && searchValue+delta >= s.values[lower+1] {
		return Entry{Index: lower + 1, Value: s.values[lower+1]}, nil
	}

	return Entry{}, ErrNotFound
}
