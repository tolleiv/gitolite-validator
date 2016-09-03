package gitolite_validator

import (
	"fmt"
	"io"
)

type ruleNode struct {
	permission string
	refex      string
	members    []user
}

func (n *ruleNode) addMember(member user) {
	if n.members == nil {
		n.members = make([]user, 0, 8)
	}
	n.members = append(n.members, member)
}

func (n *ruleNode) Type() nodeType {
	return n_rule
}

func (n *ruleNode) String() string {
	return fmt.Sprintf("{rule: %s %s}", n.permission, n.refex)
}

func (n *ruleNode) parseRefex(p *parser) error {
	for {
		t := p.next()
		switch t.t {
		case t_string:
			n.refex = t.s
		case t_assign:
			return nil
		default:
			fmt.Printf("err\n")
			return fmt.Errorf("line %d (1472910546): unexpected token type %v while parsing repo paths",t.l, t.t)
		}
	}
	return nil
}

func (n *ruleNode) parse(p *parser) error {
	if err := n.parseRefex(p); err != nil {
		return err
	}
	for {
		t := p.next()
		switch t.t {
		case t_string:
			u := user{Name:t.s}
			n.addMember(u)
		case t_group:
			u := user{Name:"@" + t.s}
			n.addMember(u)
		case t_eol:
			return nil
		default:
			return fmt.Errorf("line %d (1472910531): unexpected %v token after name, expected = (%s)",t.l,  t.t, n.permission)
		}
	}

	return nil
}

func (n *ruleNode) pretty(w io.Writer, prefix string) error {
	fmt.Fprintf(w, "%srule:\n", prefix)
	fmt.Fprintf(w, "%s%spermission: %s\n", prefix, indent, n.permission)
	if n.refex != "" {
		fmt.Fprintf(w, "%s%srefex: %s\n", prefix, indent, n.refex)
	}
	fmt.Fprintf(w, "%s%smembers:\n", prefix, indent)
	for _, u := range n.members {
		fmt.Fprintf(w, "%s%s%s%s\n", prefix, indent, indent, u.Name)
	}
	return nil
}

func (n *ruleNode) eval(ctx *context) (interface{}, error) {
	return nil, nil
}

