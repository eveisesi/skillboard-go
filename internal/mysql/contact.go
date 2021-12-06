package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type contactRepository struct {
	db       QueryExecContext
	contacts tableConf
}

const (
	ContactContactID   string = "contact_id"
	ContactContactType string = "contact_type"
	ContactStanding    string = "standing"
)

func NewContactRepository(db QueryExecContext) skillz.ContactRepository {

	return &contactRepository{
		db: db,
		contacts: tableConf{
			table: TableCharacterContacts,
			columns: []string{
				ColumnCharacterID, ContactContactID, ContactContactType, ContactStanding, ColumnCreatedAt,
			},
		},
	}

}
func (r *contactRepository) CharacterContacts(ctx context.Context, characterID uint64) ([]*skillz.CharacterContact, error) {

	query, args, err := sq.Select(r.contacts.columns...).
		From(r.contacts.table).
		Where(sq.Eq{ColumnCharacterID: characterID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "CharacterContacts", "failed to generate sql")
	}

	var contacts = make([]*skillz.CharacterContact, 0)
	err = r.db.SelectContext(ctx, &contacts, query, args...)
	return contacts, errors.Wrapf(err, prefixFormat, cloneRepositoryIdentifier, "CharacterContacts")

}

func (r *contactRepository) CreateCharacterContacts(ctx context.Context, contacts []*skillz.CharacterContact) error {

	now := time.Now()
	i := sq.Insert(r.contacts.table).Columns(r.contacts.columns...)
	for _, contact := range contacts {
		contact.CreatedAt = now
		i = i.Values(
			contact.CharacterID,
			contact.ContactID,
			contact.ContactType,
			contact.Standing,
			contact.CreatedAt,
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, contactRepositoryIdentifier, "CreateCharacterContacts", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, cloneRepositoryIdentifier, "CreateCharacterContacts")

}

func (r *contactRepository) DeleteCharacterContacts(ctx context.Context, characterID uint64) error {

	query, args, err := sq.Delete(r.contacts.table).Where(sq.Eq{ColumnCharacterID: characterID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneRepositoryIdentifier, "DeleteCharacterContacts", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}
