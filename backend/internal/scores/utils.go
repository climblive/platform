package scores

import (
	"iter"

	"github.com/climblive/platform/backend/internal/domain"
)

func FilterByClass(contenders map[domain.ResourceID]*Contender, compClassID domain.ResourceID) iter.Seq[*Contender] {
	return func(yield func(*Contender) bool) {
		for _, contender := range contenders {
			if contender.CompClassID != compClassID {
				continue
			}

			if !yield(contender) {
				return
			}
		}
	}
}
