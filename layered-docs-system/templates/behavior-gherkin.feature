# Feature: [Feature Name]
#
# This is a Gherkin feature file template for executable specifications.
# Use this format to define behaviors that can be verified automatically.
#
# Gherkin keywords:
#   Feature: Groups related scenarios
#   Background: Steps run before each scenario
#   Scenario: A single test case
#   Scenario Outline: A parameterized test case
#   Given: Preconditions (setup)
#   When: Actions taken
#   Then: Expected outcomes
#   And/But: Additional steps (takes meaning from previous keyword)
#
# Best practices:
#   - Write from the user's perspective
#   - Keep scenarios independent (no shared state)
#   - Use declarative style (what, not how)
#   - One behavior per scenario

Feature: [Feature Name]
  As a [role/persona]
  I want [goal/desire]
  So that [benefit/value]

  # Reference: [Link to specification doc]

  Background:
    # Steps that run before each scenario
    Given [common precondition]
    And [another common precondition]

  # ============================================
  # Happy Path Scenarios
  # ============================================

  Scenario: [Descriptive name of the happy path]
    Given [precondition]
    And [another precondition]
    When [action taken]
    Then [expected outcome]
    And [another expected outcome]

  Scenario: [Another happy path variation]
    Given [precondition]
    When [action taken]
    Then [expected outcome]

  # ============================================
  # Edge Cases
  # ============================================

  Scenario: [Edge case name]
    Given [unusual precondition]
    When [action taken]
    Then [expected behavior in edge case]

  # ============================================
  # Error Handling
  # ============================================

  Scenario: [Error scenario name]
    Given [precondition that will cause error]
    When [action taken]
    Then [error is handled gracefully]
    And [appropriate error message is shown]

  # ============================================
  # Parameterized Scenarios
  # ============================================

  Scenario Outline: [Parameterized scenario name]
    Given [precondition with <parameter>]
    When [action with <parameter>]
    Then [expected outcome with <expected_result>]

    Examples:
      | parameter | expected_result |
      | value1    | result1         |
      | value2    | result2         |
      | value3    | result3         |

  # ============================================
  # Integration Scenarios
  # ============================================

  @integration
  Scenario: [Integration with other feature/system]
    Given [precondition involving multiple components]
    When [action that crosses boundaries]
    Then [outcome verified across components]
