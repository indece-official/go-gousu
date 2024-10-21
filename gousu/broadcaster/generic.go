package broadcaster

import "sync"

type Generic[O comparable] struct {
	consumers      map[int64]chan O
	nextConsumerID int64
	mutexConsumers sync.Mutex
	lastValue      O
}

var _ Base = (*Generic[bool])(nil)
var _ Base = (*Generic[error])(nil)
var _ Base = (*Generic[int64])(nil)

func (b *Generic[O]) Value() O {
	return b.lastValue
}

func (b *Generic[O]) Next(val O) {
	b.mutexConsumers.Lock()
	defer b.mutexConsumers.Unlock()

	b.lastValue = val

	for _, consumer := range b.consumers {
		consumer <- val
	}
}

func (b *Generic[O]) Subscribe() (chan O, *Subscription) {
	b.mutexConsumers.Lock()
	defer b.mutexConsumers.Unlock()

	id := b.nextConsumerID
	b.nextConsumerID++

	consumer := make(chan O, 1)

	b.consumers[id] = consumer

	subscription := &Subscription{
		id:   id,
		base: b,
	}

	return consumer, subscription
}

func (b *Generic[O]) Unsubscribe(id int64) {
	b.mutexConsumers.Lock()
	defer b.mutexConsumers.Unlock()

	if _, ok := b.consumers[id]; !ok {
		return
	}

	delete(b.consumers, id)
}

func NewGeneric[O comparable](initialValue O) *Generic[O] {
	return &Generic[O]{
		consumers: map[int64]chan O{},
		lastValue: initialValue,
	}
}
