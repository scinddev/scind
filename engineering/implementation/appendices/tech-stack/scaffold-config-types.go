// Scaffold: internal/config/workspace.go and internal/config/application.go
// Migrated from specs/scind-go-stack.md

package config

// =============================================================================
// Workspace Configuration (workspace.yaml)
// =============================================================================

// WorkspaceConfig represents workspace.yaml
type WorkspaceConfig struct {
	Workspace WorkspaceSpec `yaml:"workspace" validate:"required"`
}

type WorkspaceSpec struct {
	Name         string                    `yaml:"name" validate:"required,dns_label"`
	Domain       string                    `yaml:"domain,omitempty"`
	Templates    *TemplateOverrides        `yaml:"templates,omitempty"`
	Applications map[string]ApplicationRef `yaml:"applications,omitempty"`
}

type TemplateOverrides struct {
	Hostname    string `yaml:"hostname,omitempty"`
	Alias       string `yaml:"alias,omitempty"`
	ProjectName string `yaml:"project_name,omitempty"`
}

type ApplicationRef struct {
	Path       string `yaml:"path,omitempty"`
	Repository string `yaml:"repository,omitempty"`
}

// =============================================================================
// Application Configuration (application.yaml)
// =============================================================================

// ApplicationConfig represents application.yaml
type ApplicationConfig struct {
	ExportedServices map[string]ExportedService `yaml:"exported_services" validate:"required"`
	Flavors          map[string]Flavor          `yaml:"flavors,omitempty"`
	DefaultFlavor    string                     `yaml:"default_flavor,omitempty"`
}

type ExportedService struct {
	Service string `yaml:"service,omitempty" validate:"omitempty"` // Optional: defaults to map key
	Ports   []Port `yaml:"ports" validate:"required,dive"`
}

type Port struct {
	Type       string `yaml:"type" validate:"required,oneof=proxied assigned"`
	Protocol   string `yaml:"protocol,omitempty" validate:"required_if=Type proxied,omitempty,oneof=http https tcp postgresql mysql"`
	Port       int    `yaml:"port,omitempty" validate:"omitempty,min=1,max=65535"` // Optional: inferred from Compose service
	Visibility string `yaml:"visibility,omitempty" validate:"omitempty,oneof=public protected"`
}

type Flavor struct {
	ComposeFiles []string `yaml:"compose_files" validate:"required,min=1"`
}

// =============================================================================
// Config Loading: Inference and Defaults
// =============================================================================

// The config loader (internal/config/loader.go) applies these inference rules after unmarshaling:
//
// Service name defaulting (C-3):
// - If ExportedService.Service is empty, set it to the map key from exported_services
// - Example: exported_services.web with no service: field defaults to Compose service "web"
//
// Port inference (C-2):
// - If Port.Port is zero (omitted), infer from the Compose service's ports: configuration
// - If the Compose service has exactly one port, use that port
// - If the Compose service has multiple ports, return a clear error:
//   Error: Port must be specified for exported service "web"
//     Application: frontend
//     Compose service "web" has multiple ports: 80, 443, 9229
//     Specify which port to use in application.yaml
//
// Compose file existence validation (A-10):
// - At generate time, validate that all files in Flavor.ComposeFiles exist on disk
// - If a file is missing, return a clear error:
//   Error: Flavor "full" references non-existent file: docker-compose.worker.yaml
//     Application: backend
//     Available compose files: docker-compose.yaml, docker-compose.dev.yaml
