package scores

import (
	"iter"
	"log/slog"

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

func (r *EffectRunner) RunEffects(effects iter.Seq[Effect], logger *slog.Logger) {
	if effects == nil {
		return
	}

	r.Run(func(yield func(Effect) bool) {
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
	}, logger)

	r.Run(func(yield func(Effect) bool) {
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
	}, logger)

	r.Run(func(yield func(Effect) bool) {
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
	}, logger)
}

func (r *EffectRunner) Run(effects iter.Seq[Effect], logger *slog.Logger) {
	for e := range effects {
		var chainEffects iter.Seq[Effect]

		switch effect := e.(type) {
		case EffectRankClass:
			logger.Info("re-ranking comp class", "comp_class_id", effect.CompClassID)
			r.resolver.RankCompClass(effect.CompClassID)
		case EffectScoreContender:
			logger.Info("re-scoring contender", "contender_id", effect.ContenderID)
			chainEffects = r.resolver.ScoreContender(effect.ContenderID)
		case EffectCalculatePointValues:
			logger.Info("re-calculating point values", "comp_class_id", effect.CompClassID, "problem_id", effect.ProblemID)
			chainEffects = r.resolver.CalculatePointValues(effect.CompClassID, effect.ProblemID)
		}

		if chainEffects == nil {
			continue
		}

		for chainEffect := range chainEffects {
			r.queue[chainEffect.Encode()] = chainEffect
		}
	}
}
