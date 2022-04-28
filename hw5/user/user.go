package user

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/xfiendx4life/gb_back_2_hw/hw5/manager"
	"github.com/xfiendx4life/gb_back_2_hw/hw5/pool"
)

type User struct {
	UserId int
	Name   string
	Age    int
	Spouse int
}

func (u *User) connection(m *manager.Manager, p *pool.Pool) (*sql.DB, error) {
	s, err := m.ShardById(u.UserId)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

func (u *User) Create(m *manager.Manager, p *pool.Pool) error {
	c, err := u.connection(m, p)
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO "users" VALUES ($1, $2, $3, $4)`, u.UserId,
		u.Name, u.Age, u.Spouse)
	return err
}
func (u *User) Read(m *manager.Manager, p *pool.Pool) error {
	c, err := u.connection(m, p)
	if err != nil {
		return err
	}
	r := c.QueryRow(`SELECT "name", "age", "spouse" FROM "users" WHERE "user_id" =
	$1`, u.UserId)
	return r.Scan(
		&u.Name,
		&u.Age,
		&u.Spouse,
	)
}
func (u *User) Update(m *manager.Manager, p *pool.Pool) error {
	c, err := u.connection(m, p)
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "users" SET "name" = $2, "age" = $3, "spouse" = $4
	WHERE "user_id" = $1`, u.UserId,
		u.Name, u.Age, u.Spouse)
	return err
}
func (u *User) Delete(m *manager.Manager, p *pool.Pool) error {
	c, err := u.connection(m, p)
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "users" WHERE "user_id" = $1`, u.UserId)
	return err
}
