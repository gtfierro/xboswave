package mqpb

//The drops field is unique in different places, copy the message to allow
//local modification of the drops field
func ShallowCloneMessageForDrops(m *Message) *Message {
	rv := &Message{}
	*rv = *m

	rv.Drops = make([]int64, len(m.Drops))
	copy(rv.Drops, m.Drops)
	return rv
}
