package libraries

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/makcik45/jwt-go/config"
)

type Validation struct {
	conn *sql.DB
}

func NewValid() *Validation {
	conn, err := config.DBconn()

	if err != nil {
		panic(err)
	}
	return &Validation{
		conn: conn,
	}
}

func (v *Validation) Init() (*validator.Validate, ut.Translator) {
	// memanggil package tanslator
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	// ini meregister translator english
	en_translations.RegisterDefaultTranslations(validate, trans)

	// mengubah label default nya
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := field.Tag.Get("label")
		return name
	})
	// ini untuk mengubah label nya "yang harus diisi"
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} isi dulu coy", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	// untuk melakuakn validasi agar tidak ada 2 email di dalam satu database
	validate.RegisterValidation("isunique", func(fl validator.FieldLevel) bool {
		params := fl.Param()
		split_param := strings.Split(params, "-")

		tableName := split_param[0]
		tableField := split_param[1]
		fieldValue := fl.Field().String()

		return v.chekIsUniq(tableName, tableField, fieldValue)
	})

	validate.RegisterTranslation("isunique", trans, func(ut ut.Translator) error {
		return ut.Add("isunique", "{0} sudah digunakan coy", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isunique", fe.Field())
		return t
	})
	return validate, trans
}

func (v Validation) Struct(s interface{}) interface{} {
	validate, trans := v.Init()

	vError := make(map[string]interface{})

	err := validate.Struct(s)

	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			vError[v.StructField()] = v.Translate(trans)
		}
	}
	if len(vError) > 0 {
		return vError
	}
	return nil
}

func (v *Validation) chekIsUniq(fieldName, tableField, fildValue string) bool {

	row, _ := v.conn.Query("select "+tableField+" from "+fieldName+" where "+tableField+" =?", fildValue)
	// makna select email from user2 where =fahmiillahi04@gmail.com

	defer row.Close()

	var res string
	for row.Next() {
		row.Scan(&res)
	}
	// ini di cek di email ami@gmail
	// ami@gmail
	return res != fildValue
}

// func UniqUsername(usernames string) bool {
// 	conn := db.CreateCon()

// 	script := "SELECT username FROM users"

// 	rows, err := conn.Query(script)
// 	helper.PanicErr(err)
// 	defer rows.Close()

// 	for rows.Next() {
// 		var username string
// 		err := rows.Scan(&username)
// 		helper.PanicErr(err)
// 		if username == usernames {
// 			return false
// 		}
// 	}

// 	return true
// }
