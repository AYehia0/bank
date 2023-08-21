package api

import (
	"bytes"
	"encoding/json"
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

func TestCreateTransfer(t *testing.T) {

	account1 := getRandomAccount()
	account2 := getRandomAccount()

	// setting missing values
	account1.Currency = utils.USD
	account2.Currency = utils.USD

	amount := 10

	testCases := []struct {
		testName   string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResp  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			testName: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        account1.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// we expect the getAccountById to be called twice
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountId: account1.ID,
					ToAccountId:   account2.ID,
					Amount:        int64(amount),
				}

				store.EXPECT().
					TransferTransaction(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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
			urlPath := "/transfers"
			req := httptest.NewRequest(http.MethodPost, urlPath, bytes.NewReader(data))

			// send the request and capture to the recorder
			server.router.ServeHTTP(recorder, req)
			testCase.checkResp(t, recorder)
		})
	}
}
