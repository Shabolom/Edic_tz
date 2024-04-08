package tools

import "github.com/gofrs/uuid"

type GetFields interface {
	ID() uuid.UUID
	Permissions() int
}

type GetID interface {
	ID() uuid.UUID
}

type GetPermissions interface {
	Permissions() int
}

type GetLog interface {
	Log() string
}

type GetPassword interface {
	Pass() string
}
