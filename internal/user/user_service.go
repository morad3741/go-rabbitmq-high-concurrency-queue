package user

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) AddUser(name, email, password string) error {
	newUser := User{Name: name, Email: email, Password: password}
	return us.repo.Save(&newUser)
}
