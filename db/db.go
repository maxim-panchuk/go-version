package db

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maxim-panchuk/go-version/domain"
)

type Repo struct {
	pool *pgxpool.Pool
}

func newRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{
		pool: pool,
	}
}

var (
	singleton = Repo{}
	once      sync.Once
)

func GetRepo(pool *pgxpool.Pool) *Repo {
	once.Do(func() {
		singleton = *newRepo(pool)
	})
	return &singleton
}

func (r *Repo) FindVersionSettingEntityByBoId(boId int64) (*domain.VersionSettingEntity, error) {
	conn, err := r.pool.Acquire(context.Background())
	defer conn.Release()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	query := `SELECT *
              FROM qvrn_versionsetting
              WHERE boid = $1
              LIMIT 1`

	var e domain.VersionSettingEntity
	err = conn.QueryRow(context.Background(), query, boId).Scan(&e.BoId, &e.UseVersion, &e.VersionSettingId, &e.BoSysName, &e.BoLocalName, &e.BoDescription)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println(err)
			return nil, nil
		}
		return nil, err
	}

	return &e, nil
}

func (r *Repo) FindVersionSettingEntityByBoSysName(boSysName string) (*domain.VersionSettingEntity, error) {
	conn, err := r.pool.Acquire(context.Background())
	defer conn.Release()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	query := `SELECT *
		  FROM qvrn_versionsetting
		  WHERE bosysname=$1 LIMIT 1`

	var e domain.VersionSettingEntity
	err = conn.QueryRow(context.Background(), query, boSysName).Scan(&e.BoId, &e.UseVersion, &e.VersionSettingId, &e.BoSysName, &e.BoLocalName, &e.BoDescription)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &e, nil
}

func (r *Repo) SaveVersionSettingEntity(versionSettingEntity *domain.VersionSettingEntity) error {

	if versionSettingEntity == nil {
		return nil
	}

	conn, err := r.pool.Acquire(context.Background())
	defer conn.Release()

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = conn.Exec(context.Background(),
		"INSERT INTO qvrn_versionsetting (boid, useversion, bosysname, bolocalname, bodescription) "+
			"VALUES ($1, $2, $3, $4, $5)",
		versionSettingEntity.BoId,
		versionSettingEntity.UseVersion,
		versionSettingEntity.BoSysName,
		versionSettingEntity.BoLocalName,
		versionSettingEntity.BoDescription)

	return err
}

func (r *Repo) FindAllVersionSettingEntities() ([]*domain.VersionSettingEntity, error) {
	conn, err := r.pool.Acquire(context.Background())
	defer conn.Release()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	query := `SELECT boid, useversion, versionsettingid, bosysname, bolocalname, bodescription
			  FROM qvrn_versionsetting`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*domain.VersionSettingEntity, 0)

	for rows.Next() {
		var e domain.VersionSettingEntity
		if err := rows.Scan(&e.BoId, &e.UseVersion, &e.VersionSettingId, &e.BoSysName, &e.BoLocalName, &e.BoDescription); err != nil {
			return nil, err
		}
		results = append(results, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repo) DeleteVersionSettingEntity(versionSettingEntity domain.VersionSettingEntity) error {
	conn, err := r.pool.Acquire(context.Background())
	defer conn.Release()

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = conn.Exec(context.Background(), "DELETE FROM qvrn_versionsetting WHERE versionsettingid = $1", versionSettingEntity.VersionSettingId)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) ExistsByBoIdAndObjectId(boId int64, objectId string) (bool, error) {
	conn, err := r.pool.Acquire(context.Background())
	defer conn.Release()

	if err != nil {
		log.Println(err)
		return false, err
	}

	query := "SELECT EXISTS(SELECT 1 FROM qvrn_version WHERE boid=$1 AND objectid=$2)"
	var exists bool
	err = conn.QueryRow(context.Background(), query, boId, objectId).Scan(&exists)
	if err != nil {
		return false, nil
	}
	return exists, nil
}

func (r *Repo) SaveVersionEntity(versionEntity *domain.VersionEntity) error {

	if versionEntity == nil {
		return nil
	}

	conn, err := r.pool.Acquire(context.Background())
	defer conn.Release()

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = conn.Exec(context.Background(),
		"INSERT INTO qvrn_version (boid, versionkey, data, objectid, comment, number, date, login, iscurrent, isfirst) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		versionEntity.BoId,
		versionEntity.VersionKey,
		versionEntity.Data,
		versionEntity.ObjectId,
		versionEntity.ObjectId,
		versionEntity.Comment,
		versionEntity.Number,
		versionEntity.Date,
		versionEntity.Login,
		versionEntity.IsCurrent,
		versionEntity.IsFirst)

	return err
}

func (r *Repo) FindVersionEntityByBoIdAndObjectIdAndMaxNumber(boId int64, objectId string) ([]domain.VersionEntity, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `
        SELECT a.*
        FROM qvrn_version a
        WHERE a.boid=$1 AND a.objectid=$2 AND a.number = (
            SELECT MAX(b.number)
            FROM qvrn_version b
            WHERE b.boid=a.boid AND b.objectid=a.objectid
        )
    `
	rows, err := conn.Query(context.Background(), query, boId, objectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []domain.VersionEntity
	for rows.Next() {
		var v domain.VersionEntity
		err = rows.Scan(&v.VersionId, &v.VersionKey, &v.Data, &v.ObjectId, &v.Comment, &v.Number, &v.Date, &v.Login, &v.IsCurrent, &v.IsFirst)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}
