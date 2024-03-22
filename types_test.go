package example

import (
	"fmt"
	"strings"
	"testing"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{
			"ArrayOfMaps",
			"[map[example:456 extra:4]]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("EXAMPLE_VALUES", tc.value)

			t.Run("Straight", func(t *testing.T) {
				config := &Configuration{}

				require.NoError(t, loadConfig("EXAMPLE_", "example.yml", config))

				assert.Len(t, config.Values, 1)
			})

			t.Run("ViaMap", func(t *testing.T) {
				configmap := map[string]any{}

				require.NoError(t, loadConfig("EXAMPLE_", "example.yml", &configmap))

				assert.Equal(t, map[string]any{"values": []map[string]any{{"example": "456", "extra": 4}}}, configmap)
				assert.NotEqual(t, map[string]any{"values": "[map[example:456 extra:4]]"}, configmap)

				config, err := mapToConfig(configmap)
				require.NoError(t, err)

				assert.Len(t, config.Values, 1)
			})
		})
	}
}

func loadConfig(envPrefix, path string, o any) (err error) {
	k := koanf.New(".")

	if err = k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}

	if err = k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil); err != nil {
		return fmt.Errorf("error loading env: %w", err)
	}

	if err = k.Unmarshal("", o); err != nil {
		return fmt.Errorf("error unmarshalling: %w", err)
	}

	return nil
}

func mapToConfig(o map[string]any) (config *Configuration, err error) {
	k := koanf.New(".")

	if err = k.Load(confmap.Provider(o, "."), nil); err != nil {
		return nil, fmt.Errorf("error loading map: %w", err)
	}

	config = &Configuration{}

	if err = k.Unmarshal("", config); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	return config, nil
}
