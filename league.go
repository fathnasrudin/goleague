package poker

type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if name == p.Name {
			return &l[i]
		}
	}
	return nil
}