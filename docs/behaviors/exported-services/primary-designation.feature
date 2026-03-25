# This feature verifies behaviors from:
# See: ../../specs/configuration-schemas.md
# See: ../../specs/environment-variables.md
# See: ../../decisions/0013-apex-url-primary-designation.md

Feature: Primary Export Designation
  As a developer
  I want to designate one exported service as primary
  So that my application gets a short apex URL

  Background:
    Given a workspace named "dev" with domain "scind.test"

  Scenario: Implicit primary for single export
    Given an application "frontend" with one exported service "web"
    And "web" is proxied with protocol "https"
    When configuration is validated
    Then "web" is treated as the primary export
    And apex environment variables are generated:
      | Variable                    | Value                           |
      | SCIND_FRONTEND_APEX_HOST    | dev-frontend.scind.test         |
      | SCIND_FRONTEND_APEX_PORT    | 443                             |
      | SCIND_FRONTEND_APEX_SCHEME  | https                           |
      | SCIND_FRONTEND_APEX_URL     | https://dev-frontend.scind.test |

  Scenario: Explicit primary for multi-export application
    Given an application "frontend" with exports "web" and "api"
    And export "web" has "primary: true"
    And export "api" does not have "primary: true"
    When configuration is validated
    Then "web" is the primary export
    And apex environment variables reference the "web" service

  Scenario: No primary designation for multi-export
    Given an application "shared-db" with exports "db" and "cache"
    And neither export is marked primary
    When configuration is validated
    Then no apex hostname is generated
    And no apex environment variables are generated
    And no apex Docker labels are generated

  Scenario: Validation error for multiple primaries
    Given an application "frontend" with exports "web" and "api"
    And export "web" has "primary: true"
    And export "api" has "primary: true"
    When configuration is validated
    Then a validation error is emitted
    And the error message indicates multiple primary exports

  Scenario: Primary on assigned-port export
    Given an application "shared-db" with one exported service "db"
    And "db" is of type "assigned" with port 5432
    When configuration is validated
    Then "db" is treated as the primary export
    And the apex internal alias "shared-db" is created
    But no apex hostname is generated
    And no apex environment variables are generated
    And no apex Docker labels are generated

  Scenario: Primary on proxied export in mixed-type application
    Given an application "backend" with exports:
      | name   | type     | protocol | primary |
      | api    | proxied  | https    | true    |
      | worker | assigned | -        | false   |
    When configuration is validated
    Then "api" is the primary export
    And apex hostname "dev-backend.scind.test" is generated
    And apex internal alias "backend" is created
    And apex environment variables are generated for "api"
