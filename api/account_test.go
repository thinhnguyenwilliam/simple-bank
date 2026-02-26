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
	"time"

	db "github.com/thinhcompany/simple-bank/db/sqlc"
	"github.com/thinhcompany/simple-bank/token"

	"github.com/thinhcompany/simple-bank/util"

	mockdb "github.com/thinhcompany/simple-bank/db/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomAccount(owner string) db.Accounts {
	return db.Accounts{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func TestGetAccountAPI(t *testing.T) {
	user := randomUser()
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
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
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(
					t,
					r,
					maker,
					authorizationTypeBearer,
					"another-user", // ❌ not the account owner
					time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				// ❌ do nothing → no Authorization header
			},
			buildStubs: func(store *mockdb.MockStore) {
				// middleware blocks request before hitting DB
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
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
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
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
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
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

			server := NewTestServer(t, store)

			recorder := httptest.NewRecorder()
			url := "/accounts/" + strconv.FormatInt(tc.accountID, 10)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			tc.setupAuth(t, req, server.tokenMaker)
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

func TestListAccountsAPI(t *testing.T) {
	user := randomUser()

	n := 5
	accounts := make([]db.Accounts, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.Username)
	}

	testCases := []struct {
		name          string
		query         string
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: "?page_id=1&page_size=5",
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(
					t,
					r,
					maker,
					authorizationTypeBearer,
					user.Username,
					time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  user.Username,
					Limit:  5,
					Offset: 0,
				}

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var got []db.Accounts
				err := json.Unmarshal(recorder.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Len(t, got, n)
			},
		},
		{
			name:  "Unauthorized",
			query: "?page_id=1&page_size=5",
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				// no auth
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "InvalidQueryParams",
			query: "?page_id=0&page_size=1000",
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(
					t,
					r,
					maker,
					authorizationTypeBearer,
					user.Username,
					time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t, store)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet,
				"/accounts"+tc.query,
				nil,
			)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	user := randomUser()
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(
					t,
					r,
					maker,
					authorizationTypeBearer,
					user.Username,
					time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    user.Username,
					Currency: account.Currency,
					Balance:  0,
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				// no auth header
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"currency": "INVALID",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(
					t,
					r,
					maker,
					authorizationTypeBearer,
					user.Username,
					time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t, store)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(
				http.MethodPost,
				"/accounts",
				bytes.NewReader(data),
			)
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
