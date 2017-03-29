package server

import (
	"context"
	"strconv"
	"testing"

	"github.com/lileio/account_service"
	"github.com/lileio/image_service"
	"github.com/stretchr/testify/assert"
)

func TestImageCreateDelete(t *testing.T) {
	truncate()

	ms := MockImageService{}
	ms.On("StoreSync")
	ms.On("Delete")
	is = ms

	b := []byte("sldflsdfs")

	ctx := context.Background()
	isr := &image_service.ImageStoreRequest{
		Filename: "pic.jpg",
		Data:     b,
	}

	req := &account_service.CreateAccountRequest{
		Account: &account_service.Account{
			Name:  name,
			Email: "alexbarlowis@localhost" + strconv.Itoa(9),
		},
		Password: pass,
		Image:    isr,
	}
	ac, err := as.Create(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, ac.Id)
	assert.Equal(t, len(ac.Images), 1)

	dr := &account_service.DeleteAccountRequest{
		Id: ac.Id,
	}
	_, err = as.Delete(ctx, dr)
	assert.Nil(t, err)

	ms.AssertExpectations(t)
}
