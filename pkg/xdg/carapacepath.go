package xdg

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// CarapacePathEntry is a resolved top-level carapace root directory with its group.
type CarapacePathEntry struct {
	Path  string
	Group string // "user" or "system"
}

// CarapacePaths returns ordered carapace top-level directories.
//
// If CARAPACE_PATH is set, it is a path-list (os.PathListSeparator) where each
// entry may end with "@user" or "@system". Entries starting with "~" are treated
// as user group and "~/.config" (or leading "~" + "/.config") is replaced by
// UserConfigDir() to honor XDG_CONFIG_HOME when configured. Untagged entries
// default to "system".
func CarapacePaths() ([]CarapacePathEntry, error) {
	if v, ok := os.LookupEnv("CARAPACE_PATH"); ok && v != "" {
		return parseCarapacePath(v)
	}
	return defaultCarapacePaths()
}

func parseCarapacePath(v string) ([]CarapacePathEntry, error) {
	userBase, err := UserConfigDir()
	if err != nil {
		return nil, err
	}
	userBase = filepath.FromSlash(userBase)

	var out []CarapacePathEntry
	for entry := range strings.SplitSeq(v, string(os.PathListSeparator)) {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		group := "system"
		if strings.HasSuffix(entry, "@user") {
			group = "user"
			entry = strings.TrimSuffix(entry, "@user")
		} else if strings.HasSuffix(entry, "@system") {
			group = "system"
			entry = strings.TrimSuffix(entry, "@system")
		}
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		if strings.HasPrefix(entry, "~") {
			group = "user"
			// Replace leading ~/.config with userBase (XDG config dir)
			if entry == "~/.config" || entry == "~/.config/" {
				entry = userBase
			} else if strings.HasPrefix(entry, "~/.config/") {
				rest := strings.TrimPrefix(entry, "~/.config/")
				entry = filepath.Join(userBase, filepath.FromSlash(rest))
			} else if entry == "~" {
				entry = userBase
			} else if strings.HasPrefix(entry, "~/") {
				rest := strings.TrimPrefix(entry, "~/")
				// If it looks like ".config/..." already handled above; otherwise join under userBase
				if strings.HasPrefix(rest, ".config/") || rest == ".config" {
					rest2 := strings.TrimPrefix(rest, ".config")
					rest2 = strings.TrimPrefix(rest2, "/")
					if rest2 == "" {
						entry = userBase
					} else {
						entry = filepath.Join(userBase, filepath.FromSlash(rest2))
					}
				} else {
					entry = filepath.Join(userBase, filepath.FromSlash(rest))
				}
			} else {
				entry = filepath.FromSlash(strings.TrimPrefix(entry, "~"))
				if !filepath.IsAbs(entry) {
					entry = filepath.Join(userBase, entry)
				}
			}
		} else {
			entry = filepath.FromSlash(entry)
		}

		out = append(out, CarapacePathEntry{Path: filepath.ToSlash(entry), Group: group})
	}
	return out, nil
}

func defaultCarapacePaths() ([]CarapacePathEntry, error) {
	userBase, err := UserConfigDir()
	if err != nil {
		return nil, err
	}
	userBase = filepath.ToSlash(userBase)

	var out []CarapacePathEntry
	// user: XDG config dir + "/carapace"
	out = append(out, CarapacePathEntry{Path: filepath.ToSlash(filepath.Join(filepath.FromSlash(userBase), "carapace")), Group: "user"})

	switch runtime.GOOS {
	case "windows":
		if v, ok := os.LookupEnv("PROGRAMDATA"); ok && v != "" {
			out = append(out, CarapacePathEntry{Path: filepath.ToSlash(filepath.Join(filepath.FromSlash(v), "carapace")), Group: "system"})
		}
	case "darwin":
		out = append(out, CarapacePathEntry{Path: "/Library/Application Support/carapace", Group: "system"})
	default:
		dataDirs, err := DataDirs()
		if err != nil {
			return nil, err
		}
		for _, d := range dataDirs {
			out = append(out, CarapacePathEntry{Path: filepath.ToSlash(filepath.Join(filepath.FromSlash(d), "carapace")), Group: "system"})
		}
	}
	return out, nil
}
