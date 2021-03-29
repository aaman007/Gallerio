package accounts

type SignUpForm struct {
	Username string `schema:"username"`
	Email string `schema:"email"`
	Password string `schema:"password"`
}
