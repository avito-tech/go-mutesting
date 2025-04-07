package annotation

import (
	"go/ast"
	"go/token"
)

type ChainCollector interface {
	Handle(name string, comment *ast.Comment, fset *token.FileSet, file *ast.File, fileAbs string)
	SetNext(next ChainCollector)
}

type BaseCollector struct {
	next ChainCollector
}

func (h *BaseCollector) SetNext(next ChainCollector) {
	h.next = next
}

func (h *BaseCollector) Handle(name string, comment *ast.Comment, fset *token.FileSet, file *ast.File, fileAbs string) {
	if h.next != nil {
		h.next.Handle(name, comment, fset, file, fileAbs)
	}
}

type RegexAnnotationCollector struct {
	BaseCollector
	Processor RegexAnnotation
}

type NextLineAnnotationCollector struct {
	BaseCollector
	Processor LineAnnotation
}

func (r *RegexAnnotationCollector) Handle(name string, comment *ast.Comment, fset *token.FileSet, file *ast.File, fileAbs string) {
	if name == RegexpAnnotation {
		r.Processor.collectMatchNodes(comment, fset, file, fileAbs)
	} else {
		r.BaseCollector.Handle(name, comment, fset, file, fileAbs)
	}
}

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
