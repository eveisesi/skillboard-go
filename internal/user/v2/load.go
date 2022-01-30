package user

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/eveisesi/skillz"
	"github.com/sirupsen/logrus"
)

func (s *Service) LoadCharacter(ctx context.Context, user *skillz.User, entry *logrus.Entry, mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	character, err := s.character.Character(ctx, user.CharacterID)
	if err != nil {
		entry.WithError(err).
			Error("failed to fetch character details")
		mx.Lock()
		defer mx.Unlock()
		user.Errors = append(user.Errors, fmt.Errorf("failed to fetch character details"))
		return
	}

	corporation, err := s.corporation.Corporation(ctx, character.CorporationID)
	if err != nil {
		entry.WithError(err).
			Error("failed to fetch character corporation details")
		mx.Lock()
		defer mx.Unlock()
		user.Errors = append(user.Errors, fmt.Errorf("failed to fetch  character corporation details"))
		return
	}

	character.Corporation = corporation

	if corporation.AllianceID.Valid {
		alliance, err := s.alliance.Alliance(ctx, corporation.AllianceID.Uint)
		if err != nil {
			entry.WithError(err).
				Error("failed to fetch corporation alliance details")
			mx.Lock()
			defer mx.Unlock()
			user.Errors = append(user.Errors, fmt.Errorf("failed to fetch corporation alliance details"))
			return
		}

		corporation.Alliance = alliance
	}

	user.Character = character
}

func (s *Service) LoadAttributes(ctx context.Context, user *skillz.User, entry *logrus.Entry, mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	attributes, err := s.skills.Attributes(ctx, user.CharacterID)
	if err != nil {
		entry.WithError(err).
			Error("failed to fetch character attributes")
		mx.Lock()
		defer mx.Unlock()
		user.Errors = append(user.Errors, fmt.Errorf("failed to fetch character attributes"))
		return
	}

	user.Attributes = attributes
}

func (s *Service) LoadImplants(ctx context.Context, user *skillz.User, entry *logrus.Entry, mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	implants, err := s.clones.Implants(ctx, user.CharacterID)
	if err != nil {
		entry.WithError(err).
			Error("failed to fetch character attributes")
		mx.Lock()
		defer mx.Unlock()
		user.Errors = append(user.Errors, fmt.Errorf("failed to fetch character attributes"))
		return
	}

	user.Implants = implants
}

func (s *Service) LoadSkills(ctx context.Context, user *skillz.User, entry *logrus.Entry, mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	skills, err := s.skills.Skillz(ctx, user.CharacterID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character skillz")
		mx.Lock()
		defer mx.Unlock()
		user.Errors = append(user.Errors, fmt.Errorf("failed to fetch character skillz"))
		return
	}

	user.Skills = skills
}

func (s *Service) LoadSkillGrouped(ctx context.Context, user *skillz.User, entry *logrus.Entry, mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	groups, err := s.skills.SkillsGrouped(ctx, user.CharacterID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character skillz")
		mx.Lock()
		defer mx.Unlock()
		user.Errors = append(user.Errors, fmt.Errorf("failed to fetch character skillz"))
		return
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	user.SkillsGrouped = groups
}

func (s *Service) LoadFlyable(ctx context.Context, user *skillz.User, entry *logrus.Entry, mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	flyable, err := s.skills.Flyable(ctx, user.CharacterID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character skillz")
		mx.Lock()
		defer mx.Unlock()
		user.Errors = append(user.Errors, fmt.Errorf("failed to fetch character flyable ships"))
		return
	}

	user.Flyable = flyable
}

func (s *Service) LoadSkillQueue(ctx context.Context, user *skillz.User, entry *logrus.Entry, mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	summary, err := s.skills.SkillQueue(ctx, user.CharacterID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character skillz")
		mx.Lock()
		defer mx.Unlock()
		user.Errors = append(user.Errors, fmt.Errorf("failed to fetch character skill queue"))
		return
	}

	user.QueueSummary = summary
}

func (s *Service) LoadSkillMeta(ctx context.Context, user *skillz.User, entry *logrus.Entry, mx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	meta, err := s.skills.Meta(ctx, user.CharacterID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch character skillz")
		mx.Lock()
		defer mx.Unlock()
		user.Errors = append(user.Errors, fmt.Errorf("failed to fetch character skill queue"))
		return
	}

	user.Meta = meta
}
