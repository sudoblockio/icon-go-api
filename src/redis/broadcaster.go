package redis

import (
	"sync"
	"time"
)

// BroadcasterID - type for broadcaster channel IDs
type BroadcasterID uint64

var lastBroadcasterID BroadcasterID = 0

// Broadcaster - Broadcaster channels
type Broadcaster struct {
	InputChannel chan []byte

	// Output
	OutputChannels map[BroadcasterID]chan []byte
}

var broadcasters = map[string]*Broadcaster{}
var broadcastersOnce = map[string]*sync.Once{}

func GetBroadcaster(channelName string) *Broadcaster {
	if _, ok := broadcastersOnce[channelName]; ok == false {
		broadcastersOnce[channelName] = &sync.Once{}
	}
	broadcastersOnce[channelName].Do(
		func() {
			broadcasters[channelName] = &Broadcaster{
				InputChannel:   make(chan []byte),
				OutputChannels: make(map[BroadcasterID]chan []byte),
			}

			broadcasters[channelName].Start()
		},
	)

	return broadcasters[channelName]
}

// AddBroadcastChannel - add channel to  broadcaster
func (b *Broadcaster) AddBroadcastChannel(channel chan []byte) BroadcasterID {

	id := lastBroadcasterID
	lastBroadcasterID++

	b.OutputChannels[id] = channel

	return id
}

// RemoveBroadcastChannelnel - remove channel from broadcaster
func (b *Broadcaster) RemoveBroadcastChannel(id BroadcasterID) {

	_, ok := b.OutputChannels[id]
	if ok {
		delete(b.OutputChannels, id)
	}
}

// Start - Start broadcaster go routine
func (b *Broadcaster) Start() {
	go func() {
		for {
			msg := <-b.InputChannel

			for id, channel := range b.OutputChannels {
				select {
				case channel <- msg:
				case <-time.After(time.Second * 1):
					b.RemoveBroadcastChannel(id)
				}
			}
		}
	}()
}
