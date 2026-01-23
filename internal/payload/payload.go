package payload

import (
	"fmt"
	"os"
	"path/filepath"

	baseTypes "github.com/jpalaniselvam/myna/internal/types"
	"github.com/jpalaniselvam/myna/internal/varsub"
)

// LoadPayload resolves the payload content from the 'file' or 'data' fields in the payload configuration.
// Precedence:
// 1. 'file': Reads content from the specified file path.
// 2. 'data': Uses the inline string content.
func Resolve(payload baseTypes.Payload, baseDir string, resolver *varsub.Resolver) ([]byte, error) {
	// 1. Check for 'file' payload (Precedence: High)
	if payload.File != "" {
		var path = payload.File
		// Resolve relative paths against the base directory
		if !filepath.IsAbs(path) {
			path = filepath.Join(baseDir, path)
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read payload file %s: %w", path, err)
		}
		return resolver.Resolve(b), nil
	}

	// 2. Check for 'data' payload (if file not handled)
	if payload.Data != "" {
		return []byte(payload.Data), nil
	}

	// Return nil if no payload is found (this is valid, some actions might not have payload)
	return nil, nil
}
