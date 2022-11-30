package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example/sdks"
	"example/server/database"
	"example/types"
)

var (
	ErrInsufficientCapacity = errors.New("not enough room available")
	ErrNotFound             = errors.New("not found")
)

type Services struct {
	Reservations Reservations
}

func Setup(config types.Config) *Services {
	return &Services{
		Reservations: &reservations{
			db:     database.NewClient(config.DSN),
			hotels: sdks.NewHotelsClient(config.HotelsURL),
			users:  sdks.NewUsersClient(config.UsersURL),
		},
	}
}

type Reservations interface {
	GetReservations(ctx context.Context) ([]types.Reservation, error)
	GetReservationByID(ctx context.Context, id int64) (types.Reservation, error)
	CreateReservation(ctx context.Context, reservation types.CreateReservationBody) (types.Reservation, error)
}

type reservations struct {
	db     database.Client
	hotels sdks.Hotels
	users  sdks.Users
}

func (s *reservations) GetReservations(ctx context.Context) ([]types.Reservation, error) {
	return s.db.SelectReservations(ctx)
}

func (s *reservations) GetReservationByID(ctx context.Context, id int64) (types.Reservation, error) {
	reservation, err := s.db.SelectReservationByID(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return reservation, err
		}
		return reservation, fmt.Errorf("unable to retrieve reservation for id '%d': %w", id, ErrNotFound)
	}
	return reservation, nil
}

func (s *reservations) CreateReservation(ctx context.Context, reservation types.CreateReservationBody) (types.Reservation, error) {
	var res types.Reservation
	user, err := s.users.GetUserByName(ctx, reservation.User)
	if err != nil {
		return res, fmt.Errorf("unable to retrieve user %q: %w", reservation.User, err)
	}
	hotel, err := s.hotels.GetHotelByName(ctx, reservation.Hotel)
	if err != nil {
		return res, fmt.Errorf("unable to retrieve hotel %q: %w", reservation.Hotel, err)
	}
	reservations, err := s.db.SelectReservationsByHotel(ctx, hotel.ID)
	if err != nil {
		return res, fmt.Errorf("unable to retrieve reservations from database: %w", err)
	}

	reservedRooms := int64(0)
	for _, reservation := range reservations {
		reservedRooms += reservation.RoomNumber
	}

	if (reservedRooms + reservation.Rooms) > hotel.Rooms {
		return res, fmt.Errorf("can't create reservation: %w", ErrInsufficientCapacity)
	}

	result, err := s.db.InsertReservation(ctx, types.Reservation{
		UserID:     user.ID,
		HotelID:    hotel.ID,
		RoomNumber: reservation.Rooms,
	})
	if err != nil {
		return res, fmt.Errorf("unable to create reservation on database: %w", err)
	}
	return result, nil
}
