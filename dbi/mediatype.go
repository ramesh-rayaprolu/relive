package dbi

import (
	"github.com/msproject/relive/dbmodel"
)

// MediaTypeTblDBI - testing
type MediaTypeTblDBI interface {
	// AddMediatype - testing
	AddMediaType(mtDetails *dbmodel.MediaTypeEntry) error
}
