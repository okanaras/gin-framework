package user

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"` // Kullanimi : `` arasina binding kurallari yazilir
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required"`
}
