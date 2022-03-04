package broadcaster

import "sync"

type Bool struct {
	consumers      map[int64]chan bool
	nextConsumerID int64
	mutexConsumers sync.Mutex
	lastValue      bool
}

var _ Base = (*Bool)(nil)

func (b *Bool) Value() bool {
	return b.lastValue
}

func (b *Bool) Next(val bool) {
	b.mutexConsumers.Lock()
	defer b.mutexConsumers.Unlock()

	b.lastValue = val

	for _, consumer := range b.consumers {
		consumer <- val
	}
}

func (b *Bool) Subscribe() (chan bool, *Subscription) {
	b.mutexConsumers.Lock()
	defer b.mutexConsumers.Unlock()

	id := b.nextConsumerID
	b.nextConsumerID++

	consumer := make(chan bool, 1)

	b.consumers[id] = consumer

	subscription := &Subscription{
		id:   id,
		base: b,
	}

	return consumer, subscription
}

func (b *Bool) Unsubscribe(id int64) {
	b.mutexConsumers.Lock()
	defer b.mutexConsumers.Unlock()

	if _, ok := b.consumers[id]; !ok {
		return
	}

	delete(b.consumers, id)
}

func NewBool(initialValue bool) *Bool {
	return &Bool{
		consumers: map[int64]chan bool{},
		lastValue: initialValue,
	}
}
