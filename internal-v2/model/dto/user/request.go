package user

// Registration POST /api/users
//{
//  "user":{
//    "username": "Jacob",
//    "email": "jake@jake.jake",
//    "password": "jakejake"
//  }
//}
// Required fields: email, username, password

type RegisterRequest struct {
	User struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	} `json:"user"`
}

// Authentication POST /api/users/login
//{
//  "user":{
//    "email": "jake@jake.jake",
//    "password": "jakejake"
//  }
//}
// Required fields: email, password

type LoginRequest struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	} `json:"user"`
}

// Update User PUT /api/user
// {
//  "user":{
//    "email": "jake@jake.jake",
//    "bio": "I like to skateboard",
//    "image": "https://i.stack.imgur.com/xHWG8.jpg"
//  }
//}
// Accepted fields: email, username, password, image, bio

type UpdateUserRequest struct {
	User struct { // tips: 这里用指针是因为要区分未传参和传入空字符
		Email    *string `json:"email"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
	} `json:"user"`
}
