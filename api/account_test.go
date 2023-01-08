package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/JustinWG/simpleBank/token"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/JustinWG/simpleBank/db/mock"
	db "github.com/JustinWG/simpleBank/db/sqlc"
	"github.com/JustinWG/simpleBank/util"
	"github.com/golang/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := map[string]struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"success": {
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var gotAccount db.Account
				err := json.Unmarshal(recorder.Body.Bytes(), &gotAccount)
				if err != nil {
					t.Error("bad_response_body")
				}
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, account, gotAccount)
			},
		},
		"not_found": {
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var gotAccount db.Account
				err := json.Unmarshal(recorder.Body.Bytes(), &gotAccount)
				if err != nil {
					t.Error("bad_response_body")
				}
				require.Equal(t, http.StatusNotFound, recorder.Code)
				require.Equal(t, db.Account{}, gotAccount)
			},
		},
		"internal_error": {
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var gotAccount db.Account
				err := json.Unmarshal(recorder.Body.Bytes(), &gotAccount)
				if err != nil {
					t.Error("bad_response_body")
				}
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				require.Equal(t, db.Account{}, gotAccount)
			},
		},
		"invalid_id_bad_request": {
			accountID: -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var gotAccount db.Account
				err := json.Unmarshal(recorder.Body.Bytes(), &gotAccount)
				if err != nil {
					t.Error("bad_response_body")
				}
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, db.Account{}, gotAccount)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tt.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tt.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tt.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)

			tt.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
