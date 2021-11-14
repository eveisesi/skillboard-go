package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
)

func (r *characterImplantResolver) Implant(ctx context.Context, obj *skillz.CharacterImplant) (*skillz.Type, error) {
	return r.dataloaders.TypeLoader().Load(ctx, obj.ImplantID)
}

func (r *constellationResolver) Region(ctx context.Context, obj *skillz.Constellation) (*skillz.Region, error) {
	return r.dataloaders.RegionLoader().Load(ctx, obj.RegionID)
}

func (r *groupResolver) Category(ctx context.Context, obj *skillz.Group) (*skillz.Category, error) {
	return r.dataloaders.CategoryLoader().Load(ctx, obj.CategoryID)
}

func (r *solarSystemResolver) Constellation(ctx context.Context, obj *skillz.SolarSystem) (*skillz.Constellation, error) {
	return r.dataloaders.ConstellationLoader().Load(ctx, obj.ConstellationID)
}

func (r *stationResolver) System(ctx context.Context, obj *skillz.Station) (*skillz.SolarSystem, error) {
	return r.dataloaders.SolarSystemLoader().Load(ctx, obj.SystemID)
}

func (r *structureResolver) System(ctx context.Context, obj *skillz.Structure) (*skillz.SolarSystem, error) {
	return r.dataloaders.SolarSystemLoader().Load(ctx, obj.SolarSystemID)
}

func (r *typeResolver) Group(ctx context.Context, obj *skillz.Type) (*skillz.Group, error) {
	return r.dataloaders.GroupLoader().Load(ctx, obj.GroupID)
}

// Constellation returns generated.ConstellationResolver implementation.
func (r *Resolver) Constellation() generated.ConstellationResolver { return &constellationResolver{r} }

// Group returns generated.GroupResolver implementation.
func (r *Resolver) Group() generated.GroupResolver { return &groupResolver{r} }

// SolarSystem returns generated.SolarSystemResolver implementation.
func (r *Resolver) SolarSystem() generated.SolarSystemResolver { return &solarSystemResolver{r} }

// Station returns generated.StationResolver implementation.
func (r *Resolver) Station() generated.StationResolver { return &stationResolver{r} }

// Structure returns generated.StructureResolver implementation.
func (r *Resolver) Structure() generated.StructureResolver { return &structureResolver{r} }

// Type returns generated.TypeResolver implementation.
func (r *Resolver) Type() generated.TypeResolver { return &typeResolver{r} }

type constellationResolver struct{ *Resolver }
type groupResolver struct{ *Resolver }
type solarSystemResolver struct{ *Resolver }
type stationResolver struct{ *Resolver }
type structureResolver struct{ *Resolver }
type typeResolver struct{ *Resolver }
