package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gondola/util/yaml"
)

// Constants set by gondola/cache/layer, read by gondola/app.Context
const (
	LayerCachedKey          = "___gondola_cached"
	LayerServedFromCacheKey = "___gondola_layer_served_from_cache"
)

var (
	inTest      bool
	goRun       bool
	inAppEngine bool
)

// InTest returns true iff called when running
// from go test.
func InTest() bool {
	return inTest
}

// IsGoRun returns true iff called when running
// from go run.
func IsGoRun() bool {
	return goRun
}

func InAppEngine() bool {
	return inAppEngine
}

func InAppEngineDevServer() bool {
	return os.Getenv("RUN_WITH_DEVAPPSERVER") != ""
}

func AppEngineAppId() string {
	var m map[string]interface{}
	if err := yaml.UnmarshalFile("app.yaml", &m); err == nil {
		// XXX: your-app-id is the default in app.yaml in GAE templates, found
		// in the gondolaweb repository. Keep these in sync.
		if id, ok := m["application"].(string); ok && id != "your-app-id" {
			return id
		}
	}
	return ""
}

func AppEngineAppHost() string {
	if id := AppEngineAppId(); id != "" {
		return fmt.Sprintf("http://%s.appspot.com", id)
	}
	return ""
}

func init() {
	inTest = strings.Contains(os.Args[0], string(filepath.Separator)+"_test"+string(filepath.Separator))
	goRun = strings.Contains(os.Args[0], "_obj"+string(filepath.Separator)+"exe")
}
