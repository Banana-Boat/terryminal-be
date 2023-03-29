package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		Username: "111",
		Password: "12345",
		Age:      20,
		Gender:   UsersGenderFemale,
	}

	_, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		fmt.Println("error: ", err)
	}
	require.NoError(t, err)
}
