package mapify

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	testUserBob   = TestUser{id: 1, name: "bob"}
	testUserAlice = TestUser{id: 2, name: "alice"}
	testUserFred  = TestUser{id: 3, name: "fred"}
	testUserMarie = TestUser{id: 4, name: "marie"}

	testGroupBoys  = TestGroup{id: 1001, name: "boys", members: []TestInterface{&testUserBob, &testUserFred}}
	testGroupGirls = TestGroup{id: 1002, name: "girls", members: []TestInterface{&testUserAlice, &testUserMarie}}
)

// TestFromSlice verifies the correct behaviour of the FromSlice function when
// mapifying a slice of structs.
func TestFromSlice_OfStructs(t *testing.T) {
	type TestStruct struct {
		id    string
		value any
	}

	standardKeyFn := func(s TestStruct) string {
		return s.id
	}

	testStructOne := TestStruct{id: "one", value: 1}
	testStructTwo := TestStruct{id: "two", value: 2}
	testStructThree := TestStruct{id: "three", value: 3}

	for _, tc := range []struct {
		name     string
		input    []TestStruct
		key      func(TestStruct) string
		expected map[string]TestStruct
	}{
		{
			name:     "nil-input-slice",
			input:    nil,
			key:      standardKeyFn,
			expected: map[string]TestStruct{},
		},
		{
			name:     "empty-struct-slice",
			input:    []TestStruct{},
			key:      standardKeyFn,
			expected: map[string]TestStruct{},
		},
		{
			name: "simple-struct-slice",
			input: []TestStruct{
				testStructOne,
				testStructTwo,
				testStructThree,
			},
			key: standardKeyFn,
			expected: map[string]TestStruct{
				"one":   testStructOne,
				"two":   testStructTwo,
				"three": testStructThree,
			},
		},
		{
			name: "custom-key-func",
			input: []TestStruct{
				testStructOne,
				testStructTwo,
				testStructThree,
			},
			key: func(s TestStruct) string {
				return strings.ToUpper(s.id)
			},
			expected: map[string]TestStruct{
				"ONE":   testStructOne,
				"TWO":   testStructTwo,
				"THREE": testStructThree,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := FromSlice(tc.input, tc.key)
			verifyResult(t, tc.expected, result)
		})
	}
}

// TestFromSlice verifies the correct behaviour of the FromSlice function when
// mapifying a slice of struct pointers.
func TestFromSlice_OfStructPointers(t *testing.T) {
	type TestStruct struct {
		id    string
		value any
	}

	standardKey := func(s *TestStruct) string {
		if s == nil {
			return "(nil)"
		}

		return s.id
	}

	testStructOne := TestStruct{id: "one", value: 1}
	testStructTwo := TestStruct{id: "two", value: 2}
	testStructThree := TestStruct{id: "three", value: 3}

	for _, tc := range []struct {
		name     string
		input    []*TestStruct
		key      func(*TestStruct) string
		expected map[string]*TestStruct
	}{
		{
			name:     "nil-input-slice",
			input:    nil,
			key:      standardKey,
			expected: map[string]*TestStruct{},
		},
		{
			name:     "empty-struct-slice",
			input:    []*TestStruct{},
			key:      standardKey,
			expected: map[string]*TestStruct{},
		},
		{
			name: "simple-struct-slice",
			input: []*TestStruct{
				&testStructOne,
				&testStructTwo,
				&testStructThree,
			},
			key: standardKey,
			expected: map[string]*TestStruct{
				"one":   &testStructOne,
				"two":   &testStructTwo,
				"three": &testStructThree,
			},
		},
		{
			name: "custom-key-func",
			input: []*TestStruct{
				&testStructOne,
				&testStructTwo,
				&testStructThree,
			},
			key: func(s *TestStruct) string {
				return strings.ToUpper(s.id)
			},
			expected: map[string]*TestStruct{
				"ONE":   &testStructOne,
				"TWO":   &testStructTwo,
				"THREE": &testStructThree,
			},
		},
		{
			name: "slice-with-nils-in-it",
			input: []*TestStruct{
				&testStructOne,
				nil,
				&testStructTwo,
				nil,
				&testStructThree,
			},
			key: standardKey,
			expected: map[string]*TestStruct{
				"one":   &testStructOne,
				"two":   &testStructTwo,
				"three": &testStructThree,
				"(nil)": nil,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := FromSlice(tc.input, tc.key)
			verifyResult(t, tc.expected, result)
		})
	}
}

