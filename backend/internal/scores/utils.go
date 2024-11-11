package scores

import (
	"iter"
	"maps"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

func FilterByCompClass(contenders map[domain.ContenderID]*Contender, compClassID domain.CompClassID) iter.Seq[*Contender] {
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

func CompClasses(contenders map[domain.ContenderID]*Contender) iter.Seq[domain.CompClassID] {
	compClassIDs := make(map[domain.CompClassID]struct{})

	for contender := range maps.Values(contenders) {
		compClassIDs[contender.CompClassID] = struct{}{}
	}

	return maps.Keys(compClassIDs)
}

func CompareScore(s1, s2 domain.Score) bool {
	s1.Timestamp = time.Time{}
	s2.Timestamp = time.Time{}

	return s1 == s2
}
