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
	Skills(ctx context.Context, user *skillz.User) (*skillz.CharacterSkillMeta, error)
	Attributes(ctx context.Context, user *skillz.User) (*skillz.CharacterAttributes, error)
	SkillQueue(ctx context.Context, user *skillz.User) ([]*skillz.CharacterSkillQueue, error)
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

	_, err := s.Skills(ctx, user)
	if err != nil {
		return err
	}

	_, err = s.Attributes(ctx, user)
	if err != nil {
		return err
	}

	_, err = s.SkillQueue(ctx, user)
	if err != nil {
		return err
	}

	return nil

}

func (s *Service) Scopes() []skillz.Scope {
	return s.scopes
}

func (s *Service) Skills(ctx context.Context, user *skillz.User) (*skillz.CharacterSkillMeta, error) {

	meta, err := s.cache.CharacterSkillMeta(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if meta != nil {
		skills, err := s.cache.CharacterSkills(ctx, user.CharacterID)
		if err != nil {
			return nil, err
		}

		meta.Skills = skills

		return meta, err

	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterSkills, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch tag for expiry check")
	}

	meta, err = s.skills.CharacterSkillMeta(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character skills from data store")
	}

	exists := err == nil

	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {

		mods := append(make([]esi.ModifierFunc, 0, 3), s.esi.CacheEtag(ctx, etagID), s.esi.AddAuthorizationHeader(ctx, user.AccessToken))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updateSkills, err := s.esi.GetCharacterSkills(ctx, user.CharacterID, mods...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch character skills from ESI")
		}

		if updateSkills != nil {
			switch exists {
			case true:
				err = s.skills.UpdateCharacterSkillMeta(ctx, updateSkills)
				if err != nil {
					return nil, errors.Wrap(err, "failed to update skill meta")
				}

			case false:
				err = s.skills.CreateCharacterSkillMeta(ctx, updateSkills)
				if err != nil {
					return nil, errors.Wrap(err, "failed to update skill meta")
				}
			}

			err = s.skills.CreateCharacterSkills(ctx, updateSkills.Skills)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update skills")
			}

			meta = updateSkills
		}

	}

	skills := meta.Skills
	meta.Skills = nil
	err = s.cache.SetCharacterSkillMeta(ctx, meta, time.Hour)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cache character skill meta")
	}

	err = s.cache.SetCharacterSkills(ctx, meta.CharacterID, skills, time.Hour)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cache character skills")
	}

	meta.Skills = skills
	return meta, nil

}

func (s *Service) Attributes(ctx context.Context, user *skillz.User) (*skillz.CharacterAttributes, error) {

	attributes, err := s.cache.CharacterAttributes(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if attributes != nil {
		return attributes, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterAttributes, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	attributes, err = s.skills.CharacterAttributes(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character attributes from data store")
	}

	exists := err == nil

	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {

		mods := append(make([]esi.ModifierFunc, 0, 3), s.esi.CacheEtag(ctx, etagID), s.esi.AddAuthorizationHeader(ctx, user.AccessToken))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedAttributes, err := s.esi.GetCharacterAttributes(ctx, user.CharacterID, mods...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch character attributes from ESI")
		}

		if updatedAttributes != nil {

			err = s.skills.CreateCharacterAttributes(ctx, updatedAttributes)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create/update character skill attributes")
			}

			attributes = updatedAttributes

		}

	}

	return attributes, s.cache.SetCharacterAttributes(ctx, attributes, time.Hour)

}

func (s *Service) SkillQueue(ctx context.Context, user *skillz.User) ([]*skillz.CharacterSkillQueue, error) {

	queue, err := s.cache.CharacterSkillQueue(ctx, user.CharacterID)
	if err != nil {
		return nil, err
	}

	if queue != nil {
		return queue, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterSkillQueue, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	queue, err = s.skills.CharacterSkillQueue(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character skill queue from data store")
	}

	exists := err == nil

	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {

		mods := append(make([]esi.ModifierFunc, 0, 3), s.esi.CacheEtag(ctx, etagID), s.esi.AddAuthorizationHeader(ctx, user.AccessToken))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedQueue, err := s.esi.GetCharacterSkillQueue(ctx, user.CharacterID, mods...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch character skill queue from ESI")
		}

		if updatedQueue != nil {

			err = s.skills.DeleteCharacterSkillQueue(ctx, user.CharacterID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to delete character skill queue")
			}

			err = s.skills.CreateCharacterSkillQueue(ctx, updatedQueue)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create character skill queue")
			}

			queue = updatedQueue

		}

	}

	return queue, s.cache.SetCharacterSkillQueue(ctx, user.CharacterID, queue, time.Hour)

}
