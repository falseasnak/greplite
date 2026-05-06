package fieldmatch

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// RegisterFlags adds fieldmatch CLI flags to the provided flag set.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.StringArray("field-eq", nil, "field=value exact match (repeatable)")
	fs.StringArray("field-contains", nil, "field=value substring match (repeatable)")
	fs.StringArray("field-regex", nil, "field=pattern regex match (repeatable)")
}

// FromFlags constructs a Matcher from parsed CLI flags.
func FromFlags(fs *pflag.FlagSet) (*Matcher, error) {
	var rules []Rule

	if specs, err := fs.GetStringArray("field-eq"); err == nil {
		for _, s := range specs {
			r, err := parseSpec(s, ModeExact)
			if err != nil {
				return nil, err
			}
			rules = append(rules, r)
		}
	}

	if specs, err := fs.GetStringArray("field-contains"); err == nil {
		for _, s := range specs {
			r, err := parseSpec(s, ModeContains)
			if err != nil {
				return nil, err
			}
			rules = append(rules, r)
		}
	}

	if specs, err := fs.GetStringArray("field-regex"); err == nil {
		for _, s := range specs {
			r, err := parseSpec(s, ModeRegex)
			if err != nil {
				return nil, err
			}
			rules = append(rules, r)
		}
	}

	if len(rules) == 0 {
		return None(), nil
	}
	return New(rules)
}

func parseSpec(spec string, mode Mode) (Rule, error) {
	idx := strings.IndexByte(spec, '=')
	if idx <= 0 {
		return Rule{}, fmt.Errorf("fieldmatch: invalid spec %q, expected field=value", spec)
	}
	return Rule{
		Field:   spec[:idx],
		Mode:    mode,
		Pattern: spec[idx+1:],
	}, nil
}
