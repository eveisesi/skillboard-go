package skill

import (
	"context"
	"database/sql"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
)

type API interface {
	skillz.Processor
	Meta(ctx context.Context, characterID uint64) (*skillz.CharacterSkillMeta, error)
	Skillz(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkill, error)
	Attributes(ctx context.Context, characterID uint64) (*skillz.CharacterAttributes, error)
	SkillQueue(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillQueue, error)
}

type Service struct {
	cache cache.SkillAPI
	esi   esi.SkillAPI

	skills skillz.CharacterSkillRepository

	scopes []skillz.Scope
}

var _ API = (*Service)(nil)

func New(cache cache.SkillAPI, esi esi.SkillAPI, skills skillz.CharacterSkillRepository) *Service {
	return &Service{
		cache:  cache,
		esi:    esi,
		skills: skills,

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

	meta, err = s.skills.CharacterSkillMeta(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character skills from data store")
	}

	return meta, s.cache.SetCharacterSkillMeta(ctx, meta, time.Hour)

}

func (s *Service) Skillz(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkill, error) {

	skillz, err := s.cache.CharacterSkills(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(skillz) > 0 || err == nil {
		return skillz, err
	}

	skillz, err = s.skills.CharacterSkills(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return skillz, err

}

func (s *Service) updateSkills(ctx context.Context, user *skillz.User) error {

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

	return nil
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

	return queue, nil

}

func (s *Service) updateSkillQueue(ctx context.Context, user *skillz.User) error {

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

		err = s.skills.CreateCharacterSkillQueue(ctx, updatedQueue)
		if err != nil {
			return errors.Wrap(err, "failed to create character skill queue")
		}

		err = s.cache.SetCharacterSkillQueue(ctx, user.CharacterID, updatedQueue, time.Hour)
		if err != nil {
			return errors.Wrap(err, "failed to create character skill queue")
		}
	}

	return nil

}
