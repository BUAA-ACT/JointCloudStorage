package code

import uuid "github.com/satori/go.uuid"

func GenUserId() uuid.UUID {
	return genUUIDv4()
}
