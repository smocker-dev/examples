package database

import (
	"context"
	"database/sql"

	"example/types"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Client interface {
	SelectReservations(ctx context.Context) ([]types.Reservation, error)
	SelectReservationsByHotel(ctx context.Context, hotelID int64) ([]types.Reservation, error)
	SelectReservationByID(ctx context.Context, id int64) (types.Reservation, error)
	InsertReservation(ctx context.Context, reservation types.Reservation) (types.Reservation, error)
}

type client struct {
	db *bun.DB
}

func NewClient(dsn string) Client {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return &client{db: db}
}

func (c *client) SelectReservations(ctx context.Context) ([]types.Reservation, error) {
	var reservations []types.Reservation
	err := c.db.NewSelect().Model(&reservations).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (c *client) SelectReservationsByHotel(ctx context.Context, hotelID int64) ([]types.Reservation, error) {
	var reservations []types.Reservation
	err := c.db.NewSelect().Model(&reservations).Where("hotel_id = ?", hotelID).Scan(ctx)
	if err != nil {
		return reservations, err
	}
	return reservations, nil
}

func (c *client) SelectReservationByID(ctx context.Context, id int64) (types.Reservation, error) {
	var reservation types.Reservation
	err := c.db.NewSelect().Model(&reservation).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return reservation, err
	}
	return reservation, nil
}

func (c *client) InsertReservation(ctx context.Context, reservation types.Reservation) (types.Reservation, error) {
	_, err := c.db.NewInsert().Model(&reservation).Exec(ctx)
	return reservation, err
}
