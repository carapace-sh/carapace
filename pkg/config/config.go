package config

import "github.com/rsteube/carapace/internal/config"


func Register(name string, i interface{}) { config.RegisterConfig(name, i) }
