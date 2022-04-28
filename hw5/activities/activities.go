package activities

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/xfiendx4life/gb_back_2_hw/hw5/manager"
	"github.com/xfiendx4life/gb_back_2_hw/hw5/pool"
)

type Activity struct {
	UserId int
	Date   time.Time
	Name   string
}

func (a *Activity) connection(m *manager.Manager, p *pool.Pool, master bool) (*sql.DB, error) {
	s, err := m.ShardById(a.UserId, master)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

func (a *Activity) Create(m *manager.Manager, p *pool.Pool) error {
	c, err := a.connection(m, p, true)
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO "activities" VALUES ($1, $2, $3)`, a.UserId,
		a.Date, a.Name)
	return err
}
func (a *Activity) Read(m *manager.Manager, p *pool.Pool) error {
	c, err := a.connection(m, p, false)
	if err != nil {
		return err
	}
	r := c.QueryRow(`SELECT "date", "name" FROM "activities" WHERE "user_id" =
	$1`, a.UserId)
	return r.Scan(
		&a.Date,
		&a.Name,
	)
}
func (a *Activity) Update(m *manager.Manager, p *pool.Pool) error {
	c, err := a.connection(m, p, true)
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "activities" SET "name" = $3, "date" = $4
	WHERE "user_id" = $1 AND "date" = $2`, a.UserId,
		a.Name, a.Date, a.Name)
	return err
}
func (a *Activity) Delete(m *manager.Manager, p *pool.Pool) error {
	c, err := a.connection(m, p, true)
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "activities" WHERE "user_id" = $1 AND "date" = $2`,
		a.UserId, a.Date)
	return err
}
