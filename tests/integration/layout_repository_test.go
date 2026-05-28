//go:build integration

package integration

import (
	"testing"
	"wn/internal/domain/entity"
	layoutrepo "wn/internal/infrastructure/repository/layout"

	"github.com/google/uuid"
)

func TestLayoutRepository_CreateAndGet(t *testing.T) {
	resetDB(t)
	user := seedUser(t)
	repo := layoutrepo.NewRepository(env.DB)

	layout := &entity.Layout{
		Id:         uuid.New(),
		Title:      "My Board",
		OwnerId:    user.Id,
		HaveAccess: []uuid.UUID{user.Id},
		IsMain:     true,
		Color:      "#FFFFFF",
	}

	id, err := repo.CreateLayout(env.Ctx, layout)
	if err != nil {
		t.Fatalf("CreateLayout() error = %v", err)
	}
	if id != layout.Id {
		t.Errorf("id = %v, want %v", id, layout.Id)
	}

	layouts, err := repo.GetAvailableLayouts(env.Ctx, user.Id)
	if err != nil {
		t.Fatalf("GetAvailableLayouts() error = %v", err)
	}
	if len(layouts) != 1 {
		t.Fatalf("len = %d, want 1", len(layouts))
	}
	if layouts[0].Title != "My Board" {
		t.Errorf("Title = %q, want My Board", layouts[0].Title)
	}
}

func TestLayoutRepository_DeleteLayout(t *testing.T) {
	resetDB(t)
	user := seedUser(t)
	repo := layoutrepo.NewRepository(env.DB)

	layout := &entity.Layout{
		Id:         uuid.New(),
		Title:      "Temp",
		OwnerId:    user.Id,
		HaveAccess: []uuid.UUID{user.Id},
		Color:      "#000000",
	}
	if _, err := repo.CreateLayout(env.Ctx, layout); err != nil {
		t.Fatalf("CreateLayout() error = %v", err)
	}

	if err := repo.DeleteLayoutById(env.Ctx, layout.Id); err != nil {
		t.Fatalf("DeleteLayoutById() error = %v", err)
	}

	layouts, err := repo.GetAvailableLayouts(env.Ctx, user.Id)
	if err != nil {
		t.Fatalf("GetAvailableLayouts() error = %v", err)
	}
	if len(layouts) != 0 {
		t.Errorf("len = %d, want 0", len(layouts))
	}
}

func TestLayoutRepository_UpdateLayout(t *testing.T) {
	resetDB(t)
	user := seedUser(t)
	repo := layoutrepo.NewRepository(env.DB)

	layout := &entity.Layout{
		Id:         uuid.New(),
		Title:      "Old",
		OwnerId:    user.Id,
		HaveAccess: []uuid.UUID{user.Id},
		Color:      "#111111",
	}
	if _, err := repo.CreateLayout(env.Ctx, layout); err != nil {
		t.Fatalf("CreateLayout() error = %v", err)
	}

	rows, err := repo.UpdateLayout(env.Ctx, user.Id, layout.Id, "#222222", "New Title")
	if err != nil {
		t.Fatalf("UpdateLayout() error = %v", err)
	}
	if rows != 1 {
		t.Errorf("updated rows = %d, want 1", rows)
	}

	layouts, err := repo.GetAvailableLayouts(env.Ctx, user.Id)
	if err != nil {
		t.Fatalf("GetAvailableLayouts() error = %v", err)
	}
	if layouts[0].Title != "New Title" || layouts[0].Color != "#222222" {
		t.Errorf("layout = %+v", layouts[0])
	}
}
