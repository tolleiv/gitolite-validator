package gitolite_validator

import (
	"fmt"
	"io"
)

type configNode struct {
	name  string
	value string
}

func (n *configNode) Type() nodeType {
	return n_config
}

func (n *configNode) String() string {
	return fmt.Sprintf("{config: %s}", n.name)
}

func (n *configNode) parse(p *parser) error {

	tn := p.next()
	if tn.t != t_string {
		return fmt.Errorf("line %d (1472910540): unexpected token type %v ", tn.l, tn.t)
	}
	n.name = tn.s

	ta := p.next()
	if ta.t != t_assign {
		return fmt.Errorf("line %d (1472910541): unexpected token type %v ",ta.l, ta.t)
	}

	for {
		tv := p.next()
		switch tv.t {
		case t_string:
			n.value = n.value + tv.s
		case t_eol:
			return nil
		default:
			return fmt.Errorf("line %d (1472910543): unexpected token type %v with value %s", tv.t, tv.s)
		}
	}
	return nil
}

func (n *configNode) pretty(w io.Writer, prefix string) error {
	fmt.Fprintf(w, "%s%s = %s\n", prefix + indent, n.name, n.value)
	return nil
}
