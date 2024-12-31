package blueprint

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type FieldFormatCodes string

const (
	StringCode  FieldFormatCodes = "s"
	FloatCode   FieldFormatCodes = "f"
	DateCode    FieldFormatCodes = "d"
	DefaultCode FieldFormatCodes = "default"
)

type SupportedFiles string

const (
	Xlsx SupportedFiles = "xlsx"
	Xls  SupportedFiles = "xls"
	Csv  SupportedFiles = "csv"
)

type ComponentTypes string

const (
	Hit  ComponentTypes = "hit"
	Next ComponentTypes = "next"
)

type FieldTypes string

const (
	Int        FieldTypes = "int"
	Float      FieldTypes = "float"
	Bool       FieldTypes = "bool"
	BoolStrict FieldTypes = "bool-strict"
	Date       FieldTypes = "date"
)

type File struct {
	Extension SupportedFiles `validate:"required,is-supported-file"`
	Name      string         `validate:"required"`
}

type Metadata struct {
	File *File `validate:"required"`
}

type Condition struct {
	Column []uint32 `validate:"required"`
	Is     string   `yaml:",omitempty"`
	IsNot  string   `yaml:"isNot,omitempty"`
	AnyOf  string   `yaml:"anyOf,omitempty"`
	NoneOf string   `yaml:"noneOf,omitempty"`
}

type FieldFormat struct {
	Code    FieldFormatCodes `validate:"is-valid-field-format"`
	Parameter string `validate:"required_with=Code"`
}

func (ff *FieldFormat) UnmarshalYAML(data []byte) error {
	if string(data) == "" || string(data) == "\"\"" {
		return nil
	}

	format := strings.Split(string(data), "%")

	switch format[0] {
	case "s":
		ff.Code = StringCode
	case "f":
		ff.Code = FloatCode
	case "d":
		ff.Code = DateCode
	default:
		ff.Code = DefaultCode
	}

	if len(format) == 2 {
		ff.Parameter = format[1]
	}

	return nil
}

type Field struct {
	Name     string     `validate:"required"`
	Position uint32     `validate:"required"`
	Type     FieldTypes `validate:"is-valid-field-type"`
	Format   FieldFormat
}

type Components struct {
	Name       string `validate:"required"`
	Mandatory  bool
	Table      bool
	Page       uint32
	Type       ComponentTypes `validate:"is-valid-component-type,required"`
	Conditions []*Condition   `yaml:",omitempty" validate:"dive"`
	Fields     []*Field       `validate:"required,dive"`
}

type Blueprint struct {
	Version    string        `validate:"required"`
	Metadata   *Metadata     `yaml:"meta" validate:"required"`
	Components []*Components `yaml:"blueprint" validate:"required,dive"`
}

var validate *validator.Validate

func CreateBlueprintFromFile(fileDirectory string) (*Blueprint, error) {
	yml, err := os.ReadFile(fileDirectory)
	if err != nil {
		return nil, err
	}

	return CreateBlueprintFromBytes(yml)
}

func CreateBlueprintFromBytes(yml []byte) (*Blueprint, error) {
	var blueprint Blueprint

	err := yaml.UnmarshalWithOptions(yml, &blueprint, yaml.DisallowUnknownField())
	if err != nil {
		return nil, err
	}

	newValidator()

	if err := validate.Struct(blueprint); err != nil {
		return nil, handleErrors(err.(validator.ValidationErrors))
	}

	return &blueprint, nil
}

func newValidator() {
	validate = validator.New()

	validate.RegisterStructValidation(ConditionValidator, Condition{})
	validate.RegisterValidation("is-supported-file", IsSupportedFile)
	validate.RegisterValidation("is-valid-component-type", IsValidComponentType)
	validate.RegisterValidation("is-valid-field-type", IsValidFieldType)
	validate.RegisterValidation("is-valid-field-format", IsValidFieldFormat)
}

func ConditionValidator(sl validator.StructLevel) {
	condition := sl.Current().Interface().(Condition)

	if len(condition.Is) == 0 &&
		len(condition.IsNot) == 0 &&
		len(condition.AnyOf) == 0 &&
		len(condition.NoneOf) == 0 {
		sl.ReportError(condition.Is, "is,isNot,anyOf,noneOf", "Condition", "no-condition", "")
	}
}

func IsSupportedFile(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case string(Xlsx), string(Xls), string(Csv):
		return true
	}
	return false
}

func IsValidComponentType(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case string(Hit), string(Next):
		return true
	}
	return false
}

func IsValidFieldType(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case string(Int), string(Float), string(Bool), string(BoolStrict), string(Date), "":
		return true
	}
	return false
}

func IsValidFieldFormat(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case string(StringCode), string(FloatCode), string(DateCode), "":
		return true
	}
	return false
}
