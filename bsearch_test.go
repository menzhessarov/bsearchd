package bsearchd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/menzhessarov/bsearchd"
)

func TestStore_GetIndex_EmptyFile(t *testing.T) {
	t.Parallel()

	const (
		inputFile    = "fixtures/input1.txt"
		conformation = 10
	)

	store := bsearchd.NewStore(inputFile, conformation)

	err := store.Load()
	require.ErrorIs(t, err, bsearchd.ErrEmptyFile)
}

func TestStore_GetIndex_OneLineFile(t *testing.T) {
	t.Parallel()

	const (
		inputFile    = "fixtures/input2.txt"
		conformation = 10
	)

	store := bsearchd.NewStore(inputFile, conformation)

	err := store.Load()
	require.NoError(t, err)

	cases := map[string]struct {
		value    int
		expected bsearchd.Entry
		err      error
	}{
		"non existent value": {
			value: 1000,
			err:   bsearchd.ErrNotFound,
		},
		"exact value": {
			value:    10,
			expected: bsearchd.Entry{Index: 0, Value: 10},
		},
		"lower close value": {
			value: 9,
			err:   bsearchd.ErrNotFound,
		},
		"upper close value": {
			value:    11,
			expected: bsearchd.Entry{Index: 0, Value: 10},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			entry, err := store.GetIndex(tc.value)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, entry)
		})
	}
}

func TestStore_GetIndex_TwoLineFile(t *testing.T) {
	t.Parallel()

	const (
		inputFile    = "fixtures/input3.txt"
		conformation = 10
	)

	store := bsearchd.NewStore(inputFile, conformation)

	err := store.Load()
	require.NoError(t, err)

	cases := map[string]struct {
		value    int
		expected bsearchd.Entry
		err      error
	}{
		"non existent value": {
			value: 1000,
			err:   bsearchd.ErrNotFound,
		},
		"exact value": {
			value:    10,
			expected: bsearchd.Entry{Index: 0, Value: 10},
		},
		"lower close value": {
			value: 9,
			err:   bsearchd.ErrNotFound,
		},
		"upper close value": {
			value:    11,
			expected: bsearchd.Entry{Index: 0, Value: 10},
		},
		"other exact value": {
			value:    100,
			expected: bsearchd.Entry{Index: 1, Value: 100},
		},
		"close value": {
			value:    105,
			expected: bsearchd.Entry{Index: 1, Value: 100},
		},
		"bigger close value": {
			value:    109,
			expected: bsearchd.Entry{Index: 1, Value: 100},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			entry, err := store.GetIndex(tc.value)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, entry)
		})
	}
}

func TestStore_GetIndex_ThreeLineFile(t *testing.T) {
	t.Parallel()

	const (
		inputFile    = "fixtures/input4.txt"
		conformation = 10
	)

	store := bsearchd.NewStore(inputFile, conformation)

	err := store.Load()
	require.NoError(t, err)

	cases := map[string]struct {
		value    int
		expected bsearchd.Entry
		err      error
	}{
		"non existent value": {
			value: 1000,
			err:   bsearchd.ErrNotFound,
		},
		"exact value": {
			value:    10,
			expected: bsearchd.Entry{Index: 0, Value: 10},
		},
		"lower close value": {
			value: 9,
			err:   bsearchd.ErrNotFound,
		},
		"upper close value": {
			value:    11,
			expected: bsearchd.Entry{Index: 0, Value: 10},
		},
		"other exact value": {
			value:    100,
			expected: bsearchd.Entry{Index: 1, Value: 100},
		},
		"close value": {
			value:    105,
			expected: bsearchd.Entry{Index: 1, Value: 100},
		},
		"bigger close value": {
			value:    109,
			expected: bsearchd.Entry{Index: 1, Value: 100},
		},
		"third exact value": {
			value:    200,
			expected: bsearchd.Entry{Index: 2, Value: 200},
		},
		"third lower value": {
			value:    182,
			expected: bsearchd.Entry{Index: 2, Value: 200},
		},
		"third upper value": {
			value:    222,
			expected: bsearchd.Entry{Index: 2, Value: 200},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			entry, err := store.GetIndex(tc.value)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, entry)
		})
	}
}

func TestStore_GetIndex_FullFile(t *testing.T) {
	t.Parallel()

	const (
		inputFile    = "fixtures/input5.txt"
		conformation = 10
	)

	store := bsearchd.NewStore(inputFile, conformation)

	err := store.Load()
	require.NoError(t, err)

	cases := map[string]struct {
		value    int
		expected bsearchd.Entry
		err      error
	}{
		"non existent value": {
			value: 20000000,
			err:   bsearchd.ErrNotFound,
		},
		"zero index": {
			value:    0,
			expected: bsearchd.Entry{Index: 0, Value: 0},
		},
		"zero index: higher value": {
			value: 1,
			err:   bsearchd.ErrNotFound,
		},
		"exact value": {
			value:    2300,
			expected: bsearchd.Entry{Index: 23, Value: 2300},
		},
		"lower value": {
			value:    2299,
			expected: bsearchd.Entry{Index: 22, Value: 2200},
		},
		"upper value": {
			value:    2399,
			expected: bsearchd.Entry{Index: 23, Value: 2300},
		},
		"last index": {
			value:    10000000,
			expected: bsearchd.Entry{Index: 100000, Value: 10000000},
		},
		"last index: lower value": {
			value:    9999999,
			expected: bsearchd.Entry{Index: 99999, Value: 9999900},
		},
		"last index: upper value": {
			value:    11000000,
			expected: bsearchd.Entry{Index: 100000, Value: 10000000},
		},
		"last index: error upper value": {
			value: 12000000,
			err:   bsearchd.ErrNotFound,
		},
		"last index: bigger upper value": {
			value: 20000000,
			err:   bsearchd.ErrNotFound,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			entry, err := store.GetIndex(tc.value)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, entry)
		})
	}
}
