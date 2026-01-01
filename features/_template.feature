# Contrail Behavior Specification Template
#
# This is a Gherkin template for defining executable behavior specifications.
# Each .feature file should focus on a single feature or user story.

Feature: [Feature Name]
  [Brief description of the feature and its value to users]

  Background:
    # Common setup steps that apply to all scenarios in this feature
    Given the proxy is running
    And a workspace "dev" exists

  # Happy path scenario
  Scenario: [Descriptive scenario name]
    Given [initial context]
    When [action is performed]
    Then [expected outcome]
    And [additional assertions]

  # Edge case or alternative flow
  Scenario: [Alternative scenario name]
    Given [different initial context]
    When [action is performed]
    Then [different expected outcome]

  # Scenario with examples (parameterized testing)
  Scenario Outline: [Parameterized scenario description]
    Given an application with <flavor> flavor
    When the user runs "contrail up"
    Then <compose_files> compose files are loaded

    Examples:
      | flavor  | compose_files |
      | lite    | 1             |
      | full    | 3             |

  # Error scenario
  Scenario: [Error condition description]
    Given [context that will cause an error]
    When [action is performed]
    Then the command fails with exit code <code>
    And the error message contains "<message>"
