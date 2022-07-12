package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJWTValidator_parseToken(t *testing.T) {
	type testCase struct {
		rawToken              string
		expectedPayloadStruct jwtPayload
		expectedError         error
	}

	testCases := []testCase{
		{
			rawToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJhdXRob3JpemVzIjpmYWxzZSwidXNlcl9pZCI6IjEiLCJ1c2VybmFtZSI6Im" +
				"xlaGEiLCJmaXJzdF9uYW1lIjoiTGV4ZXkiLCJsYXN0X25hbWUiOiJQb3BvdiIs" +
				"ImVtYWlsIjoibGVoYXBvcEB5YS5ydSIsImV4cGlyZWQiOjExMzYxNzEwNDV9." +
				"KxHReLCA7gnTcBcuDgHWY8_38s2MwCTUURp4dWOmHic",
			expectedPayloadStruct: jwtPayload{
				Authorized: false,
				UserID:     "1",
				UserName:   "leha",
				FirstName:  "Lexey",
				LastName:   "Popov",
				Email:      "lehapop@ya.ru",
				Expires:    time.Date(2006, 1, 2, 3, 4, 5, 0, time.UTC).Unix(),
			},
			expectedError: nil,
		},
		{
			rawToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJhdXRob3JpemVzIjpmYWxzZSwidXNlcl9pZCI6IjIiLCJ1c2VybmFtZSI6In" +
				"BpdG9rIiwiZmlyc3RfbmFtZSI6IlBldHlhIiwibGFzdF9uYW1lIjoiUGV0cm92" +
				"IiwiZW1haWwiOiJwaXRva0BnbWFpbC5jb20iLCJleHBpcmVkIjoxNjU3NjIwNjAwfQ." +
				"L6W1K1yBbGmUdlt93qc78vW2zZODsWmKh2yzI-QopLs",
			expectedPayloadStruct: jwtPayload{
				Authorized: false,
				UserID:     "2",
				UserName:   "pitok",
				FirstName:  "Petya",
				LastName:   "Petrov",
				Email:      "pitok@gmail.com",
				Expires:    time.Date(2022, 7, 12, 10, 10, 0, 0, time.UTC).Unix(),
			},
			expectedError: nil,
		},
		{
			rawToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJhdXRob3JpemVkIjp0cnVlLCJ1c2VybmFtZSI6ImxlaGEiLCJmaXJzdF9uYW" +
				"1lIjoiTGV4ZXkiLCJsYXN0X25hbWUiOiJQb3BvdiIsImVtYWlsIjoibGVoYXBv" +
				"cEB5YS5ydSIsImV4cGlyZWQiOjExMzYxNzEwNDV9." +
				"nAuQEoDVEYztULLGNDv8bRidhipo-Xn_kx9-Xk4rQwU",
			expectedPayloadStruct: jwtPayload{
				Authorized: true,
				UserName:   "leha",
				FirstName:  "Lexey",
				LastName:   "Popov",
				Email:      "lehapop@ya.ru",
				Expires:    time.Date(2006, 1, 2, 3, 4, 5, 0, time.UTC).Unix(),
			},
			expectedError: ErrEmptyUserID,
		},
		{
			rawToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjoiIiwidXNlcm5hbWUiOiJsZW" +
				"hhIiwiZmlyc3RfbmFtZSI6IkxleGV5IiwibGFzdF9uYW1lIjoiUG9wb3YiLCJl" +
				"bWFpbCI6ImxlaGFwb3BAeWEucnUiLCJleHBpcmVkIjoxMTM2MTcxMDQ1fQ." +
				"NnpQ-qDROX6rEwrnmnz5eujcE9AhnwLgidXlEO8CZBg",
			expectedPayloadStruct: jwtPayload{
				Authorized: true,
				UserID:     "",
				UserName:   "leha",
				FirstName:  "Lexey",
				LastName:   "Popov",
				Email:      "lehapop@ya.ru",
				Expires:    time.Date(2006, 1, 2, 3, 4, 5, 0, time.UTC).Unix(),
			},
			expectedError: ErrEmptyUserID,
		},
	}
	for _, tc := range testCases {
		payload, err := parseToken(tc.rawToken)
		assert.Equal(t, tc.expectedError, err)
		if payload != nil {
			assert.Equal(t, tc.expectedPayloadStruct, *payload)
		}
	}
}
