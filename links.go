package clic

// Links contains Clic instance relationships. The relationships are defined
// here to clean up the Clic type's docs.
type Links struct {
	parent *Clic
	subs   []*Clic
}

// SubCmds returns the child Clic instances associated with the Clic instance.
func (l Links) SubCmds() []*Clic {
	return l.subs
}

// ParentCmd returns the parent Clic instance associated with the Clic instance.
func (l Links) ParentCmd() *Clic {
	return l.parent
}
