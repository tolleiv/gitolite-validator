package gitolite_validator

import (
	"fmt"
	"io"
)

type groupNode struct {
	name    string
	members []user
}

func (n *groupNode) Type() nodeType {
	return n_group
}

func (n *groupNode) String() string {
	return fmt.Sprintf("{group: %s}", n.name)
}

func (n *groupNode) addMember(member user) {
	if n.members == nil {
		n.members = make([]user, 0, 8)
	}
	n.members = append(n.members, member)
}

func (n *groupNode) parse(p *parser) error {
	t := p.next()
	if t.t != t_assign {
		return fmt.Errorf("line %d (1472910529): unexpected %v token after name, expected =",t.l, t.t)
	}

	for {
		t := p.next()
		switch t.t {
		case t_string:
			u := user{Name:t.s}
			n.addMember(u)
		case t_eol:
			return nil
		default:
			return fmt.Errorf("line %d (1472910528): unexpected %v token after name, expected =", t.l, t.t)
		}
	}

	return nil
}

func (n *groupNode) pretty(w io.Writer, prefix string) error {
	fmt.Fprintf(w, "%sgroup:\n", prefix)
	fmt.Fprintf(w, "%sname: %s\n", prefix + indent, n.name)
	fmt.Fprintf(w, "%smembers:\n", prefix + indent)
	for _, u := range n.members {
		fmt.Fprintf(w, "%s%s\n", prefix + indent + indent, u.Name)
	}
	return nil
}

func (n *groupNode) eval(ctx *context) (interface{}, error) {
	return nil, nil
}

