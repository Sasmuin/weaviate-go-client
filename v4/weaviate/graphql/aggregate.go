package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/semi-technologies/weaviate/entities/models"
)

// AggregateBuilder for the aggregate GraphQL query string
type AggregateBuilder struct {
	connection                rest
	fields                    []Field
	className                 string
	includesFilterClause      bool // true if brackets behind class is needed
	groupByClausePropertyName string
	withWhereFilter           *WhereArgumentBuilder
	withNearTextFilter        *NearTextArgumentBuilder
	withNearVectorFilter      *NearVectorArgumentBuilder
	withNearObjectFilter      *NearObjectArgumentBuilder
	includesObjectLimit       bool
	objectLimit               int
}

// WithFields that should be included in the aggregation query e.g. `meta{count}`
func (ab *AggregateBuilder) WithFields(fields []Field) *AggregateBuilder {
	ab.fields = fields
	return ab
}

// WithClassName that should be aggregated
func (ab *AggregateBuilder) WithClassName(name string) *AggregateBuilder {
	ab.className = name
	return ab
}

// WithWhere adds the where filter.
func (ab *AggregateBuilder) WithWhere(where *WhereArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withWhereFilter = where
	return ab
}

// WithGroupBy adds the group by property clause as the filter.
//  The group by value/path clause still needs to be set in the WithFields field.
func (ab *AggregateBuilder) WithGroupBy(propertyName string) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.groupByClausePropertyName = propertyName
	return ab
}

// WithNearText clause to find close objects
func (ab *AggregateBuilder) WithNearText(nearText *NearTextArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearTextFilter = nearText
	return ab
}

// WithNearObject clause to find close objects
func (ab *AggregateBuilder) WithNearObject(nearObject *NearObjectArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearObjectFilter = nearObject
	return ab
}

// WithNearVector clause to find close objects
func (ab *AggregateBuilder) WithNearVector(nearVector *NearVectorArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearVectorFilter = nearVector
	return ab
}

// WithObjectLimit specifies max number of vector search results to return
func (ab *AggregateBuilder) WithObjectLimit(objectLimit int) *AggregateBuilder {
	ab.objectLimit = objectLimit
	ab.includesObjectLimit = true
	return ab
}

// Do execute the aggregation query
func (ab *AggregateBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, ab.connection, ab.build())
}

func (ab *AggregateBuilder) createFilterClause() string {
	filters := []string{}
	if len(ab.groupByClausePropertyName) > 0 {
		filters = append(filters, fmt.Sprintf(`groupBy: "%v"`, ab.groupByClausePropertyName))
	}
	if ab.withWhereFilter != nil {
		filters = append(filters, ab.withWhereFilter.build())
	}
	if ab.withNearTextFilter != nil {
		filters = append(filters, ab.withNearTextFilter.build())
	}
	if ab.withNearVectorFilter != nil {
		filters = append(filters, ab.withNearVectorFilter.build())
	}
	if ab.withNearObjectFilter != nil {
		filters = append(filters, ab.withNearObjectFilter.build())
	}
	if ab.includesObjectLimit {
		filters = append(filters, fmt.Sprintf("objectLimit: %d", ab.objectLimit))
	}

	return fmt.Sprintf("(%s)", strings.Join(filters, ", "))
}

func (ab *AggregateBuilder) createFieldsClause() string {
	if len(ab.fields) > 0 {
		fields := make([]string, len(ab.fields))
		for i := range ab.fields {
			fields[i] = ab.fields[i].build()
		}
		return strings.Join(fields, " ")
	}
	return ""
}

// build the query string
func (ab *AggregateBuilder) build() string {
	filterClause := ""
	if ab.includesFilterClause {
		filterClause = ab.createFilterClause()
	}
	fields := ab.createFieldsClause()
	return fmt.Sprintf(`{Aggregate{%v%v{%v}}}`, ab.className, filterClause, fields)
}