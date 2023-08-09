package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/AYehia0/go-bk-mst/db/mock"
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	account := getRandomAccount()

	// create a new mock controller to be able to use/build the mock's stubs
	controller := gomock.NewController(t)

	// invoked once, check to see if all methods that were expected to be called were called.
	defer controller.Finish()

	store := mockdb.NewMockStore(controller)

	store.EXPECT().
		GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start a server and handle requests using httpserver
	server := NewServer(store)

	// we don't have to start a real server, instead we can use the recorder to catch/send req/res
	recorder := httptest.NewRecorder()

	// TODO: any change in the paths won't be reflected here
	urlPath := fmt.Sprintf("/accounts/%d", account.ID)
	req := httptest.NewRequest(http.MethodGet, urlPath, nil)

	// send the request and capture to the recorder
	server.router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func getRandomAccount() db.Account {
	return db.Account{
		ID:        utils.GetRandomAmount(),
		OwnerName: utils.GetRandomOwnerName(),
		Currency:  utils.GetRandomCurrency(),
		Balance:   utils.GetRandomAmount(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// convert the body into account object by unmarshalling
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)

	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}
