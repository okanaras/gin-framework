package api

type APIErrorResponse struct {
	Message string              `json:"message"`          // struct tag
	Errors  map[string][]string `json:"errors,omitempty"` // omitempty: alan boşsa JSON çıktısında yer almaz

	// Example:
	// {
	//   "message": "Validation Failed",
	//   "errors": {
	//     "email": ["Email is required", "Email must be valid"],
	//     "password": ["Password is required"]
	//   }
	// }
}

type APISuccessResponse struct {
	Message string      `json:"messsage"`
	Data    interface{} `json:"data,omitempty"`
}
