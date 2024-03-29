package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type SkillAPI interface {
	CharacterSkillMeta(ctx context.Context, characterID uint64) (*skillz.CharacterSkillMeta, error)
	SetCharacterSkillMeta(ctx context.Context, meta *skillz.CharacterSkillMeta, expires time.Duration) error
	CharacterAttributes(ctx context.Context, characterID uint64) (*skillz.CharacterAttributes, error)
	SetCharacterAttributes(ctx context.Context, meta *skillz.CharacterAttributes, expires time.Duration) error
	CharacterSkills(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkill, error)
	SetCharacterSkills(ctx context.Context, characterID uint64, skills []*skillz.CharacterSkill, expires time.Duration) error
	CharacterGroupedSkillz(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillGroup, error)
	SetCharacterGroupedSkillz(ctx context.Context, characterID uint64, groups []*skillz.CharacterSkillGroup, expires time.Duration) error
	// CharacterSkillQueue(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillQueue, error)
	// SetCharacterSkillQueue(ctx context.Context, characterID uint64, positions []*skillz.CharacterSkillQueue, expires time.Duration) error
	CharacterSkillQueueSummary(ctx context.Context, characterID uint64) (*skillz.CharacterSkillQueueSummary, error)
	SetCharacterSkillQueueSummary(ctx context.Context, characterID uint64, summary *skillz.CharacterSkillQueueSummary, expires time.Duration) error
	CharacterFlyableShips(ctx context.Context, characterID uint64) ([]*skillz.ShipGroup, error)
	SetCharacterFlyableShips(ctx context.Context, characterID uint64, flyable []*skillz.ShipGroup, expires time.Duration) error
}

const (
	characterAttributesKeyPrefix        = "character::attributes"
	characterSkillMetaKeyPrefix         = "character::skill::meta"
	characterSkillsKeyPrefix            = "character::skills"
	characterSkillsGroupedKeyPrefix     = "character::skills-grouped"
	characterFlyableKeyPrefix           = "character::flyable"
	characterSkillQueueKeySummaryPrefix = "character::skillqueue::summary"
)

func (s *Service) CharacterSkillMeta(ctx context.Context, characterID uint64) (*skillz.CharacterSkillMeta, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(characterSkillMetaKeyPrefix, strconv.FormatUint(characterID, 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkillMeta", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var meta = new(skillz.CharacterSkillMeta)
	err = json.Unmarshal(result, meta)
	return meta, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkillMeta", "failed to decode json to structure")

}

func (s *Service) SetCharacterSkillMeta(ctx context.Context, meta *skillz.CharacterSkillMeta, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	data, err := json.Marshal(meta)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkillMeta", "failed to encode struct as json")
	}

	key := generateKey(characterSkillMetaKeyPrefix, strconv.FormatUint(meta.CharacterID, 10))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkillMeta", "failed to write cache")

}

func (s *Service) CharacterAttributes(ctx context.Context, characterID uint64) (*skillz.CharacterAttributes, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(characterAttributesKeyPrefix, strconv.FormatUint(characterID, 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterAttributes", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var meta = new(skillz.CharacterAttributes)
	err = json.Unmarshal(result, meta)
	return meta, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterAttributes", "failed to decode json to structure")

}

func (s *Service) SetCharacterAttributes(ctx context.Context, meta *skillz.CharacterAttributes, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	data, err := json.Marshal(meta)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterAttributes", "failed to encode struct as json")
	}

	key := generateKey(characterAttributesKeyPrefix, strconv.FormatUint(meta.CharacterID, 10))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterAttributes", "failed to write cache")

}

func (s *Service) CharacterSkills(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkill, error) {
	if s.disabled {
		return nil, nil
	}
	var skills = make([]*skillz.CharacterSkill, 0)

	key := generateKey(characterSkillsKeyPrefix, strconv.FormatUint(characterID, 10))
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return skills, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkills", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return skills, nil
	}

	skills = make([]*skillz.CharacterSkill, 0, len(results))

	for _, result := range results {
		var skill = new(skillz.CharacterSkill)
		err = json.Unmarshal([]byte(result), skill)
		if err != nil {
			return skills, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkills", "failed to decode json to structure")
		}

		skills = append(skills, skill)
	}

	return skills, nil

}

func (s *Service) SetCharacterSkills(ctx context.Context, characterID uint64, skills []*skillz.CharacterSkill, expires time.Duration) error {
	if s.disabled {
		return nil
	}

	members := make([]interface{}, 0, len(skills))
	for _, skill := range skills {
		data, err := json.Marshal(skill)
		if err != nil {
			return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkills", "failed to encode struct as json")
		}

		members = append(members, string(data))
	}

	if len(members) == 0 {
		return nil
	}

	key := generateKey(characterSkillsKeyPrefix, strconv.FormatUint(characterID, 10))
	err := s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkills", "failed to write cache")
	}

	err = s.redis.Expire(ctx, key, expires).Err()

	return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkills", "failed to set expiry on set")

}

