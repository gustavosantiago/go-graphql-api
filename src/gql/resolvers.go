package gql

import (
	"postgres"

	"github.com/graphql-go/graphql"
)

// Resolver struct holds a connection to our db
type Resolver struct {
	db *postgres.Db
}

// UserResolver our user query throught a dbcall to GetUsersByNAme
func (r *Resolver) UserResolver(p graphql.ResolveParams) (interface{}, error) {
	// Strip name from arguments
	name, ok := p.Args["name"].(string)

	if ok {
		users := r.db.GetUsersByName(name)
		return users, nil
	}

	return nil, nil
}
