package skill

import (
	"context"
	"database/sql"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/universe"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
)

type API interface {
	skillz.Processor
	Meta(ctx context.Context, characterID uint64) (*skillz.CharacterSkillMeta, error)
	Skillz(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkill, error)
	Attributes(ctx context.Context, characterID uint64) (*skillz.CharacterAttributes, error)
	Flyable(ctx context.Context, characterID uint64) ([]*skillz.ShipGroup, error)
	SkillQueue(ctx context.Context, characterID uint64) (*skillz.CharacterSkillQueueSummary, error)
	SkillsGrouped(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillGroup, error)
}

type Service struct {
	logger *logrus.Logger
	cache  cache.SkillAPI
	esi    esi.SkillAPI

	universe universe.API

	skills skillz.CharacterSkillRepository
}

var _ API = (*Service)(nil)

func New(logger *logrus.Logger, cache cache.SkillAPI, esi esi.SkillAPI, universe universe.API, skills skillz.CharacterSkillRepository) *Service {
	return &Service{
		logger:   logger,
		cache:    cache,
		esi:      esi,
		universe: universe,
		skills:   skills,
	}
}

func (s *Service) Process(ctx context.Context, user *skillz.User) error {

	var err error
	var funcs = []func(context.Context, *skillz.User) error{}
	for _, scope := range user.Scopes {
		switch scope {
		case skillz.ReadSkillsV1:
			funcs = append(funcs, s.updateSkills, s.updateAttributes)
		case skillz.ReadSkillQueueV1:
			funcs = append(funcs, s.updateSkillQueue)
		}
	}

	for _, f := range funcs {
		err = f(ctx, user)
		if err != nil {
			s.logger.WithError(err).Error("encounted error executing processor")
			break
		}
	}

	return err

}

