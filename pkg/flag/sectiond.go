package flag

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

type NamedFlagSets struct {
	Order    []string
	FlagSets map[string]*pflag.FlagSet
}

func NewNamedFlagSets() *NamedFlagSets {
	return &NamedFlagSets{
		Order:    make([]string, 0),
		FlagSets: make(map[string]*pflag.FlagSet),
	}
}

func (nfs *NamedFlagSets) NewFlatSet(name string) *pflag.FlagSet {
	if _, ok := nfs.FlagSets[name]; !ok {
		nfs.Order = append(nfs.Order, name)
		nfs.FlagSets[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
	}

	return nfs.FlagSets[name]
}

func PrintSections(w io.Writer, nfs *NamedFlagSets, cols int) {
	for _, name := range nfs.Order {
		fs := nfs.FlagSets[name]
		if !fs.HasFlags() {
			continue
		}

		wideFS := pflag.NewFlagSet("", pflag.ExitOnError)
		wideFS.AddFlagSet(fs)

		var zzz string
		if cols > 24 {
			zzz = strings.Repeat("z", cols-24)
			wideFS.Int(zzz, 0, strings.Repeat("z", cols-24))
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "\n%s flags:\n\n%s", strings.ToUpper(name[:1])+name[1:], wideFS.FlagUsagesWrapped(cols))

		if cols > 24 {
			i := strings.Index(buf.String(), zzz)
			lines := strings.Split(buf.String()[:i], "\n")
			fmt.Fprint(w, strings.Join(lines[:len(lines)-1], "\n"))
			fmt.Fprintln(w)
		} else {
			fmt.Fprint(w, buf.String())
		}
	}
}
