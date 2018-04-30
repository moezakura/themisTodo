package module

import (
	"time"
	"sync"
)

type SessionModule struct {
	Session sync.Map
}

type Session struct {
	Uuid     int
	Limit    int64
	LifeTime int64
}

func NewSessionModule() *SessionModule {
	self := &SessionModule{
		sync.Map{},
	}
	go self.background()
	return self
}

func (self *SessionModule) background() {
	for {
		sessionTemp := sync.Map{}
		self.Session.Range(func(key, value interface{}) bool {
			sessionTemp.Store(key, value)
			return true
		})
		now := time.Now().Unix()

		self.Session.Range(func(key, _value interface{}) bool {
			value := _value.(*Session)
			if value.Limit < now {
				self.Session.Delete(key)
			}
			return true
		})

		time.Sleep(1 * time.Second)
	}
}

func (self *SessionModule) Add(key string, uuid int, limit int64) {
	now := time.Now().Unix()
	self.Session.Store(key, &Session{uuid, now + limit, limit})
}

func (self *SessionModule) GetUuid(key string) (isExist bool, uuid int) {
	now := time.Now().Unix()

	if _value, ok := self.Session.Load(key); ok {
		value := _value.(*Session)
		if value.Limit >= now {
			return true, value.Uuid
		} else {
			return false, 0
		}
	}
	return false, 0
}

func (self *SessionModule) UpdateTime(key string) {
	now := time.Now().Unix()

	if _value, ok := self.Session.Load(key); ok {
		value := _value.(*Session)
		if value.Limit >= now {
			value.Limit = now + value.LifeTime
			self.Session.Store(key, value)
		}
	}
}
