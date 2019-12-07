package models

// Vote Информация о голосовании пользователя.
//
// swagger:model Vote
type Vote struct {
	ID int32
	// Идентификатор пользователя.
	// Required: true
	Nickname string `json:"nickname"`

	// Отданный голос.
	// Required: true
	// Enum: [-1 1]
	Voice    int32 `json:"voice"`
	ThreadId int32
}
