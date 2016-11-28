package dbi

import (
	"../dbmodel"
)

// MediaTypeTblDBI - testing
type MediaTypeTblDBI interface {
	// AddMediatype - testing
	AddMediaType(mtDetails *dbmodel.MediaTypeEntry) error
}
