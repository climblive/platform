package scores

import (
	"iter"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

func FilterByClass(contenders map[domain.ContenderID]*Contender, compClassID domain.CompClassID) iter.Seq[*Contender] {
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

func CompareScore(s1, s2 domain.Score) bool {
	s1.Timestamp = time.Time{}
	s2.Timestamp = time.Time{}

	return s1 == s2
}
