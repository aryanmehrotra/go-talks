package migrations

import (
	"context"

	"gofr.dev/pkg/gofr/migration"
)

func createTopics() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			if err := d.PubSub.CreateTopic(context.Background(), "orders"); err != nil {
				return err
			}
			return d.PubSub.CreateTopic(context.Background(), "alerts")
		},
	}
}
