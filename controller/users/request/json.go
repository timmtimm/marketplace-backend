package request

import (
	"crop_connect/business/users"
	"crop_connect/helper"
	"errors"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUser struct {
	RegionID    string `form:"regionID" json:"regionID" validate:"required"`
	Name        string `form:"name" json:"name" validate:"required"`
	Description string `form:"description" json:"description"`
	Email       string `form:"email" json:"email" validate:"required,email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" validate:"required,min=10,max=13,number"`
	Password    string `form:"password" json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
	Role        string `form:"role" json:"role" validate:"required"`
}

func (req *RegisterUser) ToDomain() (*users.Domain, error) {
	regionObjID, err := primitive.ObjectIDFromHex(req.RegionID)
	if err != nil {
		return nil, errors.New("id daerah tidak valid")
	}

	return &users.Domain{
		RegionID:    regionObjID,
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		Role:        req.Role,
	}, nil
}

func (req *RegisterUser) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
			out := make([]helper.ValidationError, len(ve))

			for i, e := range ve {
				out[i] = helper.ValidationError{
					Field:   e.Field(),
					Message: helper.MessageForTag(e.Tag()),
				}

				out[i].Message = strings.Replace(out[i].Message, "[PARAM]", e.Param(), 1)

				for _, f := range fields {
					if f.Name() == e.Field() {
						out[i].Field = f.Tag("json")
						break
					}
				}
			}
			return out
		}
	}

	return nil
}

type Login struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
}

func (req *Login) ToDomain() *users.Domain {
	return &users.Domain{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req *Login) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
			out := make([]helper.ValidationError, len(ve))

			for i, e := range ve {
				out[i] = helper.ValidationError{
					Field:   e.Field(),
					Message: helper.MessageForTag(e.Tag()),
				}

				out[i].Message = strings.Replace(out[i].Message, "[PARAM]", e.Param(), 1)

				for _, f := range fields {
					if f.Name() == e.Field() {
						out[i].Field = f.Tag("json")
						break
					}
				}
			}
			return out
		}
	}

	return nil
}

type Update struct {
	RegionID    string `form:"regionID" json:"regionID" validate:"required"`
	Name        string `form:"name" json:"name" validate:"required"`
	Description string `form:"description" json:"description"`
	Email       string `form:"email" json:"email" validate:"required,email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" validate:"required,min=10,max=13,number"`
}

func (req *Update) ToDomain() (*users.Domain, error) {
	regionObjID, err := primitive.ObjectIDFromHex(req.RegionID)
	if err != nil {
		return nil, errors.New("id daerah tidak valid")
	}

	return &users.Domain{
		RegionID:    regionObjID,
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}, nil
}

func (req *Update) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
			out := make([]helper.ValidationError, len(ve))

			for i, e := range ve {
				out[i] = helper.ValidationError{
					Field:   e.Field(),
					Message: helper.MessageForTag(e.Tag()),
				}

				out[i].Message = strings.Replace(out[i].Message, "[PARAM]", e.Param(), 1)

				for _, f := range fields {
					if f.Name() == e.Field() {
						out[i].Field = f.Tag("json")
						break
					}
				}
			}
			return out
		}
	}

	return nil
}

type RegisterValidator struct {
	RegionID    string `form:"regionID" json:"regionID" validate:"required"`
	Name        string `form:"name" json:"name" validate:"required"`
	Description string `form:"description" json:"description"`
	Email       string `form:"email" json:"email" validate:"required,email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" validate:"required,min=10,max=13,number"`
	Password    string `form:"password" json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
}

func (req *RegisterValidator) ToDomain() (*users.Domain, error) {
	regionObjID, err := primitive.ObjectIDFromHex(req.RegionID)
	if err != nil {
		return nil, errors.New("id daerah tidak valid")
	}

	return &users.Domain{
		RegionID:    regionObjID,
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}, nil
}

func (req *RegisterValidator) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
			out := make([]helper.ValidationError, len(ve))

			for i, e := range ve {
				out[i] = helper.ValidationError{
					Field:   e.Field(),
					Message: helper.MessageForTag(e.Tag()),
				}

				out[i].Message = strings.Replace(out[i].Message, "[PARAM]", e.Param(), 1)

				for _, f := range fields {
					if f.Name() == e.Field() {
						out[i].Field = f.Tag("json")
						break
					}
				}
			}
			return out
		}
	}

	return nil
}

type ChangePassword struct {
	OldPassword string `form:"oldPassword" json:"oldPassword" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
	NewPassword string `form:"newPassword" json:"newPassword" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
}

func (req *ChangePassword) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
			out := make([]helper.ValidationError, len(ve))

			for i, e := range ve {
				out[i] = helper.ValidationError{
					Field:   e.Field(),
					Message: helper.MessageForTag(e.Tag()),
				}

				out[i].Message = strings.Replace(out[i].Message, "[PARAM]", e.Param(), 1)

				for _, f := range fields {
					if f.Name() == e.Field() {
						out[i].Field = f.Tag("json")
						break
					}
				}
			}
			return out
		}
	}

	return nil
}

func (req *ChangePassword) ToDomain() *users.Domain {
	return &users.Domain{
		Password: req.OldPassword,
	}
}
