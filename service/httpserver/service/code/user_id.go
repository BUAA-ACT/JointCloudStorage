package code

import uuid "github.com/satori/go.uuid"

func GenUserID() uuid.UUID {
	return genUUIDv4()
}
