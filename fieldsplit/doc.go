// Package fieldsplit implements a pipeline stage that splits a single string
// field into multiple named fields using a configurable separator.
//
// # Usage
//
// Given a log record with a field like:
//
//	addr = "host:port"
//
// fieldsplit can expand it into:
//
//	host = "host"
//	port = "port"
//
// The original field is always preserved in the output record.
//
// # CLI flags
//
//	--split-field   source field name to split
//	--split-into    comma-separated list of destination field names
//	--split-sep     separator string (default ",")
package fieldsplit