// TestInterface is a type that provides an ID method that returns a string
// value.
type TestInterface interface {
	ID() string
}

// TestUser is a struct that implements the TestInterface type.
type TestUser struct {
	id   int
	name string
}

// ID returns a value consisting of the string `user-` and the text
// representation receiver's id field value unless the receiver is nil, in which
// case the string `(nil)` is appended.
func (u *TestUser) ID() string {
	if u == nil {
		return "user-(nil)"
	}

	return "user-" + strconv.Itoa(u.id)
}

// TestGroup is a struct to also implements the TestInterface type to allow
// testing FromSlice with slices of interfaces instead of structs.
type TestGroup struct {
	id      int
	name    string
	members []TestInterface
}

// ID works in a similar way to (*TestUser).ID except that it uses the prefix
// `group-` instead of `user-`.
func (g *TestGroup) ID() string {
	if g == nil {
		return "group-(nil)"
	}

	return "group-" + strconv.Itoa(g.id)
}

// CompoundKey is a struct used to test FromSlice with a key function that
// returns complex key values.
type CompoundKey struct {
	id        string
	timestamp time.Time
}

// TestFromSlice_WithInterface verifies that FromSlice correctly constructs a
// map from a slice of interfaces rather than structs or struct pointers.
func TestFromSlice_WithInterface(t *testing.T) {

	standardKey := func(i TestInterface) string {
		return i.ID()
	}

	for _, tc := range []struct {
		name     string
		input    []TestInterface
		expected map[string]TestInterface
	}{
		{
			name: "only-users",
			input: []TestInterface{
				&testUserBob,
				&testUserAlice,
				&testUserFred,
				&testUserMarie,
			},
			expected: map[string]TestInterface{
				"user-1": &testUserBob,
				"user-2": &testUserAlice,
				"user-3": &testUserFred,
				"user-4": &testUserMarie,
			},
		},
		{
			name: "only-groups",
			input: []TestInterface{
				&testGroupBoys,
				&testGroupGirls,
			},
			expected: map[string]TestInterface{
				"group-1001": &testGroupBoys,
				"group-1002": &testGroupGirls,
			},
		},
		{
			name: "user-group-mix",
			input: []TestInterface{
				&testUserAlice,
				&testUserMarie,
				&testGroupGirls,
			},
			expected: map[string]TestInterface{
				"group-1002": &testGroupGirls,
				"user-2":     &testUserAlice,
				"user-4":     &testUserMarie,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := FromSlice[TestInterface, string](tc.input, standardKey)
			verifyResult(t, tc.expected, result)
		})
	}
}

// TestFromSlice_WithCustomKeyFn verifies that FromSlice correctly constructs a
// map from a slice of struct pointers and a custom key function.
func TestFromSlice_WithCustomKeyFn(t *testing.T) {
	testTime := time.Now()

	for _, tc := range []struct {
		name     string
		input    []*TestUser
		key      func(u *TestUser) CompoundKey
		expected map[CompoundKey]*TestUser
	}{
		{
			name: "with-custom-key",
			input: []*TestUser{
				&testUserBob,
				&testUserAlice,
				&testUserFred,
				&testUserMarie,
			},
			key: func(u *TestUser) CompoundKey {
				return CompoundKey{
					id:        strconv.Itoa(u.id),
					timestamp: testTime,
				}
			},
			expected: map[CompoundKey]*TestUser{
				{id: "1", timestamp: testTime}: &testUserBob,
				{id: "2", timestamp: testTime}: &testUserAlice,
				{id: "3", timestamp: testTime}: &testUserFred,
				{id: "4", timestamp: testTime}: &testUserMarie,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := FromSlice(tc.input, tc.key)
			verifyResult(t, tc.expected, result)
		})
	}
}

