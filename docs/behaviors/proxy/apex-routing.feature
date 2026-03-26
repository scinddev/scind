# This feature verifies behaviors from:
# See: ../../specs/naming-conventions.md
# See: ../../specs/generated-override-files.md
# See: ../../specs/docker-labels.md

Feature: Apex URL Routing
  As a developer
  I want my application's primary service to have a short apex URL
  So that I can access it without remembering the exported service name

  Background:
    Given a workspace named "dev" with domain "scind.test"

  Scenario: Single export gets implicit apex hostname
    Given an application "frontend" with a single proxied HTTPS export "web" on port 80
    When the override file is generated
    Then the service has a Traefik router "dev-frontend-web-https" for "dev-frontend-web.scind.test"
    And the service has an apex Traefik router "dev-frontend-https" for "dev-frontend.scind.test"
    And both hostnames route to container port 80

  Scenario: Multi-export explicit primary gets apex hostname
    Given an application "frontend" with proxied exports "web" and "api"
    And export "web" is marked as primary
    When the override file is generated
    Then the "web" service has an apex Traefik router for "dev-frontend.scind.test"
    And the "api" service does not have an apex Traefik router

  Scenario: Multi-export without primary gets no apex
    Given an application "frontend" with proxied exports "web" and "api"
    And no export is marked as primary
    When the override file is generated
    Then no service has an apex Traefik router

  Scenario: Apex internal alias on workspace network
    Given an application "frontend" with a single proxied export "web" on port 80
    When the override file is generated
    Then the service has network alias "frontend-web" on the internal network
    And the service has apex network alias "frontend" on the internal network

  Scenario: Apex Docker labels on primary export
    Given an application "frontend" with a single proxied HTTPS export "web"
    When the override file is generated
    Then the service has label "scind.apex.host=dev-frontend.scind.test"
    And the service has label "scind.apex.proxy.https.url=https://dev-frontend.scind.test"

  Scenario: Apex router naming convention
    Given an application "frontend" with a single proxied HTTPS export "web"
    When the override file is generated
    Then the apex router is named "dev-frontend-https"
    And the export router is named "dev-frontend-web-https"
