package api

import (
	"bytes"
	"database/sql"
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

	testCases := []struct {
		testName   string
		accountId  int64
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			testName:  "OK",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)

			},
		},
		{
			testName:  "NotFound",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			// since the validation happens before calling any function inside the API
			testName:  "BadRequest",
			accountId: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testName:  "InternalError",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			// create a new mock controller to be able to use/build the mock's stubs
			controller := gomock.NewController(t)

			defer controller.Finish()

			store := mockdb.NewMockStore(controller)

			testCase.buildStubs(store)

			// start a server and handle requests using httpserver
			server := NewServer(store)

			// we don't have to start a real server, instead we can use the recorder to catch/send req/res
			recorder := httptest.NewRecorder()

			// TODO: any change in the paths won't be reflected here
			urlPath := fmt.Sprintf("/accounts/%d", testCase.accountId)
			req := httptest.NewRequest(http.MethodGet, urlPath, nil)

			// send the request and capture to the recorder
			server.router.ServeHTTP(recorder, req)
			testCase.checkResp(t, recorder)
		})
	}
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
