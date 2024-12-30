package blueprint_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/husamAwadhi/go-parser/internal/blueprint"
)

/*
Blueprint (HusamAwadhi\PowerParserTests\Blueprint\Blueprint)

	✔ Successfully create blueprint
	✔ Throwing exception on invalid blueprint

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
func TestBlueprint(t *testing.T) {

	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory: %v", err)
	}

	filePath := path.Join(path.Dir(path.Dir(dir)), "assets", "blueprints", "valid.yaml")
	fmt.Println(filePath)

	_, err = blueprint.CreateBlueprintFromFile(filePath)

	if err != nil {
		t.Errorf("Error creating blueprint: %v", err)
	}
}
