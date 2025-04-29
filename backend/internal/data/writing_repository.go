package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ielts/internal/domain"
	"ielts/internal/filters"
	"ielts/internal/validator"
	"strings"
	"time"
)

type WritingVariantModel struct {
	DB *sql.DB
}

func (e WritingVariantModel) Insert(WritingVariant *domain.WritingVariant) error {
	query := `INSERT INTO writing_variants (name) 
				VALUES ($1)
				RETURNING id, version`

	args := []interface{}{WritingVariant.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&WritingVariant.ID, &WritingVariant.Version)
}
func (e WritingVariantModel) Get(id int64) (*domain.WritingVariant, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, name, version
				FROM writing_variants
				WHERE id = $1`
	var writingVariant domain.WritingVariant

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, id).Scan(
		&writingVariant.ID,
		&writingVariant.Name,
		&writingVariant.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &writingVariant, nil
}

func (e WritingVariantModel) Update(writing_variant *domain.WritingVariant) error {
	query := `UPDATE writing_variants
				SET name = $1, version = version + 1
				WHERE id = $2 AND version = $3
				RETURNING version`

	args := []interface{}{
		writing_variant.Name,
		writing_variant.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&writing_variant.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (e WritingVariantModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM writing_variants
				WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := e.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (e WritingVariantModel) GetAll(mfilters filters.Filters, msearchOptions filters.WritingVariantsSearch) ([]*domain.WritingVariant, filters.Metadata, error) {
	where := []string{}
	args := []interface{}{}
	argPosition := 1

	if msearchOptions.Name != "" {
		where = append(where, fmt.Sprintf(
			"(to_tsvector('simple', name) @@ plainto_tsquery('simple', $%d) OR name ILIKE $%d)", argPosition, argPosition+1))
		args = append(args, msearchOptions.Name, "%"+msearchOptions.Name+"%")
		argPosition += 2
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " OR ")
	}

	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, name, version
	FROM writing_variants
	%s
	ORDER BY %s %s, id ASC 
	LIMIT $%d OFFSET $%d`,
		whereClause,
		mfilters.SortColumn(),
		mfilters.SortDirection(),
		argPosition,
		argPosition+1)

	args = append(args, mfilters.Limit(), mfilters.Offset())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, filters.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	writing_variants := []*domain.WritingVariant{}

	for rows.Next() {
		var writing_variant domain.WritingVariant
		err := rows.Scan(
			&totalRecords,
			&writing_variant.ID,
			&writing_variant.Name,
			&writing_variant.Version,
		)
		if err != nil {
			return nil, filters.Metadata{}, err
		}
		writing_variants = append(writing_variants, &writing_variant)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metadata := filters.CalculateMetadata(totalRecords, mfilters.Page, mfilters.PageSize)

	return writing_variants, metadata, nil
}

func ValidateWriting_variants(v *validator.Validator, writing_variant *domain.WritingVariant) {
	v.Check(writing_variant.Name != "", "Name", "must be provided")
}
