package repository

import (
	"context"
	"database/sql"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
	"github.com/giffone/forum-security/internal/object/dto"
	"github.com/giffone/forum-security/internal/object/model"
)

type repoTx struct {
	conf   *config.DriverConf
	schema *object.Query
	db     *sql.DB
}

func NewRepoTx(ctx context.Context, new New) Repo {
	return &repoTx{
		conf:   new.Driver(),
		schema: new.Query(),
		db:     new.DataBase(ctx),
	}
}

func (r *repoTx) Create(ctx context.Context, d dto.DTO) (int, object.Status) {
	ctx2, cancel := context.WithTimeout(ctx, config.TimeLimit5s)
	defer cancel()

	// make query for db
	q := d.Create().MakeQuery(r.schema)

	// new transaction for db
	tx, sts := TxBegin(ctx, r.db)
	if sts != nil {
		return 0, sts
	}
	defer TxRollBack(tx)

	// apply query
	res, err := tx.ExecContext(ctx2, q.Query, q.Fields...)
	if err != nil {
		return 0, object.ByCodeAndLog(config.Code500,
			err, "create:")
	}

	// get id of new record
	id, err := res.LastInsertId()
	if err != nil {
		return 0, object.ByCodeAndLog(config.Code500,
			err, "create: last inserted id:")
	}

	// close transaction
	TxCommit(tx)
	return int(id), nil
}

func (r *repoTx) Delete(ctx context.Context, d dto.DTO) object.Status {
	ctx2, cancel := context.WithTimeout(ctx, config.TimeLimit5s)
	defer cancel()

	q := d.Delete().MakeQuery(r.schema)

	tx, sts := TxBegin(ctx, r.db)
	if sts != nil {
		return sts
	}
	defer TxRollBack(tx)

	_, err := tx.ExecContext(ctx2, q.Query, q.Fields...)
	if err != nil {
		return object.ByCodeAndLog(config.Code500,
			err, "delete:")
	}
	// close transaction
	TxCommit(tx)
	return nil
}

func (r *repoTx) GetList(ctx context.Context, m model.Models) object.Status {
	ctx2, cancel := context.WithTimeout(ctx, config.TimeLimit5s)
	defer cancel()

	q := m.GetList().MakeQuery(r.schema)

	tx, sts := TxBegin(ctx2, r.db)
	if sts != nil {
		return sts
	}
	defer TxRollBack(tx)

	rows, sts := QueryTx(ctx2, tx, q.Query, q.Fields)
	if sts != nil {
		return sts
	}

	defer CloseRows(rows)

	sts = Rows(rows, m)
	if sts != nil {
		return sts
	}

	TxCommit(tx)
	return nil
}

func (r *repoTx) GetOne(ctx context.Context, m model.Model) object.Status {
	ctx2, cancel := context.WithTimeout(ctx, config.TimeLimit5s)
	defer cancel()

	q := m.Get().MakeQuery(r.schema)

	tx, sts := TxBegin(ctx2, r.db)
	if sts != nil {
		return sts
	}
	defer TxRollBack(tx)

	row, sts := QueryRowTx(ctx2, tx, q.Query, q.Fields)
	if sts != nil {
		return sts
	}

	sts = Row(row, m)
	if sts != nil {
		return sts
	}

	TxCommit(tx)
	return nil
}

func (r *repoTx) ExportSettings() (*sql.DB, string, *object.Query) {
	return r.db, r.conf.Port, r.schema
}