func (s *Service) CharacterGroupedSkillz(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillGroup, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(characterSkillsGroupedKeyPrefix, strconv.FormatUint(characterID, 10))
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterGroupedSkillz", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return make([]*skillz.CharacterSkillGroup, 0), nil
	}

	groups := make([]*skillz.CharacterSkillGroup, 0, len(results))
	for _, result := range results {
		var group = new(skillz.CharacterSkillGroup)
		err = json.Unmarshal([]byte(result), group)
		if err != nil {
			return groups, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterGroupedSkillz", "failed to decode json to structure")
		}

		groups = append(groups, group)
	}

	return groups, nil

}

func (s *Service) SetCharacterGroupedSkillz(ctx context.Context, characterID uint64, groups []*skillz.CharacterSkillGroup, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	members := make([]interface{}, 0, len(groups))
	for _, group := range groups {
		data, err := json.Marshal(group)
		if err != nil {
			return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterGroupedSkillz", "failed to encode struct as json")
		}
		members = append(members, string(data))
	}

	if len(members) == 0 {
		return nil
	}

	key := generateKey(characterSkillsGroupedKeyPrefix, strconv.FormatUint(characterID, 10))
	err := s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterGroupedSkillz", "failed to write cache")
	}

	err = s.redis.Expire(ctx, key, expires).Err()

	return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterGroupedSkillz", "failed to set expiry on set")

}

func (s *Service) CharacterSkillQueueSummary(ctx context.Context, characterID uint64) (*skillz.CharacterSkillQueueSummary, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(characterSkillQueueKeySummaryPrefix, strconv.FormatUint(characterID, 10))
	result, err := s.redis.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkillQueue", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var summary = new(skillz.CharacterSkillQueueSummary)

	err = json.Unmarshal([]byte(result), summary)
	if err != nil {
		return summary, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkillQueue", "failed to decode json to structure")
	}

	return summary, nil

}

func (s *Service) SetCharacterSkillQueueSummary(ctx context.Context, characterID uint64, summary *skillz.CharacterSkillQueueSummary, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	data, err := json.Marshal(summary)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkillQueueSummary", "failed to encode struct as json")
	}

	key := generateKey(characterSkillQueueKeySummaryPrefix, strconv.FormatUint(characterID, 10))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkillQueueSummary", "failed to write cache")

}

func (s *Service) CharacterFlyableShips(ctx context.Context, characterID uint64) ([]*skillz.ShipGroup, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(characterFlyableKeyPrefix, strconv.FormatUint(characterID, 10))
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkills", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return make([]*skillz.ShipGroup, 0), nil
	}

	flyables := make([]*skillz.ShipGroup, 0, len(results))
	for _, result := range results {
		var flyable = new(skillz.ShipGroup)
		err = json.Unmarshal([]byte(result), flyable)
		if err != nil {
			return flyables, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkills", "failed to decode json to structure")
		}

		flyables = append(flyables, flyable)
	}

	return flyables, nil

}

func (s *Service) SetCharacterFlyableShips(ctx context.Context, characterID uint64, flyable []*skillz.ShipGroup, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	members := make([]interface{}, 0, len(flyable))
	for _, skill := range flyable {
		data, err := json.Marshal(skill)
		if err != nil {
			return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkills", "failed to encode struct as json")
		}

		members = append(members, string(data))
	}

	if len(members) == 0 {
		return nil
	}

	key := generateKey(characterFlyableKeyPrefix, strconv.FormatUint(characterID, 10))
	err := s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkills", "failed to write cache")
	}

	err = s.redis.Expire(ctx, key, expires).Err()

	return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkills", "failed to set expiry on set")

}
