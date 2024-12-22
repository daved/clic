package flagset

import "github.com/daved/flagset"

type Flag = flagset.Flag

type recFlag struct {
	val   any
	names string
	usage string
	meta  map[string]any
}

type FlagSet struct {
	*flagset.FlagSet
	recFlags []*recFlag
}

func New(name string) *FlagSet {
	return &FlagSet{
		FlagSet: flagset.New(name),
	}
}

func (fs *FlagSet) FlagRecursive(val any, names, usage string) *flagset.Flag {
	flag := fs.FlagSet.Flag(val, names, usage)

	rFlag := recFlag{
		val:   val,
		names: names,
		usage: usage,
		meta:  flag.Meta,
	}
	fs.recFlags = append(fs.recFlags, &rFlag)

	return flag
}

func ApplyRecursiveFlags(dstFS, srcFS *FlagSet) {
	for _, recOpt := range srcFS.recFlags {
		flag := dstFS.Flag(recOpt.val, recOpt.names, recOpt.usage)
		for k, v := range recOpt.meta {
			flag.Meta[k] = v
		}
	}
}
