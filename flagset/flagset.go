package flagset

import "github.com/daved/flagset"

type Flag = flagset.Flag

type recFlag struct {
	flag  *flagset.Flag
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
		flag:  flag,
		val:   val,
		names: names,
		usage: usage,
		meta:  flag.Meta,
	}
	fs.recFlags = append(fs.recFlags, &rFlag)

	return flag
}

func ApplyRecursiveFlags(dstFS, srcFS *FlagSet) {
	for _, recFlag := range srcFS.recFlags {
		flag := dstFS.Flag(recFlag.val, recFlag.names, recFlag.usage)
		flag.TypeHint = recFlag.flag.TypeHint
		flag.DefaultHint = recFlag.flag.DefaultHint

		for k, v := range recFlag.meta {
			flag.Meta[k] = v
		}
	}
}
