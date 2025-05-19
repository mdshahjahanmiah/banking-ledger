package banking_ledger

import (
	"testing"

	"github.com/cucumber/godog"
	"github.com/mdshahjahanmiah/banking-ledger/features"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	features.InitializeScenario(ctx)
}

func TestFeatures(t *testing.T) {
	status := godog.TestSuite{
		Name:                "banking-ledger",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features/account.feature"},
		},
	}.Run()

	if status != 0 {
		t.Fatalf("non-zero status returned, failed to run feature tests")
	}
}
