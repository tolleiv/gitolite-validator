package gitolite_validator

import (
	"fmt"
	"io"
	"bytes"
)

type rootNode struct {
	children []node
}

func newRootNode() node {
	return &rootNode{children: make([]node, 0, 8)}
}

func (n *rootNode) Type() nodeType {
	return n_root
}

func (n *rootNode) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for _, child := range n.children {
		fmt.Fprintf(&buf, "%s, ", child)
	}
	if buf.Len() > 1 {
		buf.Truncate(buf.Len() - 2)
	}
	buf.WriteString("}")
	return buf.String()
}

func (n *rootNode) addChild(child node) {
	if n.children == nil {
		n.children = make([]node, 0, 8)
	}
	n.children = append(n.children, child)
}

func (n *rootNode) parse(p *parser) error {
	for {
		t := p.next()
		switch t.t {
		case t_error:
			return fmt.Errorf("line %d (1472910533): saw lex error while parsing root node: %v",t.l, t)
		case t_eof:
			return nil
		case t_eol:
			continue
		case t_string:
			if t.s == "repo" {
				nn := &repoNode{}
				if err := nn.parse(p); err != nil {
					return err
				}
				n.addChild(nn)
			} else {

				p.next()
				return fmt.Errorf("line %d (1472910535): unexpected token type %v while parsing root node (%s)", t.l, t.t, t.s)
			}
		case t_group:
			nn := &groupNode{name: t.s}
			if err := nn.parse(p); err != nil {
				return err
			}
			n.addChild(nn)
		default:
			return fmt.Errorf("line %d (1472910534): unexpected token type %v while parsing root node",t.l, t.t)
		}
	}
}

func (n *rootNode) pretty(w io.Writer, prefix string) error {
	fmt.Fprintf(w, "%sroot:\n", prefix)
	for _, child := range n.children {
		if err := child.pretty(w, prefix + indent); err != nil {
			return err
		}
	}
	return nil
}
