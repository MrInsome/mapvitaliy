package repository

import (
	"github.com/beanstalkd/go-beanstalk"
	"time"
)

func (r *Repository) NewBeanstalkConn() *Repository {
	conn, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		return nil
	}
	r.conn = conn
	return r
}

func (r *Repository) Close() error {
	return r.conn.Close()
}

func (r *Repository) Put(body []byte, priority uint32, delay, ttr time.Duration) (uint64, error) {
	return r.conn.Put(body, priority, delay, ttr)
}

func (r *Repository) Delete(id uint64) error {
	return r.conn.Delete(id)
}

func (r *Repository) Reserve(ttr time.Duration) (id uint64, body []byte, err error) {
	return r.conn.Reserve(ttr)
}
