package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/semantics"
	"github.com/semi-technologies/weaviate-go-client/weaviate/models"
	"net/http"
)

// ReferenceDeleter builder to remove a reference from a data object
type ReferenceDeleter struct {
	connection        *connection.Connection
	semanticKind      semantics.Kind
	uuid              string
	referenceProperty string
	referencePayload  *models.SingleRef
}

// WithID specifies the uuid of the object on which the reference will be deleted
func (rr *ReferenceDeleter) WithID(uuid string) *ReferenceDeleter {
	rr.uuid = uuid
	return rr
}

// WithReferenceProperty specifies the property on which the reference should be deleted
func (rr *ReferenceDeleter) WithReferenceProperty(propertyName string) *ReferenceDeleter {
	rr.referenceProperty = propertyName
	return rr
}

// WithKind specifies the semantic kind that is used for the data object
// If not called the builder defaults to `things`
func (rr *ReferenceDeleter) WithKind(semanticKind semantics.Kind) *ReferenceDeleter {
	rr.semanticKind = semanticKind
	return rr
}

// WithReference specifies reference payload of the reference about to be deleted
func (rr *ReferenceDeleter) WithReference(referencePayload *models.SingleRef) *ReferenceDeleter {
	rr.referencePayload = referencePayload
	return rr
}

// Do remove the reference defined by the payload set in this builder to the property and object defined in this builder
func (rr *ReferenceDeleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v/references/%v", string(rr.semanticKind), rr.uuid, rr.referenceProperty)
	responseData, responseErr := rr.connection.RunREST(ctx, path, http.MethodDelete, *rr.referencePayload)
	return except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 204)
}
