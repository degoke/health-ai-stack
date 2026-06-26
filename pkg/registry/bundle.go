package registry

import "embed"

//go:embed internal/bundles/r4/structure-definitions/*.json internal/bundles/r4/search-parameters/*.json
var r4BundleFS embed.FS
