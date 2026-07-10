package scores

import (
	"iter"

	"github.com/climblive/platform/backend/internal/domain"
)

type effectResolver interface {
	RankCompClass(domain.CompClassID)
	ScoreContender(domain.ContenderID) iter.Seq[Effect]
	CalculatePointValues(domain.CompClassID, domain.ProblemID) iter.Seq[Effect]
}

type EffectRunner struct {
	queue    map[EncodedEffect]Effect
	resolver effectResolver
}

func NewEffectRunner(resolver effectResolver) *EffectRunner {
	return &EffectRunner{
		queue:    make(map[EncodedEffect]Effect),
		resolver: resolver,
	}
}

func (r *EffectRunner) RunEffects(effects iter.Seq[Effect]) {
	if effects == nil {
		return
	}

	r.run(func(yield func(Effect) bool) {
		for effect := range effects {
			switch effect.(type) {
			case EffectCalculatePointValues:
			default:
				r.queue[effect.Encode()] = effect
				continue
			}

			if !yield(effect) {
				return
			}
		}
	})

	r.run(func(yield func(Effect) bool) {
		for _, effect := range r.queue {
			switch effect.(type) {
			case EffectScoreContender:
			default:
				continue
			}

			if !yield(effect) {
				return
			}
		}
	})

	r.run(func(yield func(Effect) bool) {
		for _, effect := range r.queue {
			switch effect.(type) {
			case EffectRankClass:
			default:
				continue
			}

			if !yield(effect) {
				return
			}
		}
	})
}

func (r *EffectRunner) run(effects iter.Seq[Effect]) {
	for e := range effects {
		var chainEffects iter.Seq[Effect]

		switch effect := e.(type) {
		case EffectRankClass:
			r.resolver.RankCompClass(effect.CompClassID)
		case EffectScoreContender:
			chainEffects = r.resolver.ScoreContender(effect.ContenderID)
		case EffectCalculatePointValues:
			chainEffects = r.resolver.CalculatePointValues(effect.CompClassID, effect.ProblemID)
		}

		if chainEffects == nil {
			continue
		}

		for chainEffect := range chainEffects {
			if chainEffect.Type() <= e.Type() {
				continue
			}

			r.queue[chainEffect.Encode()] = chainEffect
		}
	}
}
