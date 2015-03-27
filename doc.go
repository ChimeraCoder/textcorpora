// Package textcorpora proivdes an interface for various corpora used in natural language processing.
//
// Currently the package provides two corpora: the Carnegie Mellon Pronouncing Dictionary Corpus, and the Enron Corpus (containing over 600,000 emails from Enron employees).
//
// The location for each corpus is stored in a location provided by appdirs. For example, on Linux, the current version of the CMU corpus will be downloaded and saved to ~/.local/share/cmudict/.1/cmudict.0.7a.corpus.
//
// TextCorpora is a helper package that provides an interface for various corpora. It was originally written for use in the ReadingLevel library. It is provided as a separate package for convenience - both to faciliate use of corpora in other applications and libraries, and also to allow users of the ReadingLevel library the ability to plug in an alternative corpus if desired.
package textcorpora
