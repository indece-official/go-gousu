package broadcaster

type Subscription struct {
	id   int64
	base Base
}

func (s *Subscription) Unsubscribe() {
	s.base.Unsubscribe(s.id)
}
