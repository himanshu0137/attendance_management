package app

import (
	"context"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"hunrx.com/attendance_management/pkg/db"
	"hunrx.com/attendance_management/pkg/events"
	"hunrx.com/attendance_management/pkg/httpx"
)

// FeatureModule is the contract every feature module registers with the core.
// Migrate installs the module's PostgreSQL schema and Mount wires its handlers.
type FeatureModule interface {
	Name() string
	Migrate(ctx context.Context, db *db.DB) error
	Mount(core *CoreContext, rg *gin.RouterGroup)
}

// registered holds the modules initialised by their package init() functions.
var registered []FeatureModule

// Register appends a module to the registry, called from each module's init().
func Register(m FeatureModule) {
	registered = append(registered, m)
}

// ResolveModules returns the modules to activate. An empty or omitted list
// activates every registered module, deterministically ordered by Name().
func ResolveModules(requested []string) ([]FeatureModule, error) {
	if len(requested) == 0 {
		return registered, nil
	}
	valid := map[string]bool{}
	for _, m := range registered {
		valid[m.Name()] = true
	}
	active := []FeatureModule{}
	for _, name := range requested {
		if !valid[name] {
			return nil, unknownModule{name, registered}
		}
		active = append(active, byName(name))
	}
	sort.Slice(active, func(i, j int) bool {
		return active[i].Name() < active[j].Name()
	})
	return active, nil
}

func byName(name string) FeatureModule {
	for _, m := range registered {
		if m.Name() == name {
			return m
		}
	}
	return nil
}

type unknownModule struct {
	name    string
	modules []FeatureModule
}

func (u unknownModule) Error() string {
	names := []string{}
	for _, m := range u.modules {
		names = append(names, m.Name())
	}
	return "unknown module " + u.name + "; registered modules are " + join(names, ", ")
}

func join(names []string, sep string) string {
	if len(names) == 0 {
		return ""
	}
	b := names[0]
	for _, n := range names[1:] {
		b += sep + n
	}
	return b
}

// CoreContext is the shared core injected into every active module.
type CoreContext struct {
	Engine   *gin.Engine
	Pool     *pgxpool.Pool
	Events   events.Bus
	Modules  []FeatureModule
	TenantID string
	Keycloak *httpx.KeycloakConfig
}
