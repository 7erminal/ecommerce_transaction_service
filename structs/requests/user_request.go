package requests

type GetUserRequest struct {
	UserId string
}

type GetUserWithUsernameRequest struct {
	Username string
}

type GetCustomerRequest struct {
	CustomerId string
}

type GetBranchRequest struct {
	BranchId string
}

type GetCustomerWithUsernameRequest struct {
	Username string
}

type UpdateUserPasswordRequest struct {
	UserId      int64  `validate:"required"`
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required"`
}
