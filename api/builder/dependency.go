package builder

import "github.com/Ayut0/coffee-journal/api/config"

// Dependency holds all shared dependencies wired at startup.
// Issue #17 adds Pool *pgxpool.Pool.
type Dependency struct {
	Cfg *config.Config
}
