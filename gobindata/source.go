// gobindata creates migrate source from go-bindata
package gobindata

import (
	"github.com/golang-migrate/migrate/v4/source"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

// NewSourceFrom creates a source builder from go bindata assets.
func NewSourceFrom(
	AssetNames func() []string,
	Asset func(string) ([]byte, error),
) func() (string, source.Driver, error) {
	s := bindata.Resource(
		AssetNames(),
		Asset,
	)

	return func() (string, source.Driver, error) {
		driver, err := bindata.WithInstance(s)
		if err != nil {
			return "go-bindata", nil, err
		}
		return "go-bindata", driver, nil
	}
}
