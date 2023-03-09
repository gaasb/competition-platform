// Code generated by SQLBoiler 4.14.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testTeams(t *testing.T) {
	t.Parallel()

	query := Teams()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testTeamsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTeamsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Teams().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTeamsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TeamSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTeamsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := TeamExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Team exists: %s", err)
	}
	if !e {
		t.Errorf("Expected TeamExists to return true, but got false.")
	}
}

func testTeamsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	teamFound, err := FindTeam(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if teamFound == nil {
		t.Error("want a record, got nil")
	}
}

func testTeamsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Teams().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testTeamsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Teams().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testTeamsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	teamOne := &Team{}
	teamTwo := &Team{}
	if err = randomize.Struct(seed, teamOne, teamDBTypes, false, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}
	if err = randomize.Struct(seed, teamTwo, teamDBTypes, false, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = teamOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = teamTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Teams().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testTeamsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	teamOne := &Team{}
	teamTwo := &Team{}
	if err = randomize.Struct(seed, teamOne, teamDBTypes, false, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}
	if err = randomize.Struct(seed, teamTwo, teamDBTypes, false, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = teamOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = teamTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func teamBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Team) error {
	*o = Team{}
	return nil
}

func teamAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Team) error {
	*o = Team{}
	return nil
}

func teamAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Team) error {
	*o = Team{}
	return nil
}

func teamBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Team) error {
	*o = Team{}
	return nil
}

func teamAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Team) error {
	*o = Team{}
	return nil
}

func teamBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Team) error {
	*o = Team{}
	return nil
}

func teamAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Team) error {
	*o = Team{}
	return nil
}

func teamBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Team) error {
	*o = Team{}
	return nil
}

func teamAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Team) error {
	*o = Team{}
	return nil
}

func testTeamsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Team{}
	o := &Team{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, teamDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Team object: %s", err)
	}

	AddTeamHook(boil.BeforeInsertHook, teamBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	teamBeforeInsertHooks = []TeamHook{}

	AddTeamHook(boil.AfterInsertHook, teamAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	teamAfterInsertHooks = []TeamHook{}

	AddTeamHook(boil.AfterSelectHook, teamAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	teamAfterSelectHooks = []TeamHook{}

	AddTeamHook(boil.BeforeUpdateHook, teamBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	teamBeforeUpdateHooks = []TeamHook{}

	AddTeamHook(boil.AfterUpdateHook, teamAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	teamAfterUpdateHooks = []TeamHook{}

	AddTeamHook(boil.BeforeDeleteHook, teamBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	teamBeforeDeleteHooks = []TeamHook{}

	AddTeamHook(boil.AfterDeleteHook, teamAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	teamAfterDeleteHooks = []TeamHook{}

	AddTeamHook(boil.BeforeUpsertHook, teamBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	teamBeforeUpsertHooks = []TeamHook{}

	AddTeamHook(boil.AfterUpsertHook, teamAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	teamAfterUpsertHooks = []TeamHook{}
}

func testTeamsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTeamsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(teamColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTeamToOneBracketUsingBracket(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Team
	var foreign Bracket

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, bracketDBTypes, false, bracketColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Bracket struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.BracketID, foreign.ID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Bracket().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	ranAfterSelectHook := false
	AddBracketHook(boil.AfterSelectHook, func(ctx context.Context, e boil.ContextExecutor, o *Bracket) error {
		ranAfterSelectHook = true
		return nil
	})

	slice := TeamSlice{&local}
	if err = local.L.LoadBracket(ctx, tx, false, (*[]*Team)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Bracket == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Bracket = nil
	if err = local.L.LoadBracket(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Bracket == nil {
		t.Error("struct should have been eager loaded")
	}

	if !ranAfterSelectHook {
		t.Error("failed to run AfterSelect hook for relationship")
	}
}

func testTeamToOneSetOpBracketUsingBracket(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Team
	var b, c Bracket

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, teamDBTypes, false, strmangle.SetComplement(teamPrimaryKeyColumns, teamColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, bracketDBTypes, false, strmangle.SetComplement(bracketPrimaryKeyColumns, bracketColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, bracketDBTypes, false, strmangle.SetComplement(bracketPrimaryKeyColumns, bracketColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Bracket{&b, &c} {
		err = a.SetBracket(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Bracket != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Teams[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.BracketID, x.ID) {
			t.Error("foreign key was wrong value", a.BracketID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.BracketID))
		reflect.Indirect(reflect.ValueOf(&a.BracketID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.BracketID, x.ID) {
			t.Error("foreign key was wrong value", a.BracketID, x.ID)
		}
	}
}

func testTeamToOneRemoveOpBracketUsingBracket(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Team
	var b Bracket

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, teamDBTypes, false, strmangle.SetComplement(teamPrimaryKeyColumns, teamColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, bracketDBTypes, false, strmangle.SetComplement(bracketPrimaryKeyColumns, bracketColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetBracket(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveBracket(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Bracket().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Bracket != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.BracketID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.Teams) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testTeamsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testTeamsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TeamSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testTeamsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Teams().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	teamDBTypes = map[string]string{`ID`: `bigint`, `TeamAlias`: `character varying`, `BracketID`: `uuid`}
	_           = bytes.MinRead
)

func testTeamsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(teamPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(teamAllColumns) == len(teamPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, teamDBTypes, true, teamPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testTeamsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(teamAllColumns) == len(teamPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Team{}
	if err = randomize.Struct(seed, o, teamDBTypes, true, teamColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, teamDBTypes, true, teamPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(teamAllColumns, teamPrimaryKeyColumns) {
		fields = teamAllColumns
	} else {
		fields = strmangle.SetComplement(
			teamAllColumns,
			teamPrimaryKeyColumns,
		)
		fields = strmangle.SetComplement(fields, teamGeneratedColumns)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := TeamSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testTeamsUpsert(t *testing.T) {
	t.Parallel()

	if len(teamAllColumns) == len(teamPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Team{}
	if err = randomize.Struct(seed, &o, teamDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Team: %s", err)
	}

	count, err := Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, teamDBTypes, false, teamPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Team struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Team: %s", err)
	}

	count, err = Teams().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
