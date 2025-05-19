Feature: Account Creation

  Scenario: Successfully create a new account
    Given an existing account with ID "acc123"
    And I have a new account with ID "acc124" and user ID "user456"
    When I create the account
    Then the account should be created successfully
    And the account balance should be 0
    And the account status should be "active"

  Scenario: Fail to create an account with duplicate ID
    Given an existing account with ID "acc123"
    And I have a new account with ID "acc123" and user ID "user789"
    When I create the account
    Then the creation should fail with error "duplicate account"
