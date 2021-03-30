package accounts

type SignUpForm struct {
	Name string `schema:"name"`
	Username string `schema:"username"`
	Email string `schema:"email"`
	Password string `schema:"password"`
}

type SignInForm struct {
	Email string `schema:"email"`
	Password string `schema:"password"`
}

