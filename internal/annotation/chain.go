package annotation

import (
	"go/ast"
	"go/token"
)

// ChainCollector defines the interface for handlers in the annotation processing chain.
// Implementations should handle specific annotation types and pass unhandled cases to the next handler in the chain.
type ChainCollector interface {
	// Handle processes an annotation if it matches the handler's type,
	// otherwise delegates to the next handler in the chain.
	Handle(name string, comment *ast.Comment, fset *token.FileSet, file *ast.File, fileAbs string)
	// SetNext establishes the next handler in the chain of responsibility.
	SetNext(next ChainCollector)
}

// BaseCollector provides default chain handling behavior
type BaseCollector struct {
	next ChainCollector
}

// SetNext sets the next handler in the chain of responsibility.
func (h *BaseCollector) SetNext(next ChainCollector) {
	h.next = next
}

// Handle implements the default chain behavior by delegating to the next handler.
func (h *BaseCollector) Handle(name string, comment *ast.Comment, fset *token.FileSet, file *ast.File, fileAbs string) {
	if h.next != nil {
		h.next.Handle(name, comment, fset, file, fileAbs)
	}
}

// RegexAnnotationCollector implements the ChainCollector interface for "mutator-disable-regexp" annotations.
type RegexAnnotationCollector struct {
	BaseCollector
	Processor RegexAnnotation
}

// NextLineAnnotationCollector implements the ChainCollector interface for "mutator-disable-next-line" annotations.
type NextLineAnnotationCollector struct {
	BaseCollector
	Processor LineAnnotation
}

// Handle processes regex pattern annotations, delegating other types to the next handler.
func (r *RegexAnnotationCollector) Handle(name string, comment *ast.Comment, fset *token.FileSet, file *ast.File, fileAbs string) {
	if name == RegexpAnnotation {
		r.Processor.collectMatchNodes(comment, fset, file, fileAbs)
	} else {
		r.BaseCollector.Handle(name, comment, fset, file, fileAbs)
	}
}

// Handle processes regex pattern annotations, delegating other types to the next handler.
func (n *NextLineAnnotationCollector) Handle(name string, comment *ast.Comment, fset *token.FileSet, file *ast.File, fileAbs string) {
	if name == NextLineAnnotation {
		n.Processor.collectNodesOnNextLine(comment, fset, file)
	} else {
		n.BaseCollector.Handle(name, comment, fset, file, fileAbs)
	}
}

func (p *Processor) buildChain() ChainCollector {
	regexHandler := &RegexAnnotationCollector{Processor: p.RegexAnnotation}
	nextLineHandler := &NextLineAnnotationCollector{Processor: p.LineAnnotation}
	regexHandler.SetNext(nextLineHandler)

	return regexHandler
}
