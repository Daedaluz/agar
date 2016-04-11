package intern

import "log"

type Broadcaster struct {
	reader chan interface{}
	registered []chan interface{}
}

func NewBroadcaster() *Broadcaster {
	b := &Broadcaster{
		reader: make(chan interface{}, 10),
		registered: make([]chan interface{},0, 5),
	}

	go func() {
		for tosend := range b.reader {
			for _, reg := range b.registered {
				reg <- tosend
			}
		}
		for x, _ := range b.registered {
			close(b.registered[x])
		}
		log.Println("closed!")
	}()

	return b
}

func (b *Broadcaster) Send(data ...interface{}) {
	for _, v := range data {
		b.reader <- v
	}
}

func (b *Broadcaster) Register() chan interface{} {
	res := make(chan interface{})
	b.registered = append(b.registered, res)
	return res
}

func (b *Broadcaster) Close() {
	close(b.reader)
}
