package valid

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kms9/publicyc/pkg/server/ogin"
)

type ValidatorErr struct {
	Field string `json:"field"`
	Err   error  `json:"err"`
}

func NewValidatorErr(field string, err error) *ValidatorErr {
	return &ValidatorErr{field, err}
}

func (v *ValidatorErr) Error() string {
	if v == nil || v.Err == nil {
		return ""
	}
	return fmt.Sprintf("filed:%s err:[%s]", v.Field, v.Err.Error())
}

func (v *ValidatorErr) String() string {
	if v == nil || v.Err == nil {
		return ""
	}
	return fmt.Sprintf("filed:%s err:[%s]", v.Field, v.Err.Error())
}

func (v *ValidatorErr) ParamsErr(c *gin.Context) {
	c.JSON(http.StatusBadRequest, ogin.ErrMsg{
		Err:   v.String(),
		Msg:   "params err",
		Debug: "",
	})
}

// Validator 自定义
type Validator struct {
	*validator.Validate
}

func New() *Validator {
	v := validator.New()
	return &Validator{v}
}

// CheckFields 检查多字段
func (v *Validator) CheckFields(err ...*ValidatorErr) *ValidatorErr {
	var fields strings.Builder
	var errs strings.Builder
	var firstErr bool
	for _, validatorErr := range err {
		if validatorErr != nil && validatorErr.Err != nil {
			if firstErr {
				_, _ = fields.WriteString(", \n")
				_, _ = errs.WriteString(", \n")
			}
			_, _ = fields.WriteString(validatorErr.Field)
			_, _ = errs.WriteString(validatorErr.Err.Error())
			firstErr = true
		}
	}
	if errs.String() == "" {
		return nil
	}
	return NewValidatorErr(fields.String(), fmt.Errorf(errs.String()))
}

type EachParams struct {
	Field       string // 必须为 Value.interface 中 key;
	ValidateTag string
}

// fieldToLower 字段转成小写
func (v *Validator) fieldToLower(field string) string {
	return strings.ReplaceAll(strings.ToLower(field), "_", "")
}

// CheckEach 检查每一个类型
func (v *Validator) CheckEach(params []*EachParams, value []interface{}) *ValidatorErr {
	err := make([]*ValidatorErr, 0)
	checkTag := map[string]string{}
	for _, param := range params {
		checkTag[param.Field] = param.ValidateTag
	}

	for i := range value {
		val := structs.Map(value[i])
		for k := range val {
			val[v.fieldToLower(k)] = val[k]
		}

		for _, param := range params {
			filed := v.fieldToLower(param.Field)
			if _, exists := val[filed]; !exists {
				err = append(err, NewValidatorErr(param.Field, fmt.Errorf("column not found")))
			}
			err = append(err, v.Check(param.Field, param.ValidateTag, val[filed]))
		}
	}

	return v.CheckFields(err...)
}

// Check 检查单一字段
func (v *Validator) Check(field string, validateTag string, value interface{}) *ValidatorErr {
	err := v.Validate.Var(value, validateTag)
	if err != nil {
		return NewValidatorErr(field, err)
	}
	return nil
}

// CheckStruct 检查struct
func (v *Validator) CheckStruct(name string, s interface{}) *ValidatorErr {
	err := v.Validate.Struct(s)
	if err != nil {
		return NewValidatorErr(name, err)
	}
	return nil
}
