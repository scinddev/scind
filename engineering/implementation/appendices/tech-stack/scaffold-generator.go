// internal/generator/generator.go
// Override file generation logic

package generator

import (
	"fmt"

	"github.com/spf13/afero"
)

// Generator creates override files for workspaces
type Generator struct {
	fs afero.Fs
}

// New creates a new Generator with the given filesystem
func New(fs afero.Fs) *Generator {
	return &Generator{fs: fs}
}

// Generate creates override files for all applications in a workspace
func (g *Generator) Generate(workspacePath string) error {
	// 1. Load workspace.yaml
	// 2. For each application in the workspace:
	//    a. Load application.yaml
	//    b. Resolve active flavor
	//    c. Generate override file with:
	//       - Docker labels for Traefik routing
	//       - Network attachments
	//       - Port mappings for assigned ports
	//    d. Generate manifest.yaml with computed values
	// 3. Write files to .generated/ directory

	return nil
}

// GenerateApp creates override files for a single application
func (g *Generator) GenerateApp(workspacePath, appName string) error {
	// Implementation for single-app generation
	return nil
}

// internal/generator/override.go
// Docker Compose override file builder

// OverrideBuilder constructs Docker Compose override content
type OverrideBuilder struct {
	services map[string]ServiceOverride
}

type ServiceOverride struct {
	Labels   map[string]string
	Networks []string
	Ports    []string
}

// NewOverrideBuilder creates a new override builder
func NewOverrideBuilder() *OverrideBuilder {
	return &OverrideBuilder{
		services: make(map[string]ServiceOverride),
	}
}

// AddTraefikLabels adds Traefik routing labels for a service
func (b *OverrideBuilder) AddTraefikLabels(service string, labels map[string]string) {
	svc := b.services[service]
	if svc.Labels == nil {
		svc.Labels = make(map[string]string)
	}
	for k, v := range labels {
		svc.Labels[k] = v
	}
	b.services[service] = svc
}

// AddNetwork attaches a service to a network
func (b *OverrideBuilder) AddNetwork(service, network string) {
	svc := b.services[service]
	svc.Networks = append(svc.Networks, network)
	b.services[service] = svc
}

// AddPort adds a port mapping for assigned ports
func (b *OverrideBuilder) AddPort(service string, hostPort, containerPort int) {
	svc := b.services[service]
	svc.Ports = append(svc.Ports, fmt.Sprintf("%d:%d", hostPort, containerPort))
	b.services[service] = svc
}

// Build returns the override content as YAML
func (b *OverrideBuilder) Build() ([]byte, error) {
	// Marshal to YAML format
	return nil, nil
}

// internal/generator/traefik.go
// Traefik label generation

// GenerateTraefikLabels creates Traefik routing labels for a service
func GenerateTraefikLabels(workspace, app, service, hostname, protocol string, port int) map[string]string {
	routerName := fmt.Sprintf("%s-%s-%s", workspace, app, service)

	labels := map[string]string{
		"traefik.enable": "true",
		fmt.Sprintf("traefik.http.routers.%s.rule", routerName): fmt.Sprintf("Host(`%s`)", hostname),
		fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", routerName): fmt.Sprintf("%d", port),
	}

	// Add TLS configuration for HTTPS
	if protocol == "https" {
		labels[fmt.Sprintf("traefik.http.routers.%s.tls", routerName)] = "true"
	}

	return labels
}
