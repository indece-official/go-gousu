package broadcaster

import "sync"

type Generic struct {
	consumers      map[int64]chan interface{}
	nextConsumerID int64
	mutexConsumers sync.Mutex
	lastValue      interface{}
}

var _ Base = (*Generic)(nil)

func (b *Generic) Value() interface{} {
	return b.lastValue
}

func (b *Generic) Next(val interface{}) {
	b.mutexConsumers.Lock()
	defer b.mutexConsumers.Unlock()

	b.lastValue = val

	for _, consumer := range b.consumers {
		consumer <- val
	}
}

func (b *Generic) Subscribe() (chan interface{}, *Subscription) {
	b.mutexConsumers.Lock()
	defer b.mutexConsumers.Unlock()

	id := b.nextConsumerID
	b.nextConsumerID++

	consumer := make(chan interface{}, 1)

	b.consumers[id] = consumer

	subscription := &Subscription{
		id:   id,
		base: b,
	}

	return consumer, subscription
}

func (b *Generic) Unsubscribe(id int64) {
	b.mutexConsumers.Lock()
	defer b.mutexConsumers.Unlock()

	if _, ok := b.consumers[id]; !ok {
		return
	}

	delete(b.consumers, id)
}

func NewGeneric(initialValue interface{}) *Generic {
	return &Generic{
		consumers: map[int64]chan interface{}{},
		lastValue: initialValue,
	}
}
