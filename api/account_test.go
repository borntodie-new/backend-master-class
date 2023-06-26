package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/borntodie-new/backend-master-class/db/mock"
	db "github.com/borntodie-new/backend-master-class/db/sqlc"
	"github.com/borntodie-new/backend-master-class/token"
	"github.com/borntodie-new/backend-master-class/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListAccountAPI(t *testing.T) {

	var accounts []db.Account
	user, _ := randomUser(t)
	var lastAccount db.Account
	for i := 0; i < 5; i++ {
		lastAccount = randomGetAccountAPI(user.Username)
		accounts = append(accounts, lastAccount)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	server := newTestServer(t, store)

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore, pageID int32, pageSize int32)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageID:   int32(1),
			pageSize: int32(5),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, pageID int32, pageSize int32) {
				arg := db.ListAccountsParams{
					Owner:  lastAccount.Owner,
					Offset: pageID,
					Limit:  pageSize,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name:     "InternalError",
			pageID:   int32(1),
			pageSize: int32(5),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, pageID int32, pageSize int32) {
				arg := db.ListAccountsParams{
					Owner:  lastAccount.Owner,
					Offset: pageID,
					Limit:  pageSize,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidParams",
			pageID:   int32(0),
			pageSize: int32(500),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, pageID int32, pageSize int32) {
				arg := db.ListAccountsParams{
					Owner:  lastAccount.Owner,
					Offset: pageID,
					Limit:  pageSize,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store, tc.pageID, tc.pageSize)

			// 发起请求
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", tc.pageID, tc.pageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}

}

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomCreateAccountAPI(user.Username)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"currency": account.Currency,
				"owner":    account.Owner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// mock 测试
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  account.Balance,
					Currency: account.Currency,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), arg).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"currency": account.Currency,
				"owner":    account.Owner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// mock 测试
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  account.Balance,
					Currency: account.Currency,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), arg).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidParams",
			body: gin.H{
				"currency": "CNC",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// mock 测试
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	server := newTestServer(t, store)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mock 测试
			tc.buildStubs(store)
			// 发起请求
			recorder := httptest.NewRecorder()
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			// 校验结果
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomGetAccountAPI(user.Username)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	server := newTestServer(t, store)

	testCases := []struct {
		name          string
		accountId     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountId: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// build Stubs
			// 构建代理
			tc.buildStubs(store)

			// start test server and send request
			// 发起请求
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)

			// check response
			// 测试结果
			tc.checkResponse(t, recorder)
		})
	}

}

func randomGetAccountAPI(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func randomCreateAccountAPI(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  0,
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
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
	var gotAccount []db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccount)
}
