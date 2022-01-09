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
	Flyable(ctx context.Context, characterID uint64) ([]*skillz.CharacterFlyableShip, error)
	SkillQueue(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillQueue, error)
	SkillsGrouped(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillGroup, error)
}

type Service struct {
	logger *logrus.Logger
	cache  cache.SkillAPI
	esi    esi.SkillAPI

	universe universe.API

	skills skillz.CharacterSkillRepository

	scopes []skillz.Scope
}

var _ API = (*Service)(nil)

func New(logger *logrus.Logger, cache cache.SkillAPI, esi esi.SkillAPI, universe universe.API, skills skillz.CharacterSkillRepository) *Service {
	return &Service{
		logger:   logger,
		cache:    cache,
		esi:      esi,
		universe: universe,
		skills:   skills,

		scopes: []skillz.Scope{skillz.ReadSkillsV1, skillz.ReadSkillQueueV1},
	}
}

func (s *Service) Process(ctx context.Context, user *skillz.User) error {

	var err error
	var funcs = []func(context.Context, *skillz.User) error{
		s.updateSkills, s.updateAttributes, s.updateSkillQueue,
	}

	for _, f := range funcs {
		err = f(ctx, user)
		if err != nil {
			break
		}

		time.Sleep(time.Second)
	}

	return err

}

func (s *Service) Scopes() []skillz.Scope {
	return s.scopes
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

	return meta, s.cache.SetCharacterSkillMeta(ctx, meta, time.Hour)

}

