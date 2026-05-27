package scores

import "iter"

type EffectRunner struct {
	queue  map[EncodedEffect]Effect
	driver *ScoreEngineDriver
}

func NewEffectRunner(driver *ScoreEngineDriver) *EffectRunner {
	return &EffectRunner{
		queue:  make(map[EncodedEffect]Effect),
		driver: driver,
	}
}

func (r *EffectRunner) RunChainEffects(effects iter.Seq[Effect]) {
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
	})

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
	})

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
	})
}

func (r *EffectRunner) Run(effects iter.Seq[Effect]) {
	for e := range effects {
		var chainEffects iter.Seq[Effect]

		switch effect := e.(type) {
		case EffectRankClass:
			r.driver.logger.Info("re-ranking comp class", "comp_class_id", effect.CompClassID)
			r.driver.engine.RankCompClass(effect.CompClassID)
		case EffectScoreContender:
			r.driver.logger.Info("re-scoring contender", "contender_id", effect.ContenderID)
			chainEffects = r.driver.engine.ScoreContender(effect.ContenderID)
		case EffectCalculatePointValues:
			r.driver.logger.Info("re-calculating point values", "comp_class_id", effect.CompClassID, "problem_id", effect.ProblemID)
			chainEffects = r.driver.engine.CalculatePointValues(effect.CompClassID, effect.ProblemID)
		}

		if chainEffects == nil {
			continue
		}

		for chainEffect := range chainEffects {
			r.queue[chainEffect.Encode()] = chainEffect
		}
	}
}
