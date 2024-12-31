package blueprint_test

import (
	"os"
	"path"
	"testing"

	"github.com/husamAwadhi/go-parser/internal/blueprint"
	"github.com/stretchr/testify/assert"
)

/*
Blueprint Builder (HusamAwadhi\PowerParserTests\Blueprint\BlueprintBuilder)

	✔ Create from string
	✔ Create from path
	✔ Throwing exception on empty stream

Components (HusamAwadhi\PowerParserTests\Blueprint\Components\Components)

	✔ Throwing exception on invalid blueprint
	✔ Create from parameters
	✔ Create from invalid parameters

Conditions (HusamAwadhi\PowerParserTests\Blueprint\Components\Conditions)

	✔ Create from parameters
	✔ Create from invalid parameters

Fields (HusamAwadhi\PowerParserTests\Blueprint\Components\Fields)

	✔ Create from parameters
	✔ Create from invalid parameters
*/
func TestSuccessfullyCreateBlueprint(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory: %v", err)
	}

	var testCases = []struct {
		name  string
		input string
	}{
		{"Create Simple Blueprint", "valid.yaml"},
		{"Create Blueprint with processors", "valid_with_processors.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath := path.Join(path.Dir(path.Dir(dir)), "assets", "blueprints", tc.input)
			bp, err := blueprint.CreateBlueprintFromFile(filePath)

			assert.Nil(t, err)
			assert.IsType(t, &blueprint.Blueprint{}, bp)
		})
	}
}

func TestSuccessfullyCreateBlueprintFromBytes(t *testing.T) {
	yaml := `
%YAML 1.2
---
version: "1.0"
meta:
  file:
    extension: csv
    name: "sample"
blueprint:
  - name: header_info
    type: hit
    conditions:
      - column: [1]
        is: "Cashier Number"
    fields:
      - name: currency
        position: 6
`

	bp, err := blueprint.CreateBlueprintFromBytes([]byte(yaml))

	assert.Nil(t, err)
	assert.IsType(t, &blueprint.Blueprint{}, bp)
}

func TestErrorOnInvalidBlueprint(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory: %v", err)
	}

	var testCases = []struct {
		name  string
		input string
		error []string
	}{
		{"Blueprint does not exist", "does_not_exist.yaml", []string{"no such file or directory"}},
		{"invalid Blueprint 1", "invalid_blueprint_1.yaml", []string{"[count=4]", "Metadata", "Components[0].Name", "Components[1].Name", "Components[2].Name"}},
		{"invalid Blueprint 2", "invalid_blueprint_2.yaml", []string{"[count=1]", "Components"}},
		{"invalid Blueprint 3", "invalid_blueprint_3.yaml", []string{"[count=3]", "Version", "Conditions[0].is,isNot,anyOf,noneOf", "Components[0].Type"}},
		{"invalid Blueprint 4", "invalid_blueprint_4.yaml", []string{"[count=2]", "Conditions[0].is,isNot,anyOf,noneOf", "Metadata.File.Name"}},
		{"invalid Blueprint 5", "invalid_with_processors.yaml", []string{"[count=4]", "Fields[0].Format.Code", "Fields[0].Format.Parameter", "Fields[5].Format.Parameter", "Fields[6].Format.Code"}},
		{"invalid Blueprint 6", "invalid_exe.yaml", []string{"[count=2]", "Blueprint.Metadata.File.Extension", "Fields[0].Type"}},
		{"invalid Blueprint 6", "invalid.yaml", []string{"unknown field"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			filePath := path.Join(path.Dir(path.Dir(dir)), "assets", "blueprints", tc.input)
			bp, err := blueprint.CreateBlueprintFromFile(filePath)

			assert.Nil(t, bp)
			assert.NotNil(t, err)
			for _, e := range tc.error {
				assert.ErrorContains(t, err, e)
			}
		})
	}
}
