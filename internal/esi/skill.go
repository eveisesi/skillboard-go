package esi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type SkillAPI interface {
	skills
	etags
	modifiers
}

type skills interface {
	GetCharacterAttributes(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*skillz.CharacterAttributes, error)
	GetCharacterSkills(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*skillz.CharacterSkillMeta, error)
	GetCharacterSkillQueue(ctx context.Context, characterID uint64, mods ...ModifierFunc) ([]*skillz.CharacterSkillQueue, error)
}

func (s *Service) GetCharacterAttributes(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*skillz.CharacterAttributes, error) {

	var attributes = new(skillz.CharacterAttributes)
	var out = new(out)
	out.Data = attributes
	endpoint := fmt.Sprintf(endpoints[GetCharacterAttributes], characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	if err == nil {
		attributes.CharacterID = characterID

	}

	return attributes, errors.Wrap(err, "failed to execute request to ESI for Character Attribute Data")

}

func (s *Service) GetCharacterSkills(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*skillz.CharacterSkillMeta, error) {

	var skills = new(skillz.CharacterSkillMeta)
	var out = new(out)
	out.Data = skills
	endpoint := fmt.Sprintf(endpoints[GetCharacterSkills], characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	if err == nil {
		skills.CharacterID = characterID
		for _, skill := range skills.Skills {
			skill.CharacterID = characterID
		}
	}

	return skills, errors.Wrap(err, "failed to execute request to ESI for Character Skill Data")

}

func (s *Service) GetCharacterSkillQueue(ctx context.Context, characterID uint64, mods ...ModifierFunc) ([]*skillz.CharacterSkillQueue, error) {

	var queue = make([]*skillz.CharacterSkillQueue, 0)
	var out = new(out)
	out.Data = &queue
	endpoint := fmt.Sprintf(endpoints[GetCharacterSkillQueue], characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	if err == nil {
		for _, position := range queue {
			position.CharacterID = characterID
		}
	}

	return queue, errors.Wrap(err, "failed to execute request to ESI for Character Skill Queue Data")

}
