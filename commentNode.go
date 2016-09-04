package gitolite_validator

import (
	"fmt"
	"io"
	"bufio"
	"strings"
)

type commentNode struct {
	body string
}

func (n *commentNode) Type() nodeType {
	return n_comment
}

func (n *commentNode) String() string {
	return fmt.Sprintf("{comment: %s}", n.body)
}

func (n *commentNode) parse(p *parser) error {
	return nil
}

func (n *commentNode) pretty(w io.Writer, prefix string) error {
	fmt.Fprintf(w, "%scomment:\n", prefix)
	r := bufio.NewReader(strings.NewReader(n.body))
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			if line != "" {
				fmt.Fprintf(w, "%s%s%s\n", prefix, indent, line)
			}
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%s%s%s\n", prefix, indent, line)
	}
	return nil
}
