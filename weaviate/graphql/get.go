package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate/entities/models"
)

// GetBuilder for GraphQL
type GetBuilder struct {
	connection rest
	className  string
	withFields []Field

	includesFilterClause bool // true if brackets behind class is needed
	includesLimit        bool
	limit                int
	includesOffset       bool
	offset               int
	includesAfter        bool
	after                string
	consistencyLevel     string
	withWhereFilter      *filters.WhereBuilder
	withNearTextFilter   *NearTextArgumentBuilder
	withNearVectorFilter *NearVectorArgumentBuilder
	withNearObjectFilter *NearObjectArgumentBuilder
	withGroupFilter      *GroupArgumentBuilder
	withAskFilter        *AskArgumentBuilder
	withNearImageFilter  *NearImageArgumentBuilder
	withSort             *SortBuilder
	withBM25             *BM25ArgumentBuilder
	withHybrid           *HybridArgumentBuilder
	withGenerativeSearch *GenerativeSearchBuilder
	withGroupBy          *GroupByArgumentBuilder
}

// WithAfter is part of the Cursor API. It can be used to extract all elements
// by specfiying the last ID from the previous "page". Cannot be combined with
// any other filters or search operators other than limit. Requires
// WithClassName() and WithLimit() to be set.
func (gb *GetBuilder) WithAfter(id string) *GetBuilder {
	gb.includesAfter = true
	gb.after = id
	return gb
}

// WithClassName that should be queried
func (gb *GetBuilder) WithClassName(name string) *GetBuilder {
	gb.className = name
	return gb
}

// WithFields included in the result set
func (gb *GetBuilder) WithFields(fields ...Field) *GetBuilder {
	gb.withFields = fields
	return gb
}

// WithWhere filter
func (gb *GetBuilder) WithWhere(where *filters.WhereBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withWhereFilter = where
	return gb
}

// WithLimit of objects in the result set
func (gb *GetBuilder) WithLimit(limit int) *GetBuilder {
	gb.includesFilterClause = true
	gb.includesLimit = true
	gb.limit = limit
	return gb
}

// WithOffset of objects in the result set
func (gb *GetBuilder) WithOffset(offset int) *GetBuilder {
	gb.includesFilterClause = true
	gb.includesOffset = true
	gb.offset = offset
	return gb
}

// WithBM25 to search the inverted index
func (gb *GetBuilder) WithBM25(bm25 *BM25ArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withBM25 = bm25
	return gb
}

// WithHybrid to combine multiple searches
func (gb *GetBuilder) WithHybrid(hybrid *HybridArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withHybrid = hybrid
	return gb
}

// WithNearText clause to find close objects
func (gb *GetBuilder) WithNearText(nearText *NearTextArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withNearTextFilter = nearText
	return gb
}

// WithNearObject clause to find close objects
func (gb *GetBuilder) WithNearImage(nearImage *NearImageArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withNearImageFilter = nearImage
	return gb
}

// WithNearVector clause to find close objects
func (gb *GetBuilder) WithNearVector(nearVector *NearVectorArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withNearVectorFilter = nearVector
	return gb
}

// WithGroup statement
func (gb *GetBuilder) WithGroup(group *GroupArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withGroupFilter = group
	return gb
}

// WithAsk clause to find an aswer to the question
func (gb *GetBuilder) WithAsk(ask *AskArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withAskFilter = ask
	return gb
}

// WithNearObject clause to find close objects
func (gb *GetBuilder) WithNearObject(nearObject *NearObjectArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withNearObjectFilter = nearObject
	return gb
}

// WithSort included in the result set
func (gb *GetBuilder) WithSort(sort ...Sort) *GetBuilder {
	gb.includesFilterClause = true
	gb.withSort = &SortBuilder{sort}
	return gb
}

// To Use OpenAI to generate a response based on the results
func (gb *GetBuilder) WithGenerativeSearch(s *GenerativeSearchBuilder) *GetBuilder {
	gb.withGenerativeSearch = s
	return gb
}

