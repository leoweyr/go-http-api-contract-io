package client

import (
	"net/http"
)

// PreprocessorChain sequences a set of preprocessors applied in declaration order around a base round tripper.
type PreprocessorChain struct {
	preprocessors []Preprocessor
}

// NewPreprocessorChain creates a PreprocessorChain from the given preprocessors, applied so the first declared runs outermost.
func NewPreprocessorChain(preprocessors ...Preprocessor) PreprocessorChain {
	var chain PreprocessorChain = PreprocessorChain{preprocessors: preprocessors}

	return chain
}

// Use returns a new PreprocessorChain with the given preprocessor appended to the sequence.
func (chain PreprocessorChain) Use(preprocessor Preprocessor) PreprocessorChain {
	var extended []Preprocessor = make([]Preprocessor, 0, len(chain.preprocessors)+1)

	extended = append(extended, chain.preprocessors...)
	extended = append(extended, preprocessor)

	var extendedChain PreprocessorChain = PreprocessorChain{preprocessors: extended}

	return extendedChain
}

// Then wraps the base round tripper with every preprocessor so the first declared runs outermost, and returns the composed round tripper.
func (chain PreprocessorChain) Then(base http.RoundTripper) http.RoundTripper {
	var composed http.RoundTripper = base

	var index int

	for index = len(chain.preprocessors) - 1; index >= 0; index = index - 1 {
		composed = chain.preprocessors[index].Call(composed)
	}

	return composed
}
