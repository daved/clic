package flagset

import "github.com/daved/flagset"

type Flag = flagset.Flag

type recFlag struct {
	flag  *flagset.Flag
	val   any
	names string
	desc  string
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

func (fs *FlagSet) FlagRecursive(val any, names, desc string) *flagset.Flag {
	flag := fs.FlagSet.Flag(val, names, desc)

	rFlag := recFlag{
		flag:  flag,
		val:   val,
		names: names,
		desc:  desc,
		meta:  flag.Meta,
	}
	fs.recFlags = append(fs.recFlags, &rFlag)

	return flag
}

func ApplyRecursiveFlags(dstFS, srcFS *FlagSet) {
	for _, recFlag := range srcFS.recFlags {
		flag := dstFS.Flag(recFlag.val, recFlag.names, recFlag.desc)
		flag.TypeName = recFlag.flag.TypeName
		flag.DefaultText = recFlag.flag.DefaultText

		for k, v := range recFlag.meta {
			flag.Meta[k] = v
		}
	}
}