// WithGroupBy to perform group by operation
func (gb *GetBuilder) WithGroupBy(groupBy *GroupByArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withGroupBy = groupBy
	return gb
}

func (gb *GetBuilder) WithConsistencyLevel(lvl string) *GetBuilder {
	gb.includesFilterClause = true
	gb.consistencyLevel = lvl
	return gb
}

// Do execute the GraphQL query
func (gb *GetBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, gb.connection, gb.build())
}

// build the GraphQL query string (not needed when Do is executed)
func (gb *GetBuilder) build() string {
	filterClause := ""
	if gb.includesFilterClause {
		filterClause = gb.createFilterClause()
	}
	fieldsClause := gb.createFieldsClause()

	query := fmt.Sprintf("{Get {%v %v {%v}}}", gb.className, filterClause, fieldsClause)

	return query
}

func (gb *GetBuilder) createFilterClause() string {
	filters := []string{}
	if gb.withWhereFilter != nil {
		filters = append(filters, gb.withWhereFilter.String())
	}
	if gb.withNearTextFilter != nil {
		filters = append(filters, gb.withNearTextFilter.build())
	}
	if gb.withBM25 != nil {
		filters = append(filters, gb.withBM25.build())
	}
	if gb.withHybrid != nil {
		filters = append(filters, gb.withHybrid.build())
	}
	if gb.withNearVectorFilter != nil {
		filters = append(filters, gb.withNearVectorFilter.build())
	}
	if gb.withNearObjectFilter != nil {
		filters = append(filters, gb.withNearObjectFilter.build())
	}
	if gb.withAskFilter != nil {
		filters = append(filters, gb.withAskFilter.build())
	}
	if gb.withNearImageFilter != nil {
		filters = append(filters, gb.withNearImageFilter.build())
	}
	if gb.withGroupFilter != nil {
		filters = append(filters, gb.withGroupFilter.build())
	}
	if gb.consistencyLevel != "" {
		filters = append(filters, fmt.Sprintf("consistencyLevel: %s", gb.consistencyLevel))
	}
	if gb.includesLimit {
		filters = append(filters, fmt.Sprintf("limit: %v", gb.limit))
	}
	if gb.includesOffset {
		filters = append(filters, fmt.Sprintf("offset: %v", gb.offset))
	}
	if gb.includesAfter {
		filters = append(filters, fmt.Sprintf("after: %q", gb.after))
	}
	if gb.withSort != nil {
		filters = append(filters, gb.withSort.build())
	}
	if gb.withGroupBy != nil {
		filters = append(filters, gb.withGroupBy.build())
	}
	return fmt.Sprintf("(%s)", strings.Join(filters, ", "))
}

func (gb *GetBuilder) createFieldsClause() string {
	if len(gb.withFields) == 0 && gb.withGenerativeSearch == nil {
		return ""
	}

	if gb.withGenerativeSearch == nil {
		return joinFields(gb.withFields)
	}

	generate := gb.withGenerativeSearch.build()
	generateAdditional := Field{Name: "_additional", Fields: []Field{generate}}

	if len(gb.withFields) == 0 {
		return generateAdditional.build()
	}

	// check if _additional field exists. If missing just add new _additional with generate,
	// if exists merge generate into present one
	posAdditional := -1
	for i := range gb.withFields {
		if gb.withFields[i].Name == "_additional" {
			posAdditional = i
			break
		}
	}

	if posAdditional == -1 {
		return joinFields(append(gb.withFields, generateAdditional))
	}

	mergedAdditional := Field{
		Name:   "_additional",
		Fields: append(gb.withFields[posAdditional].Fields, generate),
	}

	fields := append(gb.withFields[:posAdditional], gb.withFields[posAdditional+1:]...)
	return joinFields(append(fields, mergedAdditional))
}

func joinFields(fields []Field) string {
	strFields := []string{}
	for i := range fields {
		if str := strings.TrimSpace(fields[i].build()); len(str) > 0 {
			strFields = append(strFields, str)
		}
	}
	return strings.Join(strFields, " ")
}
