package evaluator

type Flag int

type EvaluatorOptions struct {
	Flag Flag
}

type EvaluatorOption func(opts *EvaluatorOptions)
