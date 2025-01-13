// Package flagset wraps the [flagset] package.
package flagset

import "github.com/daved/flagset"

// Flag is an alias of [flagset.Flag].
type Flag = flagset.Flag

type recFlag struct {
	flag  *flagset.Flag
	val   any
	names string
	desc  string
	meta  map[string]any
}

// FlagSet wraps [flagset.FlagSet]. This type can store recursive flag info.
type FlagSet struct {
	*flagset.FlagSet
	recFlags []*recFlag
}

// New returns an instance of FlagSet.
func New(name string) *FlagSet {
	return &FlagSet{
		FlagSet: flagset.New(name),
	}
}

// FlagRecursive adds a flag option to the FlagSet. Data necessary for applying
// a flag recursively is stored. It is the responsibility of callers to also
// call [ApplyRecursiveFlags] with the appropriate [FlagSet] instances.
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

// ApplyRecursiveFlags will apply any recursive flags from the source [FlagSet]
// to the destination [FlagSet].
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
