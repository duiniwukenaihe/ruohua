package resp

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
)

// 自定义 shouldbing error 信息
// example: [field: required]
func validatorForTag(structBody interface{}, err validator.ValidationErrors) map[string]string {
	result := map[string]string{}
	bodyType := reflect.TypeOf(structBody)

	for _, v := range err {
		if field, ok := v.(validator.FieldError); ok {
			structField, _ := bodyType.FieldByName(field.Field())
			result[structField.Tag.Get("json")] = field.Tag()
		}
	}
	return result
}

// 自定义 json 字段值类型错误提示信息
func unmarshalErrorMsg(err *json.UnmarshalTypeError) map[string]string {
	result := map[string]string{}
	result[err.Field] = fmt.Sprintf("the type should be: %v, but received: %s", err.Type, err.Value)
	return result
}

// error 类型判断，返回不同的提示信息
func GetBindingError(ctx *gin.Context, bindBody interface{}, err error) {
	switch err.(type) {
	case validator.ValidationErrors:
		errResult := validatorForTag(bindBody, err.(validator.ValidationErrors))
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": errResult})
	case *json.UnmarshalTypeError:
		errResult := unmarshalErrorMsg(err.(*json.UnmarshalTypeError))
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": errResult})
	case *json.SyntaxError:
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": "Please fill in the correct json parameters"})
	}
}
