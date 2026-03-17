package types

type AuthenticationResult struct {
	AccessToken  string `json:"accessToken" doc:"JWT access token"`
	RefreshToken string `json:"refreshToken" doc:"JWT refresh token"`
	User         User   `json:"user" doc:"Authenticated user"`
}

type Login struct {
	Identifier string `json:"identifier" validate:"user_identifier,required" doc:"Email, username or phone number"`
	Password   string `json:"password" validate:"password,required" doc:"User password"`
}

type Register struct {
	PhoneNumber string    `json:"phoneNumber" validate:"required,e164" doc:"Phone number"`
	Email       string    `json:"email" validate:"required,email" doc:"Email address"`
	Username    string    `json:"username" validate:"required,username" doc:"Username"`
	Role        UserRoles `json:"role" validate:"required,eq=2|eq=3" doc:"Role (2=Client, 3=Contractor)"`
	FirstName   string    `json:"firstName" validate:"required,min=3,max=255" doc:"First name"`
	LastName    string    `json:"lastName" validate:"required,min=2,max=255" doc:"Last name"`
	Password    string    `json:"password" validate:"required,password" doc:"Password"`
}
