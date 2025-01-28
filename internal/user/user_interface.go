package user

// UserServiceInterface defines the methods that a user service must implement.
type UserServiceInterface interface {
	AddUser(name, email, password string) error
}
