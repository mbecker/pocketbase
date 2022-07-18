package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const (
	OAuthTokenCollectionName         = "oauthtokens"
	OAuthTokenCollectionNameProvider = "provider"
)

func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		// add up queries...
		// inserts the system profiles collection
		// -----------------------------------------------------------
		profileOwnerRule := fmt.Sprintf("%s = @request.user.id", models.ProfileCollectionUserFieldName)
		collection := &models.Collection{
			Name:       OAuthTokenCollectionName,
			System:     false,
			CreateRule: &profileOwnerRule,
			ListRule:   &profileOwnerRule,
			ViewRule:   &profileOwnerRule,
			UpdateRule: &profileOwnerRule,
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     OAuthTokenCollectionNameProvider,
					Type:     schema.FieldTypeText,
					Unique:   false,
					Required: true,
					System:   false,
				},
				&schema.SchemaField{
					Name:     models.ProfileCollectionUserFieldName,
					Type:     schema.FieldTypeUser,
					Unique:   false,
					Required: true,
					System:   false,
					Options: &schema.UserOptions{
						MaxSelect:     1,
						CascadeDelete: true,
					},
				},
				&schema.SchemaField{
					Name:     "access_token",
					Type:     schema.FieldTypeText,
					Options:  &schema.TextOptions{},
					Required: true,
				},
				&schema.SchemaField{
					Name:     "refresh_token",
					Type:     schema.FieldTypeText,
					Options:  &schema.TextOptions{},
					Required: true,
				},
				&schema.SchemaField{
					Name:     "token_type",
					Type:     schema.FieldTypeText,
					Options:  &schema.TextOptions{},
					Required: true,
				},
				&schema.SchemaField{
					Name:     "expiry",
					Type:     schema.FieldTypeDate,
					Options:  &schema.TextOptions{},
					Required: true,
				},
			),
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		tables := []string{
			OAuthTokenCollectionName,
		}

		for _, name := range tables {
			if _, err := db.DropTable(name).Execute(); err != nil {
				return err
			}
		}

		return nil
	})
}
