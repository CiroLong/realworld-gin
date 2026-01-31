package dto

// User
// {
//  "user": {
//    "email": "jake@example.com",
//    "token": "jwt.token.here",
//    "username": "jake",
//    "bio": "I work at statefarm",
//    "image": null
//  }
//}

type UserResponse struct {
	User UserDTO `json:"user"`
}

type UserDTO struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type ProfileDTO struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}
type ProfileResponse struct {
	Profile ProfileDTO `json:"profile"`
}
