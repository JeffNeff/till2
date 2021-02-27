package build

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/graph"
	"bridgedl/translate"
)

// GraphBuilder builds a graph by applying a series of sequential
// transformation steps.
type GraphBuilder struct {
	Bridge      *config.Bridge
	Translators *translate.TranslatorProviders
}

// Build iterates over the transformation steps of the GraphBuilder to build a graph.
func (b *GraphBuilder) Build() (*graph.DirectedGraph, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	steps := []GraphTransformer{
		// Add all blocks as graph vertices
		&AddComponentsTransformer{
			Bridge: b.Bridge,
		},

		// Attach block translators
		&AttachTranslatorsTransformer{
			Translators: b.Translators,
		},

		// Resolve references and connect vertices
		&ConnectReferencesTransformer{},
	}

	g := graph.NewDirectedGraph()

	for _, step := range steps {
		trsfDiags := step.Transform(g)
		diags = diags.Extend(trsfDiags)
	}

	return g, diags
}
