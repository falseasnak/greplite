package fieldtemplate

import (
	"fmt"

	"github.com/spf13/pflag"
)

// RegisterFlags adds the --template-field and --template-expr flags to
// the provided flag set.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.String("template-field", "", "name of the field to write the rendered template into")
	fs.String("template-expr", "", "Go text/template expression evaluated against each record")
}

// FromFlags constructs an Applier from the registered flags. Returns
// None() when neither flag is set. Returns an error when only one of
// the two flags is provided or when the template fails to parse.
func FromFlags(fs *pflag.FlagSet) (*Applier, error) {
	field, err := fs.GetString("template-field")
	if err != nil {
		return nil, err
	}
	expr, err := fs.GetString("template-expr")
	if err != nil {
		return nil, err
	}

	switch {
	case field == "" && expr == "":
		return None(), nil
	case field == "" && expr != "":
		return nil, fmt.Errorf("--template-field is required when --template-expr is set")
	case field != "" && expr == "":
		return nil, fmt.Errorf("--template-expr is required when --template-field is set")
	default:
		return New(field, expr)
	}
}
