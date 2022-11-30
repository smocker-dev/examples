package services

import (
	"context"
	"reflect"
	"testing"

	"example/types"

	sdksMocks "example/sdks/mocks"
	databaseMocks "example/server/database/mocks"

	"github.com/stretchr/testify/mock"
)

func Test_reservations_CreateReservation(t *testing.T) {
	hotelsClient := sdksMocks.NewHotels(t)
	hotelsClient.On("GetHotelByName", mock.Anything, "hotel1").Return(types.Hotel{
		ID:    1,
		Name:  "hotel1",
		Rooms: 4,
	}, nil)
	usersClient := sdksMocks.NewUsers(t)
	usersClient.On("GetUserByName", mock.Anything, "user1").Return(types.User{
		ID:   1,
		Name: "user1",
	}, nil)
	dbClient := databaseMocks.NewClient(t)
	dbClient.On("SelectReservationsByHotel", mock.Anything, int64(1)).Return([]types.Reservation{}, nil)
	dbClient.On("InsertReservation", mock.Anything, mock.Anything).Return(types.Reservation{
		ID:      1,
		HotelID: 1,
		UserID:  1,
	}, nil)

	service := reservations{
		db:     dbClient,
		hotels: hotelsClient,
		users:  usersClient,
	}
	ctx := context.Background()
	got, err := service.CreateReservation(ctx, types.CreateReservationBody{
		Hotel: "hotel1",
		User:  "user1",
		Rooms: 2,
	})
	if err != nil {
		t.Errorf("unable to create reservation: %v", err)
		return
	}
	expected := types.Reservation{
		ID:      1,
		HotelID: 1,
		UserID:  1,
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("reservations.CreateReservation() = %v+, want %v+", got, expected)
	}
}
