package sqlx

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql/v2"
)

// ScanToStruct prepares a given set of Queries and assigns the resulting
// *sqlx.Stmt or *sqlx.NamedStmt statements to the fields of a given struct, matching based on the name
// in the `query` tag in the struct field names.
func ScanToStruct(obj interface{}, q goyesql.Queries, db *sqlx.DB) error {
	ob := reflect.ValueOf(obj)
	if ob.Kind() == reflect.Ptr {
		ob = ob.Elem()
	}

	if ob.Kind() != reflect.Struct {
		return fmt.Errorf("Failed to apply SQL statements to struct. Non struct type: %T", ob)
	}

	// Go through every field in the struct and look for it in the Args map.
	for i := 0; i < ob.NumField(); i++ {
		f := ob.Field(i)

		if f.IsValid() {
			if tag := ob.Type().Field(i).Tag.Get("query"); tag != "" && tag != "-" {
				// Extract the value of the `query` tag.
				var (
					tg   = strings.Split(tag, ",")
					name string
				)
				if len(tg) == 2 {
					if tg[0] != "-" && tg[0] != "" {
						name = tg[0]
					}
				} else {
					name = tg[0]
				}

				// Query name found in the field tag is not in the map.
				if _, ok := q[name]; !ok {
					return fmt.Errorf("query '%s' not found in query map", name)
				}

				if !f.CanSet() {
					return fmt.Errorf("query field '%s' is unexported", ob.Type().Field(i).Name)
				}

				switch f.Type().String() {
				case "string":
					// Unprepared SQL query.
					f.Set(reflect.ValueOf(q[name].Query))
				case "*sqlx.Stmt":
					// Prepared query.
					stmt, err := db.Preparex(q[name].Query)
					if err != nil {
						return fmt.Errorf("Error preparing query '%s': %v", name, err)
					}
					f.Set(reflect.ValueOf(stmt))
				case "*sqlx.NamedStmt":
					// Prepared query.
					stmt, err := db.PrepareNamed(q[name].Query)
					if err != nil {
						return fmt.Errorf("Error preparing query '%s': %v", name, err)
					}
					f.Set(reflect.ValueOf(stmt))
				}
			}
		}
	}

	return nil
}

// ScanWithContext prepares a given set of Queries and assigns the resulting
// *sqlx.Stmt or *sqlx.NamedStmt statement. It calls sqlx methods with context
// to pass metadata.
func ScanWithContext(ctx context.Context, obj interface{}, q goyesql.Queries, db *sqlx.DB) error {
	ob := reflect.ValueOf(obj)
	if ob.Kind() == reflect.Ptr {
		ob = ob.Elem()
	}

	if ob.Kind() != reflect.Struct {
		return fmt.Errorf("Failed to apply SQL statements to struct. Non struct type: %T", ob)
	}

	// Go through every field in the struct and look for it in the Args map.
	for i := 0; i < ob.NumField(); i++ {
		f := ob.Field(i)

		if f.IsValid() {
			if tag := ob.Type().Field(i).Tag.Get("query"); tag != "" && tag != "-" {
				// Extract the value of the `query` tag.
				var (
					tg   = strings.Split(tag, ",")
					name string
				)
				if len(tg) == 2 {
					if tg[0] != "-" && tg[0] != "" {
						name = tg[0]
					}
				} else {
					name = tg[0]
				}

				// Query name found in the field tag is not in the map.
				if _, ok := q[name]; !ok {
					return fmt.Errorf("query '%s' not found in query map", name)
				}

				if !f.CanSet() {
					return fmt.Errorf("query field '%s' is unexported", ob.Type().Field(i).Name)
				}

				switch f.Type().String() {
				case "string":
					// Unprepared SQL query.
					f.Set(reflect.ValueOf(q[name].Query))
				case "*sqlx.Stmt":
					// Prepared query.
					stmt, err := db.PreparexContext(ctx, q[name].Query)
					if err != nil {
						return fmt.Errorf("Error preparing query '%s': %v", name, err)
					}
					f.Set(reflect.ValueOf(stmt))
				case "*sqlx.NamedStmt":
					// Prepared query.
					stmt, err := db.PrepareNamedContext(ctx, q[name].Query)
					if err != nil {
						return fmt.Errorf("Error preparing query '%s': %v", name, err)
					}
					f.Set(reflect.ValueOf(stmt))
				}
			}
		}
	}

	return nil
}
