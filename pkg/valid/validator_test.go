package valid

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/kms9/publicyc/pkg/util/debug"
)

func TestValidator_Check(t *testing.T) {
	type fields struct {
		Validate *validator.Validate
	}
	type args struct {
		field       string
		validateTag string
		value       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ValidatorErr
	}{
		{
			name: "check",
			fields: fields{
				Validate: validator.New(),
			},
			args: args{
				field:       "subsectionId",
				validateTag: "required,uuid",
				value:       "37f004fa-57f7-11e7-913c-bf2058c40ec6",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				Validate: tt.fields.Validate,
			}
			got := v.Check(tt.args.field, tt.args.validateTag, tt.args.value)
			debug.Print(got)
			debug.Print(got.String())
		})
	}
}

func TestValidator_CheckFields(t *testing.T) {
	type fields struct {
		Validate *validator.Validate
	}
	type args struct {
		err []*ValidatorErr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ValidatorErr
	}{
		{
			name   : "checkFields",
			fields: fields{
				Validate: validator.New(),
			},
			args   : args{err: []*ValidatorErr{
				New().Check("subsectionId", "required,uuid", "37f004"),
				New().Check("uid", "required,len=24", "uid1"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				Validate: tt.fields.Validate,
			}
			got := v.CheckFields(tt.args.err...)
			debug.Print(got.String())
		})
	}
}