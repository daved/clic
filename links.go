package clic

type Links struct {
	self   *Clic
	parent *Clic
	subs   []*Clic
}

// ResolvedCmd returns the command that was selected during Parse processing.
func (l Links) ResolvedCmd() *Clic {
	return lastCalled(l.self)
}

func (l Links) SubCmds() []*Clic {
	return l.subs
}

func (l Links) ParentCmd() *Clic {
	return l.parent
}
