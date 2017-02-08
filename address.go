package dataset

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// this regex makes sure we have snake_case.addresses.with.dot_separators.1_and_only_alphanumeric_characters
// `^[a-z0-9-_/]+(\.[a-z0-9-_/]+)?$`
var addressRegex = regexp.MustCompile(`^[a-z-_]+[\.[a-z0-9-_/]+]?[a-z0-9-_]$`)

// check for a valid namespce address
func ValidAddressString(s string) bool {
	return addressRegex.MatchString(s)
}

// a address is a string slice that divides the global namspace
type Address []string

// Create a new address from one or more strings. all strings are divided by any dot separators.
// So the internal array would map as:
// 	NewAddress("user.dataset","table") => ["user","dataset","table"]
// Which is the eqivelent to:
// 	NewAddress("user", "dataset", "table") => ["user", "dataset", "table"]
func NewAddress(strs ...string) (p Address) {
	for _, str := range strs {
		for _, s := range strings.Split(str, ".") {
			if s != "" {
				p = append(p, s)
			}
		}
	}

	return
}

// IsEmpty is a convenience method to see if the address is assigned
func (a Address) IsEmpty() bool {
	return len(a) == 0 || a[0] == ""
}

// Conform to stringer interface
func (p Address) String() string {
	return strings.Join(p, ".")
}

func (a Address) PathString(path ...string) string {
	return fmt.Sprintf("/%s/%s", strings.Join(a, "/"), strings.Join(path, "/"))
	// return "/" + strings.Join(a, "/")
}

func (a Address) Equal(b Address) bool {
	if len(a) != len(b) {
		return false
	}

	for i, ds := range a {
		if ds != b[i] {
			return false
		}
	}

	return true
}

func (a Address) Parent() Address {
	if len(a) <= 1 {
		return NewAddress("")
	}

	return a[:len(a)-1]
}

func (a Address) IsParent(b Address) bool {
	if len(a)+1 != len(b) {
		return false
	}

	for i, ds := range a {
		if ds != b[i] {
			return false
		}
	}

	return true
}

func (a Address) IsAncestor(b Address) bool {
	if len(a) >= len(b) {
		return false
	}

	if a.IsEmpty() {
		return true
	}

	for i, ds := range a {
		if b[i] != ds {
			return false
		}
	}

	return true
}

func (p Address) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, p.String())), nil
}

func (ad *Address) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Data Type should be a string, got %s", data)
	}

	*ad = NewAddress(s)
	return nil
}
