package elastic

/*
*
Documentation of the KNN query
*/
type KnnQuery struct {
	boost         *float32
	field         string
	filter        []Query
	k             int
	numCandidates int
	queryVector   []float64
	// TODO: Support QueryVectorBuilder
}

// NewKnnQuery creates and initializes a new KnnQuery.
func NewKnnQuery(field string, k int, numCandidates int, queryVector []float64) *KnnQuery {
	return &KnnQuery{
		field:         field,
		k:             k,
		numCandidates: numCandidates,
		queryVector:   queryVector,
	}
}

func (q *KnnQuery) Boost(boost *float32) *KnnQuery {
	q.boost = boost
	return q
}

func (q *KnnQuery) Filter(filter []Query) *KnnQuery {
	q.filter = filter
	return q
}

func (q *KnnQuery) Field(field string) *KnnQuery {
	q.field = field
	return q
}

func (q *KnnQuery) K(k int) *KnnQuery {
	q.k = k
	return q
}

func (q *KnnQuery) NumCandidates(numCandidates int) *KnnQuery {
	q.numCandidates = numCandidates
	return q
}

func (q *KnnQuery) QueryVector(queryVector []float64) *KnnQuery {
	q.queryVector = queryVector
	return q
}

func (q *KnnQuery) Source() (interface{}, error) {
	// {
	//   "knn": {
	//     "field": "field",
	//     "query_vector": [1.0, 2.0, 3.0],
	//     "k": 10,
	//     "num_candidates": 100,
	//     "boost": 1.0,
	//     "filter": [
	//       {
	//         "term": {
	//           "field": "value"
	//         }
	//       }
	//     ]
	//   }
	// }
	source := make(map[string]interface{})
	knn := make(map[string]interface{})
	source["knn"] = knn

	knn["field"] = q.field
	knn["k"] = q.k
	knn["num_candidates"] = q.numCandidates
	knn["query_vector"] = q.queryVector

	if q.boost != nil {
		knn["boost"] = *q.boost
	}

	if len(q.filter) > 0 {
		if len(q.filter) == 1 {
			knn["filter"] = q.filter[0]
		} else {
			knn["filter"] = q.filter
		}
	}

	return source, nil
}
