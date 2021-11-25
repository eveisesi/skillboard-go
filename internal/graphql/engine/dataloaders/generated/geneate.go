//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 AllianceLoader uint *github.com/eveisesi/skillz.Alliance
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 CategoryLoader uint *github.com/eveisesi/skillz.Category
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 CharacterLoader uint64 *github.com/eveisesi/skillz.Character
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 CloneLoader *github.com/eveisesi/skillz.User *github.com/eveisesi/skillz.CharacterCloneMeta
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 ConstellationLoader uint *github.com/eveisesi/skillz.Constellation
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 CorporationLoader uint *github.com/eveisesi/skillz.Corporation
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 DeathCloneLoader *github.com/eveisesi/skillz.User *github.com/eveisesi/skillz.CharacterDeathClone
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 GroupLoader uint *github.com/eveisesi/skillz.Group
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 ImplantLoader *github.com/eveisesi/skillz.User []*github.com/eveisesi/skillz.CharacterImplant
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 JumpCloneLoader *github.com/eveisesi/skillz.User []*github.com/eveisesi/skillz.CharacterJumpClone
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 RegionLoader uint *github.com/eveisesi/skillz.Region
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 SolarSystemLoader uint *github.com/eveisesi/skillz.SolarSystem
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 StationLoader uint *github.com/eveisesi/skillz.Station
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 StructureLoader uint64 *github.com/eveisesi/skillz.Structure
//go:generate go run github.com/ddouglas/dataloaden@v0.4.0 TypeLoader uint *github.com/eveisesi/skillz.Type

package generated
