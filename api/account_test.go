// simple-bank\api\account_test.go

package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	mockdb "simple-bank/db/mock"
	db "simple-bank/db/sqlc"
	"simple-bank/util"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomAccount() db.Accounts {
	return db.Accounts{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Accounts{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Accounts{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)

			recorder := httptest.NewRecorder()
			url := "/accounts/" + strconv.FormatInt(tc.accountID, 10)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchAccount(
	t *testing.T,
	body *bytes.Buffer,
	account db.Accounts,
) {
	var gotAccount db.Accounts
	err := json.Unmarshal(body.Bytes(), &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func TestCreateAccountAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// expected input
	arg := db.CreateAccountParams{
		Owner:    "alice",
		Currency: "USD",
		Balance:  0,
	}

	// mock behavior
	store.EXPECT().
		CreateAccount(gomock.Any(), gomock.Eq(arg)).
		Times(1).
		Return(db.Accounts{
			ID:       1,
			Owner:    "alice",
			Currency: "USD",
			Balance:  0,
		}, nil)

	server := NewServer(store)

	body := map[string]string{
		"owner":    "alice",
		"currency": "USD",
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
}
