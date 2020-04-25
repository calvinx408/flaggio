package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
)

func TestVariantRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// drop database first
	if err := mongoDB.Drop(ctx); err != nil {
		t.Fatalf("failed drop database: %s", err)
	}

	// create new repo
	flgRepo, err := mongo_repo.NewFlagRepository(ctx, mongoDB)
	assert.NoError(t, err, "failed to create flag repository")
	repo := mongo_repo.NewVariantRepository(flgRepo.(*mongo_repo.FlagRepository))

	// create a flag
	flgID, err := flgRepo.Create(ctx, flaggio.NewFlag{Key: "test"})
	assert.NoError(t, err, "failed to create flag")

	// create the first variant
	vrnt1ID, err := repo.Create(ctx, flgID, flaggio.NewVariant{Value: 2.1})
	assert.NoError(t, err, "failed to create first variant")

	// checks the variant was created
	vrnt, err := repo.FindByID(ctx, flgID, vrnt1ID)
	assert.NoError(t, err, "failed to find first variant")
	assert.Equal(t, &flaggio.Variant{ID: vrnt1ID, Value: 2.1}, vrnt)

	// create the second variant
	vrnt2ID, err := repo.Create(ctx, flgID, flaggio.NewVariant{Value: "a"})
	assert.NoError(t, err, "failed to create second variant")

	// find the created variant
	vrnt, err = repo.FindByID(ctx, flgID, vrnt2ID)
	assert.NoError(t, err, "failed to find second variant")
	assert.Equal(t, &flaggio.Variant{ID: vrnt2ID, Value: "a"}, vrnt)

	// update the second variant
	err = repo.Update(ctx, flgID, vrnt2ID, flaggio.UpdateVariant{Value: false})
	assert.NoError(t, err, "failed to update second variant")

	// find second variant
	vrnt, err = repo.FindByID(ctx, flgID, vrnt2ID)
	assert.NoError(t, err, "failed to find second variant again")
	assert.Equal(t, &flaggio.Variant{ID: vrnt2ID, Value: false}, vrnt)

	// delete the first variant
	err = repo.Delete(ctx, flgID, vrnt1ID)
	assert.NoError(t, err, "failed to delete first variant")

	// find first variant
	vrnt, err = repo.FindByID(ctx, flgID, vrnt1ID)
	assert.EqualError(t, err, "variant: not found")
	assert.Nil(t, vrnt)
}
