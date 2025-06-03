package service

import (
	"gorm.io/gorm"
)

type Session struct {
	tx *gorm.DB
}

func NewSession(tx *gorm.DB) *Session {
	return &Session{tx: tx}
}

func (sess *Session) Tx(work func(tx *gorm.DB) error) error {
	return sess.tx.Transaction(work)
}

func (sess *Session) DB() *gorm.DB {
	return sess.tx
}
