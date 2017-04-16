package subscribers

import (
	"context"
	"fmt"
	"time"

	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	"github.com/lileio/lile/pubsub"
)

type AccountServiceSubscribers struct {
	pubsub.Subscriber
	DB database.Database
}

func (as *AccountServiceSubscribers) AccountCreated(ctx context.Context, i account_service.Account) error {
	fmt.Printf("i = %+v\n", i)
	return nil
}

func (as *AccountServiceSubscribers) Setup(c *pubsub.Client) {
	c.On("account_service.created", as.AccountCreated, 30*time.Second, true)
}
