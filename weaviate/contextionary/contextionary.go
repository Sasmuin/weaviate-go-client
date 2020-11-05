package contextionary

import (
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/models"
)

// API for the contextionary endpoint
type API struct {
	connection *connection.Connection
}

// New Contextionary api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// ConceptsGetter get builder to query weaviate concepts
func (c11y *API) ConceptsGetter() *ConceptGetter {
	return &ConceptGetter{
		connection: c11y.connection,
	}
}

// ExtensionCreator get a builder to extend weaviates contextionary
func (c11y *API) ExtensionCreator() *ExtensionCreator {
	return &ExtensionCreator{
		connection: c11y.connection,
		extension: &models.C11yExtension{
			Weight: 1.0,
		},
	}
}