// verifyResult is a convenience function to verify that an expected map of keys
// K to values V matches the actual map of the same type.
func verifyResult[V, K comparable](t *testing.T, expected, actual map[K]V) {
	if actual == nil {
		t.Log("actual is expected to be not nil")
		t.Fail()
	}

	if len(actual) != len(expected) {
		t.Logf("length of actual is expected to be %d but was %d", len(expected), len(actual))
		t.Fail()
	}

	for k, v := range expected {
		if av, ok := actual[k]; !ok {
			t.Logf("actual is expected to contain %v for key %v but it did not have that key", v, k)
			t.Fail()
		} else if av != v {
			t.Logf("actual is expected to contain %v for key %v but contained %v", v, k, av)
			t.Fail()
		}
	}
}

// TestFromSliceWithDuplicates verifies that the FromSliceWithDuplicates
// function correctly creates a map of keys to slices of elements, where the
// key function does not returns unique values for each element.
func TestFromSliceWithDuplicates(t *testing.T) {
	testUserAlanna := TestUser{
		id:   2,
		name: "Alanna",
	}
	testUserMelanie := TestUser{
		id:   4,
		name: "Melanie",
	}

	standardKeyFn := func(u *TestUser) string {
		if u == nil {
			return "3"
		}

		return strconv.Itoa(u.id)
	}

	for _, tc := range []struct {
		name     string
		input    []*TestUser
		expected map[string][]*TestUser
	}{
		{
			name: "with-duplicates",
			input: []*TestUser{
				&testUserBob,
				&testUserAlice,
				&testUserFred,
				&testUserMarie,
				&testUserAlanna,
				&testUserMelanie,
			},
			expected: map[string][]*TestUser{
				"1": {&testUserBob},
				"2": {&testUserAlanna, &testUserAlice},
				"3": {&testUserFred},
				"4": {&testUserMarie, &testUserMelanie},
			},
		},
		{
			name: "with-duplicates-that-are-nil",
			input: []*TestUser{
				&testUserBob,
				&testUserAlice,
				&testUserFred,
				nil,
				&testUserMarie,
				nil,
				&testUserAlanna,
				&testUserMelanie,
			},
			expected: map[string][]*TestUser{
				"1": {&testUserBob},
				"2": {&testUserAlanna, &testUserAlice},
				"3": {&testUserFred, nil, nil},
				"4": {&testUserMarie, &testUserMelanie},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := FromSliceWithDuplicates(tc.input, standardKeyFn)
			verifyResultDuplicates(t, tc.expected, result)
		})
	}
}

// verifyResultDuplicates is a convenience function to verify that an expected
// map of keys K to values of slices of elements E matches an actual map of the
// same type.
func verifyResultDuplicates[E, K comparable](t *testing.T, expected, actual map[K][]E) {
	if actual == nil {
		t.Log("actual expected to not be nil")
		t.Fail()
	}

	if len(actual) != len(expected) {
		t.Logf("length of actual expected to be %d but was %d", len(expected), len(actual))
		t.Fail()
	}

	for emk, emv := range expected {
		if emv == nil {
			t.Fatalf("expected has nil slice mapped to key %v", emk)
		}

		if amv, ok := actual[emk]; !ok {
			t.Logf("actual is expected to contain the key %v", emk)
			t.Fail()
		} else {
			if amv == nil && emv != nil {
				t.Logf("actual is expected to contain a slice for key %v, but is nil", emk)
				t.Fail()
			}

			if len(amv) != len(emv) {
				t.Logf("actual slice for key %v has length %d but is expected to have %d", emk, len(amv), len(emv))
				t.Fail()
			}

			for _, emvE := range emv {
				found := false
				for _, amvE := range amv {
					if amvE == emvE {
						found = true
						break
					}
				}

				if !found {
					t.Logf("actual slice for key %v is expected to contain value %v", emk, emvE)
				}
			}
		}
	}
}
