package types

type AuthenticationResult struct {
	AccessToken  string `json:"accessToken" doc:"JWT access token"`
	RefreshToken string `json:"refreshToken" doc:"JWT refresh token"`
	User         User   `json:"user" doc:"Authenticated user"`
}

type Login struct {
	Identifier string `json:"identifier" validate:"user_identifier,required" doc:"Email, username or phone number" minLength:"1" maxLength:"255"`
	Password   string `json:"password" validate:"password,required" doc:"User password" minLength:"1" maxLength:"255"`
}

type Register struct {
	PhoneNumber string    `json:"phoneNumber" validate:"required,phone_number" doc:"Phone number" minLength:"1" maxLength:"13"`
	Email       string    `json:"email" validate:"required,email" doc:"Email address" format:"email" maxLength:"255"`
	Username    string    `json:"username" validate:"required,username" doc:"Username" minLength:"1" maxLength:"255"`
	Role        UserRoles `json:"role" validate:"required,eq=2|eq=3" doc:"Role (2=Client, 3=Contractor)" enum:"2,3"`
	FirstName   string    `json:"firstName" validate:"required,min=3,max=255" doc:"First name" minLength:"3" maxLength:"255"`
	LastName    string    `json:"lastName" validate:"required,min=2,max=255" doc:"Last name" minLength:"2" maxLength:"255"`
	Password    string    `json:"password" validate:"required,password" doc:"Password" minLength:"1" maxLength:"255"`
}
