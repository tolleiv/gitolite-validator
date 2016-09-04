package gitolite_validator

import (
	"io"
)

type nodeType int

const (
	n_error nodeType = iota
	n_root
	n_comment
	n_group
	n_repo
	n_rule
	n_config
)

var indent = "  "

type context struct {
	public  map[string]interface{}
	private map[string]interface{}
}

type node interface {
	Type() nodeType
	parse(*parser) error
	pretty(io.Writer, string) error
}

type user struct {
	Name string
}

// -------------------------------------------------------------------------------------------------------------

