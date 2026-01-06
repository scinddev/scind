// scaffold-config-types.go
// Configuration type scaffolds for Contrail
// Create as: internal/config/workspace.go and internal/config/application.go

// --- internal/config/workspace.go ---

package config

// WorkspaceConfig represents workspace.yaml
type WorkspaceConfig struct {
    Workspace WorkspaceSpec `yaml:"workspace" validate:"required"`
}

type WorkspaceSpec struct {
    Name         string                      `yaml:"name" validate:"required,dns_label"`
    Domain       string                      `yaml:"domain,omitempty"`
    Templates    *TemplateOverrides          `yaml:"templates,omitempty"`
    Applications map[string]ApplicationRef   `yaml:"applications,omitempty"`
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

// --- internal/config/application.go ---

// ApplicationConfig represents application.yaml
type ApplicationConfig struct {
    ExportedServices map[string]ExportedService `yaml:"exported_services" validate:"required"`
    Flavors          map[string]Flavor          `yaml:"flavors,omitempty"`
    DefaultFlavor    string                     `yaml:"default_flavor,omitempty"`
}

type ExportedService struct {
    Service string `yaml:"service,omitempty" validate:"omitempty"`  // Optional: defaults to map key
    Ports   []Port `yaml:"ports" validate:"required,dive"`
}

type Port struct {
    Type       string `yaml:"type" validate:"required,oneof=proxied assigned"`
    Protocol   string `yaml:"protocol,omitempty" validate:"required_if=Type proxied,omitempty,oneof=http https tcp postgresql mysql"`
    Port       int    `yaml:"port,omitempty" validate:"omitempty,min=1,max=65535"`  // Optional: inferred from Compose service
    Visibility string `yaml:"visibility,omitempty" validate:"omitempty,oneof=public protected"`  // Defaults to "protected"
}

type Flavor struct {
    ComposeFiles []string `yaml:"compose_files" validate:"required,min=1"`
}