func (s *Service) Meta(ctx context.Context, characterID uint64) (*skillz.CharacterSkillMeta, error) {

	meta, err := s.cache.CharacterSkillMeta(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if meta != nil {
		return meta, nil
	}

	meta, err = s.skills.CharacterSkillMeta(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character skills from data store")
	}

	defer func() {
		err = s.cache.SetCharacterSkillMeta(ctx, meta, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache character skill meta")
		}
	}()

	return meta, nil

}

func (s *Service) Flyable(ctx context.Context, characterID uint64) ([]*skillz.ShipGroup, error) {

	flyableShipGroups, err := s.cache.CharacterFlyableShips(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(flyableShipGroups) > 0 {
		return flyableShipGroups, nil
	}

	flyable, err := s.skills.CharacterFlyableShips(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character skills from data store")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	flyableShipsMap := make(map[uint]struct{})
	for _, entry := range flyable {
		flyableShipsMap[entry.ShipTypeID] = struct{}{}
	}

	groups, err := s.universe.TypeGroupsHydrated(ctx, universe.CategoryShips)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch hydrated ship groups")
	}

	shipGroups := make([]*skillz.ShipGroup, 0, len(groups))
	for _, group := range groups {
		shipGroup := &skillz.ShipGroup{Group: group, Ships: make([]*skillz.ShipType, 0, len(group.Types))}
		for _, ship := range group.Types {
			shipType := &skillz.ShipType{Type: ship, Flyable: false}
			if _, ok := flyableShipsMap[ship.ID]; ok {
				shipType.Flyable = true
			}
			shipGroup.Ships = append(shipGroup.Ships, shipType)
		}
		shipGroups = append(shipGroups, shipGroup)
	}

	defer func() {
		err := s.cache.SetCharacterFlyableShips(ctx, characterID, shipGroups, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache character flyable ships")
		}
	}()

	return shipGroups, nil

}

func (s *Service) Skillz(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkill, error) {

	skills, err := s.cache.CharacterSkills(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(skills) > 0 && err == nil {
		return skills, err
	}

	skills, err = s.skills.CharacterSkills(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	skillInfo, err := s.universe.SkillTypesHydrated(ctx)
	if err != nil {
		return nil, err
	}

	mapSkillInfo := make(map[uint]*skillz.Type)
	for _, info := range skillInfo {
		mapSkillInfo[info.ID] = info
	}

	for _, skill := range skills {
		if _, ok := mapSkillInfo[skill.SkillID]; !ok {
			continue
		}

		skill.Info = mapSkillInfo[skill.SkillID]
	}

	defer func(characterID uint64, skills []*skillz.CharacterSkill) {
		err = s.cache.SetCharacterSkills(ctx, characterID, skills, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache character skillz")
		}
	}(characterID, skills)

	return skills, nil

}

func (s *Service) SkillsGrouped(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillGroup, error) {

	groupedSkillz, err := s.cache.CharacterGroupedSkillz(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(groupedSkillz) > 0 {
		return groupedSkillz, err
	}

	groups, err := s.universe.TypeGroupsHydrated(ctx, universe.CategorySkills)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch hydrated skill groups from universe")
	}

	skills, err := s.Skillz(ctx, characterID)
	if err != nil {
		return nil, err
	}

	skillzMap := make(map[uint]*skillz.CharacterSkill)
	for _, skill := range skills {
		skillzMap[skill.SkillID] = skill
	}

	var characterSkillGroups = make([]*skillz.CharacterSkillGroup, 0, 50)
	for _, group := range groups {
		skillGroup := &skillz.SkillGroup{Group: group}
		characterSkillGroup := &skillz.CharacterSkillGroup{SkillGroup: skillGroup}
		skillGroup.Skills = make([]*skillz.SkillType, 0, len(group.Types))
		for _, t := range group.Types {
			skillType := &skillz.SkillType{Type: t}

			for _, attribute := range t.Attributes {
				if attribute.AttributeID == 275 {
					skillType.Rank = attribute
					break
				}
			}

			if skill, ok := skillzMap[t.ID]; ok {
				skillType.Skill = skill
				characterSkillGroup.TotalGroupSP += skill.SkillpointsInSkill
			}

			skillGroup.Skills = append(skillGroup.Skills, skillType)

		}

		characterSkillGroups = append(characterSkillGroups, characterSkillGroup)

	}

	defer func() {
		err = s.cache.SetCharacterGroupedSkillz(ctx, characterID, characterSkillGroups, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache grouped skillz")
		}
	}()

	return characterSkillGroups, nil
}

func (s *Service) updateSkills(ctx context.Context, user *skillz.User) error {

	s.logger.WithFields(logrus.Fields{
		"service": "skill",
		"userID":  user.ID.String(),
	}).Info("updating skills")

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterSkills, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return errors.Wrap(err, "failed to fetch tag for expiry check")
	}

	if etag != nil && etag.CachedUntil.Unix() > time.Now().Unix() {
		return nil
	}

	mods := s.esi.BaseCharacterModifiers(ctx, user, etagID, etag)

	updateSkills, err := s.esi.GetCharacterSkills(ctx, user.CharacterID, mods...)
	if err != nil {
		return errors.Wrap(err, "failed to fetch character skills from ESI")
	}

	if updateSkills != nil {
		err = s.skills.CreateCharacterSkillMeta(ctx, updateSkills)
		if err != nil {
			return errors.Wrap(err, "failed to update skill meta")
		}

		err = s.skills.CreateCharacterSkills(ctx, updateSkills.Skills)
		if err != nil {
			return errors.Wrap(err, "failed to update skills")
		}

	}

	return s.processFlyableShips(ctx, user)

}

func (s *Service) Attributes(ctx context.Context, characterID uint64) (*skillz.CharacterAttributes, error) {

	attributes, err := s.cache.CharacterAttributes(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if attributes != nil {
		return attributes, nil
	}

	attributes, err = s.skills.CharacterAttributes(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character attributes from data store")
	}

	defer func() {
		err = s.cache.SetCharacterAttributes(ctx, attributes, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache character attributes")
		}
	}()

	return attributes, nil

}

func (s *Service) updateAttributes(ctx context.Context, user *skillz.User) error {

	s.logger.WithFields(logrus.Fields{
		"service": "skill",
		"userID":  user.ID.String(),
	}).Info("updating attributes")

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterAttributes, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	if etag != nil && etag.CachedUntil.Unix() < time.Now().Unix() {
		return nil
	}

	mods := s.esi.BaseCharacterModifiers(ctx, user, etagID, etag)

	updatedAttributes, err := s.esi.GetCharacterAttributes(ctx, user.CharacterID, mods...)
	if err != nil {
		return errors.Wrap(err, "failed to fetch character attributes from ESI")
	}

	if updatedAttributes != nil {

		err = s.skills.CreateCharacterAttributes(ctx, updatedAttributes)
		if err != nil {
			return errors.Wrap(err, "failed to create/update character skill attributes")
		}

	}

	return nil

}

func (s *Service) SkillQueue(ctx context.Context, characterID uint64) (*skillz.CharacterSkillQueueSummary, error) {

	summary, err := s.cache.CharacterSkillQueueSummary(ctx, characterID)
	if err != nil {
		return nil, err
	}

	if summary != nil {
		return summary, nil
	}

	queue, err := s.skills.CharacterSkillQueue(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character skill queue from data store")
	}

	skillInfo, err := s.universe.SkillTypesHydrated(ctx)
	if err != nil {
		return nil, err
	}

	mapSkillInfo := make(map[uint]*skillz.Type)
	for _, info := range skillInfo {
		mapSkillInfo[info.ID] = info
	}

	mapGroupSummary := make(map[uint]*skillz.QueueGroupSummary)

	for _, position := range queue {
		if _, ok := mapSkillInfo[position.SkillID]; !ok {
			continue
		}

		t := mapSkillInfo[position.SkillID]
		position.Type = t

		var gs *skillz.QueueGroupSummary
		if _, ok := mapGroupSummary[t.GroupID]; !ok {
			mapGroupSummary[t.GroupID] = new(skillz.QueueGroupSummary)
		}
		gs = mapGroupSummary[t.GroupID]
		if t.Group != nil {
			gs.Group = t.Group
		}

		gs.Count += 1
		if position.LevelEndSp.Valid && position.LevelStartSp.Valid {
			gs.Skillpoints += position.LevelEndSp.Uint - position.LevelStartSp.Uint
		}
		if position.FinishDate.Valid && position.StartDate.Valid {
			gs.Duration += time.Duration(position.FinishDate.Time.Unix() - position.StartDate.Time.Unix())
		}
	}

	summary = new(skillz.CharacterSkillQueueSummary)
	summary.Summary = make([]*skillz.QueueGroupSummary, 0, len(mapGroupSummary))
	for _, entry := range mapGroupSummary {
		summary.Summary = append(summary.Summary, entry)
	}
	summary.Queue = queue

	if len(queue) > 0 {
		defer func(summary *skillz.CharacterSkillQueueSummary) {
			err = s.cache.SetCharacterSkillQueueSummary(ctx, characterID, summary, time.Hour)
			if err != nil {
				s.logger.WithError(err).Error("failed to character skill queue")
			}
		}(summary)
	}

	return summary, nil

}

func (s *Service) updateSkillQueue(ctx context.Context, user *skillz.User) error {

	s.logger.WithFields(logrus.Fields{
		"service": "skill",
		"userID":  user.ID.String(),
	}).Info("updating skill queue")

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterSkillQueue, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	if etag != nil && etag.CachedUntil.Unix() > time.Now().Unix() {
		return nil
	}

	mods := s.esi.BaseCharacterModifiers(ctx, user, etagID, etag)

	updatedQueue, err := s.esi.GetCharacterSkillQueue(ctx, user.CharacterID, mods...)
	if err != nil {
		return errors.Wrap(err, "failed to fetch character skill queue from ESI")
	}

	if updatedQueue != nil {

		err = s.skills.DeleteCharacterSkillQueue(ctx, user.CharacterID)
		if err != nil {
			return errors.Wrap(err, "failed to delete character skill queue")
		}

		if len(updatedQueue) > 0 {

			err = s.skills.CreateCharacterSkillQueue(ctx, updatedQueue)
			if err != nil {
				return errors.Wrap(err, "failed to create character skill queue")
			}

		}
	}

	return nil

}

func (s *Service) processFlyableShips(ctx context.Context, user *skillz.User) error {

	s.logger.WithFields(logrus.Fields{
		"service": "skill",
		"userID":  user.ID.String(),
	}).Info("updating flyable ships")

	err := s.skills.DeleteCharacterFlyableShips(ctx, user.CharacterID)
	if err != nil {
		return errors.Wrap(err, "failed to remove character flyables")
	}

	skills, err := s.Skillz(ctx, user.CharacterID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch character skillz")
	}

	skillzMap := make(map[uint]*skillz.CharacterSkill)
	for _, skill := range skills {
		skillzMap[skill.SkillID] = skill
	}

	groups, err := s.universe.TypeGroupsHydrated(ctx, universe.CategoryShips)
	if err != nil {
		return errors.Wrap(err, "failed to fetch ship groups by category")
	}

	flyableShips := make([]*skillz.CharacterFlyableShip, 0, 550)

	for _, group := range groups {
		shipGroup := &skillz.ShipGroup{
			Group: group,
			Ships: make([]*skillz.ShipType, 0, len(group.Types)),
		}

	SHIPLOOP:
		for _, ship := range group.Types {
			mapShipDogma := make(map[uint]*skillz.TypeDogmaAttribute)
			for _, attribute := range ship.Attributes {
				mapShipDogma[attribute.AttributeID] = attribute
			}

			shipType := &skillz.ShipType{
				Type: ship,
			}

			shipGroup.Ships = append(shipGroup.Ships, shipType)

			for _, nameAttributeID := range skillNameDogmaSlice {
				// Check to see if this ships dogma attributes contain an entry for this dogma attribute
				// If it doesn't, then that means that this ship does not have this level of skill requirements
				// We can break and say that the character is able to fly this ship.
				// Human translation of this is the ship might have a single required skill
				// This will fail when looking for a second required skill
				if _, ok := mapShipDogma[nameAttributeID]; !ok {
					break
				}

				// Okay, so we've confirm that this ship has this level of skill requirements. Now we need
				// to find the level of the skill that is required. This shouldn't be missing, but if it is,
				// then we cannot properly determine if the ship is flyable or not, so we will continue  the
				// outer loop to the next ship
				if _, ok := mapShipDogma[skillNameToLevelDogmaMap[nameAttributeID]]; !ok {
					continue SHIPLOOP
				}

				// This turns the dogma value into a comparable uint
				skillID := uint(mapShipDogma[nameAttributeID].Value)
				level := uint(mapShipDogma[skillNameToLevelDogmaMap[nameAttributeID]].Value)

				// Check to see if this character has the required skill
				skill := skillFromSkillSlice(skillID, skills)
				// Character does not have the skill, ship is not flytable, continue to next ship
				if skill == nil {
					continue SHIPLOOP
				}

				// Character has previously trained the skill to the required level.
				// Since we don't care about Omega status exactly, we can look at the "Omega" skill level
				// as opposed to the "Alpha" skil level
				if skill.TrainedSkillLevel < level {
					continue SHIPLOOP
				}

			}

			flyableShips = append(flyableShips, &skillz.CharacterFlyableShip{
				CharacterID: user.CharacterID,
				ShipTypeID:  ship.ID,
			})
			shipType.Flyable = true

		}

	}

	if len(flyableShips) > 0 {
		err = s.skills.CreateCharacterFlyableShips(ctx, flyableShips)
		if err != nil {
			return errors.Wrap(err, "failed to save flyable ships to data store")
		}

	}
	return nil

}

var skillNameDogmaSlice = []uint{182, 183, 184, 1285, 1289, 1290}

var skillNameToLevelDogmaMap = map[uint]uint{
	182:  277,
	183:  278,
	184:  279,
	1285: 1286,
	1289: 1287,
	1290: 1288,
}

func skillFromSkillSlice(skillID uint, skills []*skillz.CharacterSkill) *skillz.CharacterSkill {
	for _, skill := range skills {
		if skill.SkillID == skillID {
			return skill
		}
	}

	return nil

}

func PlushHelperHasSkillCount(group *skillz.SkillGroup) uint {
	var out = uint(0)

	for _, t := range group.Skills {
		if t == nil {
			continue
		}

		if t.Skill != nil {
			out++
		}
	}
	return out
}

func PlushHelperMissingSkill(group *skillz.SkillGroup) uint {
	var out = uint(0)
	for _, t := range group.Skills {
		if t.Skill == nil {
			out++
		}
	}
	return out
}

func PlushHelperLevelVSkillCount(group *skillz.SkillGroup) uint {
	var out = uint(0)

	for _, t := range group.Skills {
		if t.Skill != nil && t.Skill.TrainedSkillLevel == 5 {
			out++
		}
	}

	return out
}

func PlushHelperLevelVSPTotal(group *skillz.SkillGroup) uint {
	var out = uint(0)

	for _, t := range group.Skills {
		if t.Skill != nil && t.Skill.TrainedSkillLevel == 5 {
			out += t.Skill.SkillpointsInSkill
		}
	}

	return out
}

func PlushHelperPossibleSPByRank(rank interface{}) float64 {

	switch rank := rank.(type) {
	case float64:
		return float64(rank * 256000)
	case int:
		return float64(rank * 256000)

	}

	return 0
}
