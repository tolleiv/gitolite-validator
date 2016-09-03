package gitolite_validator

import (
	"fmt"
	"io"
	"strings"
)

type repoNode struct {
	paths   []string
	configs []node
	options []node
	rules   []node
}

func (n *repoNode) addRule(child node) {
	if n.rules == nil {
		n.rules = make([]node, 0, 8)
	}
	n.rules = append(n.rules, child)
}
func (n *repoNode) addPath(path string) {
	if n.paths == nil {
		n.paths = make([]string, 0, 8)
	}
	n.paths = append(n.paths, path)
}
func (n *repoNode) addConfig(config node) {
	if n.configs == nil {
		n.configs = make([]node, 0, 8)
	}
	n.configs = append(n.configs, config)
}
func (n *repoNode) addOption(option node) {
	if n.options == nil {
		n.options = make([]node, 0, 8)
	}
	n.options = append(n.options, option)
}

func (n *repoNode) Type() nodeType {
	return n_repo
}

func (n *repoNode) String() string {
	return fmt.Sprintf("{repo: %s}", strings.Join(n.paths, ", "))
}

func (n *repoNode) parsePaths(p *parser) error {
	for {
		t := p.next()
		switch t.t {
		case t_string:
			n.addPath(t.s)
		case t_group:
			n.addPath("@" + t.s)
		case t_eol:
			return nil
		default:
			return fmt.Errorf("line %d: parse error 1472910536: unexpected token type %v while parsing repo paths",t.l, t.t)
		}
	}
	return nil
}

func (n *repoNode) parse(p *parser) error {
	if err := n.parsePaths(p); err != nil {
		return err
	}
	for {

		t := p.next()
		if t.t == t_string {
			switch t.s {
			case "repo":
				p.unread(t)
				return nil
			case "config":
				nn := &configNode{}
				if err := nn.parse(p); err != nil {
					return err
				}
				n.addConfig(nn)
			case "option":
				nn := &configNode{}
				if err := nn.parse(p); err != nil {
					return err
				}
				n.addOption(nn)
			case "R", "RW", "RW+", "RW+D", "C", "RW+C", "RW+CD", "RWD", "RWC":
				nn := &ruleNode{permission: t.s}
				if err := nn.parse(p); err != nil {
					return err
				}
				n.addRule(nn)
			default:
				return fmt.Errorf("line %d: parse error 1472910546: unexpected %v token, expected (%s)", t.l, t.t, t.s)
			}
		} else if t.t == t_dash {
			nn := &ruleNode{permission: "-"}
			if err := nn.parse(p); err != nil {
				return err
			}
			n.addRule(nn)
		} else if t.t == t_eof {
			return nil
		} else if t.t == t_eol {
			continue
		} else {
			return fmt.Errorf("line %d: parse error 1472910537: unexpected %v token, expected", t.l, t.t)
		}
	}
	return nil
}

func (n *repoNode) pretty(w io.Writer, prefix string) error {
	fmt.Fprintf(w, "%srepo:\n", prefix)
	fmt.Fprintf(w, "%s%spaths:\n", prefix, indent)
	for _, path := range n.paths {
		fmt.Fprintf(w, "%s%s\n", prefix + indent + indent, path)
	}
	fmt.Fprintf(w, "%s%srules:\n", prefix, indent)
	for _, rule := range n.rules {
		if err := rule.pretty(w, prefix + indent + indent); err != nil {
			return err
		}
	}

	if len(n.configs) > 0 {
		fmt.Fprintf(w, "%sconfigs:\n", prefix + indent)
		for _, config := range n.configs {
			if err := config.pretty(w, prefix + indent + indent); err != nil {
				return err
			}
		}
	}

	if len(n.options) > 0 {
		fmt.Fprintf(w, "%s%soptions:\n", prefix, indent)
		for _, option := range n.options {
			if err := option.pretty(w, prefix + indent + indent); err != nil {
				return err
			}
		}
	}
	return nil
}
func (n *repoNode) eval(ctx *context) (interface{}, error) {
	return nil, nil
}
