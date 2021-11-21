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
	CharacterSkillQueue(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillQueue, error)
	SetCharacterSkillQueue(ctx context.Context, characterID uint64, positions []*skillz.CharacterSkillQueue, expires time.Duration) error
}

const (
	characterAttributesKeyPrefix = "character::attributes"
	characterSkillMetaKeyPrefix  = "character::skill::meta"
	characterSkillsKeyPrefix     = "character::skills"
	characterSkillQueueKeyPrefix = "character::skillqueue"
)

func (s *Service) CharacterSkillMeta(ctx context.Context, characterID uint64) (*skillz.CharacterSkillMeta, error) {

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

	data, err := json.Marshal(meta)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkillMeta", "failed to encode struct as json")
	}

	key := generateKey(characterSkillMetaKeyPrefix, strconv.FormatUint(meta.CharacterID, 10))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkillMeta", "failed to write cache")

}

func (s *Service) CharacterAttributes(ctx context.Context, characterID uint64) (*skillz.CharacterAttributes, error) {

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

	data, err := json.Marshal(meta)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterAttributes", "failed to encode struct as json")
	}

	key := generateKey(characterAttributesKeyPrefix, strconv.FormatUint(meta.CharacterID, 10))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterAttributes", "failed to write cache")

}

func (s *Service) CharacterSkills(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkill, error) {

	key := generateKey(characterSkillsKeyPrefix, strconv.FormatUint(characterID, 10))
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkills", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return nil, nil
	}

	skills := make([]*skillz.CharacterSkill, 0, len(results))
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

func (s *Service) CharacterSkillQueue(ctx context.Context, characterID uint64) ([]*skillz.CharacterSkillQueue, error) {

	key := generateKey(characterSkillQueueKeyPrefix, strconv.FormatUint(characterID, 10))
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkillQueue", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return nil, nil
	}

	queue := make([]*skillz.CharacterSkillQueue, 0, len(results))
	for _, result := range results {
		var position = new(skillz.CharacterSkillQueue)
		err = json.Unmarshal([]byte(result), position)
		if err != nil {
			return queue, errors.Wrapf(err, errorFFormat, skillAPI, "CharacterSkillQueue", "failed to decode json to structure")
		}

		queue = append(queue, position)
	}

	return queue, nil

}

func (s *Service) SetCharacterSkillQueue(ctx context.Context, characterID uint64, positions []*skillz.CharacterSkillQueue, expires time.Duration) error {

	members := make([]string, 0, len(positions))
	for _, position := range positions {
		data, err := json.Marshal(position)
		if err != nil {
			return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkillQueue", "failed to encode struct as json")
		}

		members = append(members, string(data))
	}

	key := generateKey(characterSkillQueueKeyPrefix, strconv.FormatUint(characterID, 10))
	err := s.redis.SAdd(ctx, key, members).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkillQueue", "failed to write cache")
	}

	err = s.redis.Expire(ctx, key, expires).Err()

	return errors.Wrapf(err, errorFFormat, skillAPI, "SetCharacterSkillQueue", "failed to set expiry on set")

}