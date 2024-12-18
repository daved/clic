package clic

import "github.com/daved/flagset"

type recOpt struct {
	val   any
	names string
	usage string
	meta  map[string]any
}

type FlagSet struct {
	*flagset.FlagSet
	recOpts []*recOpt
}

func (fs *FlagSet) RecursiveOpt(val any, names, usage string) *flagset.Opt {
	opt := fs.FlagSet.Opt(val, names, usage)

	recOpt := recOpt{
		val:   val,
		names: names,
		usage: usage,
		meta:  opt.Meta,
	}
	fs.recOpts = append(fs.recOpts, &recOpt)

	return opt
}

func applyRecursiveOpts(c *Clic, recOpts []*recOpt) {
	recOpts = append(c.FlagSet.recOpts, recOpts...)

	for _, sub := range c.Subs {
		for _, recOpt := range recOpts {
			subOpt := sub.FlagSet.FlagSet.Opt(recOpt.val, recOpt.names, recOpt.usage)
			for k, v := range recOpt.meta {
				subOpt.Meta[k] = v
			}
		}

		applyRecursiveOpts(sub, recOpts)
	}
}
