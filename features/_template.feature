# Feature Template
#
# This file provides a template for writing Gherkin feature specifications.
# Features describe expected behavior from a user's perspective.

Feature: [Feature Name]
  As a [type of user]
  I want [to perform some action]
  So that [I can achieve some goal]

  Background:
    # Common setup steps that run before each scenario
    Given a workspace named "dev" exists
    And the workspace has an application "app-one"

  Scenario: [Scenario Name]
    # Describe a specific user flow
    Given [some precondition]
    When [some action is taken]
    Then [some result should occur]
    And [another result should occur]

  Scenario: [Another Scenario]
    Given [some different precondition]
    When [some action is taken]
    Then [some result should occur]

  Scenario Outline: [Parameterized Scenario]
    # Use scenario outlines for testing multiple variations
    Given a port of type "<type>" with protocol "<protocol>"
    When the override file is generated
    Then the environment variable should contain port "<expected_port>"

    Examples:
      | type     | protocol | expected_port |
      | proxied  | https    | 443           |
      | proxied  | http     | 80            |
      | assigned | -        | 5432          |
