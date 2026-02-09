package types

type AuthenticationResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

type Login struct {
	Identifier string `json:"identifier" validate:"user_identifier,required"`
	Password   string `json:"password" validate:"password,required"`
}

type Register struct {
	PhoneNumber string    `json:"phoneNumber" validate:"required,phone_number"`
	Email       string    `json:"email" validate:"required,email"`
	Username    string    `json:"username" validate:"required,username"`
	Role        UserRoles `json:"role" validate:"required,eq=2|eq=3"`
	FirstName   string    `json:"firstName" validate:"required,min=3,max=255"`
	LastName    string    `json:"lastName" validate:"required,min=2,max=255"`
	Password    string    `json:"password" validate:"required,password"`
}
