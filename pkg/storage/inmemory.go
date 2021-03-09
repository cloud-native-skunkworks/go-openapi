package storage

import "github.com/hashicorp/go-memdb"

func LoadLocalDB() (*memdb.MemDB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": &memdb.TableSchema{
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"username": &memdb.IndexSchema{
						Name:    "username",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
					},
					"firstname": &memdb.IndexSchema{
						Name:    "firstname",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "FirstName"},
					},
					"lastname": &memdb.IndexSchema{
						Name:    "lastname",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "LastName"},
					},
					"phone": &memdb.IndexSchema{
						Name:    "phone",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Phone"},
					},
					"password": &memdb.IndexSchema{
						Name:    "password",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Password"},
					},
					"email": &memdb.IndexSchema{
						Name:    "email",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"userstatus": &memdb.IndexSchema{
						Name:    "userstatus",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "UserStatus"},
					},
				},
			},
		},
	}
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	return db, err
}