func (s *Service) Flyable(ctx context.Context, characterID uint64) ([]*skillz.CharacterFlyableShip, error) {

	flyable, err := s.cache.CharacterFlyableShips(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if flyable != nil {
		return flyable, nil
	}

	flyable, err = s.skills.CharacterFlyableShips(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character skills from data store")
	}

	for _, ship := range flyable {

		shipType, err := s.universe.Type(ctx, ship.ShipTypeID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch ship info")
		}

		shipGroup, err := s.universe.Group(ctx, shipType.GroupID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch ship group info")
		}

		shipType.Group = shipGroup
		ship.Ship = shipType

	}

	return flyable, s.cache.SetCharacterFlyableShips(ctx, characterID, flyable, time.Hour)

}

func (s *Service) Skillz(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkill, error) {

	skillz, err := s.cache.CharacterSkills(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(skillz) > 0 && err == nil {
		return skillz, err
	}

	skillz, err = s.skills.CharacterSkills(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return skillz, s.cache.SetCharacterSkills(ctx, characterID, skillz, time.Hour)

}

func (s *Service) SkillsGrouped(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillGroup, error) {

	groupedSkillz, err := s.cache.CharacterGroupedSkillz(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(groupedSkillz) > 0 && err == nil {
		return groupedSkillz, err
	}

	skills, err := s.Skillz(ctx, characterID)
	if err != nil {
		return nil, err
	}

	skillzMap := make(map[uint]*skillz.CharacterSkill)
	for _, skill := range skills {
		skillzMap[skill.SkillID] = skill
	}

	groups, err := s.universe.GroupsByCategory(ctx, 16)
	if err != nil {
		return nil, err
	}

	groupedSkillz = make([]*skillz.CharacterSkillGroup, 0, len(groups))

	for _, group := range groups {
		group.Types, err = s.universe.TypesByGroup(ctx, group.ID)
		if err != nil {
			return nil, err
		}

		groupedSkill := &skillz.CharacterSkillGroup{
			Info:         group,
			Skills:       make([]*skillz.CharacterSkill, 0, len(group.Types)),
			TotalGroupSP: 0,
		}

		for _, t := range group.Types {
			skill, ok := skillzMap[t.ID]
			if !ok {
				continue
			}

			skill.Info = t

			groupedSkill.Skills = append(groupedSkill.Skills, skill)
			groupedSkill.TotalGroupSP += skill.SkillpointsInSkill
		}

		groupedSkillz = append(groupedSkillz, groupedSkill)
	}

	defer func() {
		err = s.cache.SetCharacterGroupedSkillz(ctx, characterID, groupedSkillz, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache grouped skillz")
		}
	}()

	return groupedSkillz, nil
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

		err = s.cache.SetCharacterSkills(ctx, updateSkills.CharacterID, updateSkills.Skills, time.Hour)
		if err != nil {
			return errors.Wrap(err, "failed to cache character skills")
		}

		updateSkills.Skills = nil

		err = s.cache.SetCharacterSkillMeta(ctx, updateSkills, time.Hour)
		if err != nil {
			return errors.Wrap(err, "failed to cache character skill meta")
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

		err = s.cache.SetCharacterAttributes(ctx, updatedAttributes, time.Hour)
		if err != nil {
			return errors.Wrap(err, "failed to cache character attributes")
		}

	}

	return nil

}

func (s *Service) SkillQueue(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillQueue, error) {

	queue, err := s.cache.CharacterSkillQueue(ctx, characterID)
	if err != nil {
		return nil, err
	}

	if queue != nil {
		return queue, nil
	}

	queue, err = s.skills.CharacterSkillQueue(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character skill queue from data store")
	}

	for _, position := range queue {

		positionType, err := s.universe.Type(ctx, position.SkillID)
		if err != nil {
			return queue, errors.Wrap(err, "failed to fetch skill info for queue position")
		}

		positionGroup, err := s.universe.Group(ctx, positionType.GroupID)
		if err != nil {
			return queue, errors.Wrap(err, "failed to fetch skill group info for queue position")
		}

		positionType.Group = positionGroup
		position.Type = positionType

	}

	return queue, s.cache.SetCharacterSkillQueue(ctx, characterID, queue, time.Hour)

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

			for _, position := range updatedQueue {
				t, err := s.universe.Type(ctx, position.SkillID)
				if err != nil {
					s.logger.WithError(err).Error("failed to fetch type for skill queue position")
				}

				position.Type = t
			}

			err = s.cache.SetCharacterSkillQueue(ctx, user.CharacterID, updatedQueue, time.Hour)
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

	groups, err := s.universe.GroupsByCategory(ctx, 6)
	if err != nil {
		return errors.Wrap(err, "failed to fetch ship groups by category")
	}

	shipData := make([]*skillz.Type, 0)
	for _, group := range groups {

		groupShips, err := s.universe.TypesByGroup(ctx, group.ID)
		if err != nil {
			return errors.Wrap(err, "failed to fetch types by group id")
		}

		for _, ship := range groupShips {

			dogma, err := s.universe.TypeAttributes(ctx, ship.ID)
			if err != nil {
				return errors.Wrap(err, "failed to fetch attributes for type")
			}

			ship.Attributes = dogma
			ship.Group = group
			shipData = append(shipData, ship)
		}

	}

	flyableShips := make([]*skillz.CharacterFlyableShip, 0, len(shipData))

OUTER:
	for _, ship := range shipData {
		mapShipDogma := make(map[uint]*skillz.TypeDogmaAttribute)
		for _, attribute := range ship.Attributes {
			mapShipDogma[attribute.AttributeID] = attribute
		}

		flyable := &skillz.CharacterFlyableShip{
			CharacterID: user.CharacterID,
			ShipTypeID:  ship.ID,
			Ship:        ship,
		}

		for _, nameAttributeID := range skillNameDogmaSlice {
			if _, ok := mapShipDogma[nameAttributeID]; !ok {
				// Skill Level Name Attribute for the level is missing, can break and save as flyable
				break
			}

			if _, ok := mapShipDogma[skillNameToLevelDogmaMap[nameAttributeID]]; !ok {
				// Skill Level Attribute is missing, this is an error,
				// log the missing attribute, continue outer loop
				// Save that the ship is not flyable
				flyableShips = append(flyableShips, flyable)
				continue OUTER
			}

			skillID := uint(mapShipDogma[nameAttributeID].Value)
			level := uint(mapShipDogma[skillNameToLevelDogmaMap[nameAttributeID]].Value)

			skill := skillFromSkillSlice(skillID, skills)
			if skill == nil {
				flyableShips = append(flyableShips, flyable)
				continue OUTER
			}

			if skill.TrainedSkillLevel < level {
				flyableShips = append(flyableShips, flyable)
				continue OUTER
			}

		}

		flyable.Flyable = true
		flyableShips = append(flyableShips, flyable)
	}

	err = s.skills.CreateCharacterFlyableShips(ctx, flyableShips)
	if err != nil {
		return errors.Wrap(err, "failed to save flyable ships to data store")
	}

	return s.cache.SetCharacterFlyableShips(ctx, user.CharacterID, flyableShips, time.Hour)

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
