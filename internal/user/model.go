package user

type User struct {
	ID      string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string   `json:"name" bson:"name"`
	Age     int      `json:"age" bson:"age"`
	Friends []string `json:"friends" bson:"friends"`
}

type CreateUserDTO struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type CreateFriendDTO struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
}

type UpdateUserAgeDTO struct {
	Age int `json:"age"`
}

func NewUser(CreateUserDTO *CreateUserDTO) *User {
	return &User{
		Name:    CreateUserDTO.Name,
		Age:     CreateUserDTO.Age,
		Friends: make([]string, 0),
	}
}
