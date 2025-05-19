package features

import (
	"fmt"
	"github.com/cucumber/godog"
)

var (
	existingAccounts = map[string]bool{}
	currentAccountID string
	accountBalance   int
	accountStatus    string
	createError      error
)

func AnExistingAccountWithID(id string) error {
	existingAccounts[id] = true
	return nil
}

func IHaveANewAccountWithIDAndUserID(id, userID string) error {
	// Just save the account ID for this scenario
	currentAccountID = id
	// For this example, userID is ignored but can be saved if needed
	return nil
}

func ICreateTheAccount() error {
	if existingAccounts[currentAccountID] {
		createError = fmt.Errorf("duplicate account")
		return nil
	}
	existingAccounts[currentAccountID] = true
	accountBalance = 0
	accountStatus = "active"
	createError = nil
	return nil
}

func TheAccountShouldBeCreatedSuccessfully() error {
	if createError != nil {
		return fmt.Errorf("expected success but got error: %v", createError)
	}
	return nil
}

func TheAccountBalanceShouldBe(expected int) error {
	if accountBalance != expected {
		return fmt.Errorf("expected balance %d but got %d", expected, accountBalance)
	}
	return nil
}

func TheAccountStatusShouldBe(expected string) error {
	if accountStatus != expected {
		return fmt.Errorf("expected status %q but got %q", expected, accountStatus)
	}
	return nil
}

func TheCreationShouldFailWithError(expected string) error {
	if createError == nil || createError.Error() != expected {
		return fmt.Errorf("expected error %q but got %v", expected, createError)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^an existing account with ID "([^"]*)"$`, AnExistingAccountWithID)
	ctx.Step(`^I have a new account with ID "([^"]*)" and user ID "([^"]*)"$`, IHaveANewAccountWithIDAndUserID)
	ctx.Step(`^I create the account$`, ICreateTheAccount)
	ctx.Step(`^the account should be created successfully$`, TheAccountShouldBeCreatedSuccessfully)
	ctx.Step(`^the account balance should be (\d+)$`, TheAccountBalanceShouldBe)
	ctx.Step(`^the account status should be "([^"]*)"$`, TheAccountStatusShouldBe)
	ctx.Step(`^the creation should fail with error "([^"]*)"$`, TheCreationShouldFailWithError)
}
