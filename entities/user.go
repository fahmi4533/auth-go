package entities

type User struct {
	Id          int64
	NamaLengkap string `validate:"required" label:"nama lengkap"`
	Email       string `validate:"required,email,isunique=user2-email"`
	Username    string `validate:"required,gte=3,isunique=user2-username"`
	Password    string `validate:"required,gte=6"`
	Cpassword   string `validate:"required,eqfield=Password" label:"konfirmasi password woy"`
}
