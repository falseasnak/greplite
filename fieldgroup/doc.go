// Package fieldgroup provides a pipeline stage that merges multiple record
// fields into a single composite field by concatenating their string
// representations with a configurable separator.
//
// # Usage
//
// Create a Grouper with New, specifying the destination field name, an ordered
// list of source field names, and the separator string:
//
//	g, err := fieldgroup.New("full_name", []string{"first", "last"}, " ")
//
// Apply it to a parsed log record:
//
//	outRec := g.Apply(inRec)
//
// Source fields that are absent from the record are treated as empty strings,
// so the separator positions are always preserved.
//
// # CLI flags
//
// Register the flags with RegisterFlags and construct via FromFlags:
//
//	--group-dest   destination field name
//	--group-src    comma-separated ordered source field names
//	--group-sep    separator (default: space)
package fieldgroup
