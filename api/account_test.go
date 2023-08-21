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
	"github.com/gin-gonic/gin"
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
			server := newTestServer(t, store)

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

func TestCreateAccountAPI(t *testing.T) {
	account := getRandomAccount()
	account.Balance = 0

	testCases := []struct {
		testName   string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			testName: "OK",
			body: gin.H{
				"owner_name": account.OwnerName,
				"currency":   account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.CreateAccountParams{
					OwnerName: account.OwnerName,
					Currency:  account.Currency,
					Balance:   0,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchAccount(t, recorder.Body, account)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			testName: "BadRequest/InvalidPropertyName",
			body: gin.H{
				"owner":    account.OwnerName, // instead of owner_name
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testName: "BadRequest/InvalidCurrency",
			body: gin.H{
				"owner_name": account.OwnerName, // instead of owner_name
				"currency":   "Invalid",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testName: "InternalError",
			body: gin.H{
				"owner_name": account.OwnerName,
				"currency":   account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
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
			server := newTestServer(t, store)

			// creating the json
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			// we don't have to start a real server, instead we can use the recorder to catch/send req/res
			recorder := httptest.NewRecorder()

			// TODO: any change in the paths won't be reflected here
			urlPath := "/accounts"
			req := httptest.NewRequest(http.MethodPost, urlPath, bytes.NewReader(data))

			// send the request and capture to the recorder
			server.router.ServeHTTP(recorder, req)
			testCase.checkResp(t, recorder)
		})
	}
}

func TestGetAccountsAPI(t *testing.T) {

	numAccounts := 5
	accounts := []db.Account{}

	for range make([]struct{}, numAccounts) {
		accounts = append(accounts, getRandomAccount())
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		testName   string
		query      Query
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			testName: "OK",
			query: Query{
				pageID:   1,
				pageSize: numAccounts,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.GetAccountsParams{
					Limit:  int32(numAccounts),
					Offset: 0,
				}
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			testName: "BadRequest",
			query: Query{
				pageSize: numAccounts,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testName: "BadRequest/InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: -numAccounts,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			testName: "BadRequest/InvalidPageId",
			query: Query{
				pageID:   0,
				pageSize: numAccounts,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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
			server := newTestServer(t, store)

			// we don't have to start a real server, instead we can use the recorder to catch/send req/res
			recorder := httptest.NewRecorder()

			// TODO: any change in the paths won't be reflected here
			urlPath := "/accounts"
			req := httptest.NewRequest(http.MethodGet, urlPath, nil)

			// adding the query params
			q := req.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", testCase.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", testCase.query.pageSize))
			req.URL.RawQuery = q.Encode() // encoding

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

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)

	require.NoError(t, err)

	require.Equal(t, accounts, gotAccounts)
}
