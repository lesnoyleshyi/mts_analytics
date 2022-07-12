package http

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/mailru/easyjson"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/dto"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/interfaces"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type JWTValidator struct {
	transport interfaces.JWTValidator
	logger    *zap.Logger
}

//easyjson:json
type jwtPayload struct {
	Authorized bool   `json:"authorized,omitempty"`
	UserID     string `json:"user_id"`
	UserName   string `json:"username,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Email      string `json:"email,omitempty"`
	Expires    int64  `json:"expired"`
}

const (
	ACCESS_TOKEN  = `access_token`
	REFRESH_TOKEN = `refresh_token`
	CTX_USER      = `user`
)

var (
	ErrInvalidToken = errors.New("auth service goes wrong: invalid token struct")
	ErrEmptyUserID  = errors.New("auth service goes wrong: token doesn't have 'user_id' field")
)

func NewJWTValidator(transport interfaces.JWTValidator, logger *zap.Logger) JWTValidator {
	return JWTValidator{transport: transport, logger: logger}
}

func (v JWTValidator) Validate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()

		// retrieve tokens from Cookies in client's request
		tokens, err := getTokens(r)
		if err != nil {
			v.logger.Warn("can't get tokens from request", zap.Error(err))
			// TODO presenter writes response to client
			return
		}

		// send client's tokens to auth service, validate them and
		// receive updated tokens or error in case of invalid client's tokens
		updatedTokens, err := v.transport.Validate(ctx, tokens)
		if err != nil {
			v.logger.Warn("can't validate tokens", zap.Error(err))
			// TODO presenter writes response to client
			return
		}

		// retrieve user info from token to pass it through request's context
		accessPayload, err := parseToken(updatedTokens.Access)
		refreshPayload, err := parseToken(updatedTokens.Access)

		// set new cookies to let client update tokens (bad practice)
		newAccess := http.Cookie{
			Name:     ACCESS_TOKEN,
			Value:    updatedTokens.Access,
			Expires:  time.Unix(accessPayload.Expires, 0),
			HttpOnly: true,
		}
		newRefresh := http.Cookie{
			Name:     REFRESH_TOKEN,
			Value:    updatedTokens.Refresh,
			Expires:  time.Unix(refreshPayload.Expires, 0),
			HttpOnly: true,
		}
		http.SetCookie(w, &newAccess)
		http.SetCookie(w, &newRefresh)

		// inject user info into request's context
		useredCtx := context.WithValue(r.Context(), CTX_USER, accessPayload.UserID)

		next.ServeHTTP(w, r.WithContext(useredCtx))
	}

	return http.HandlerFunc(fn)
}

func getTokens(r *http.Request) (dto.TokenPair, error) {
	accessToken, err := r.Cookie(ACCESS_TOKEN)
	if err != nil {
		return dto.TokenPair{}, err
	}
	refreshToken, err := r.Cookie(REFRESH_TOKEN)
	if err != nil {
		return dto.TokenPair{}, nil
	}

	return dto.TokenPair{Access: accessToken.Value, Refresh: refreshToken.Value}, nil
}

func parseToken(token string) (*jwtPayload, error) {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		return nil, ErrInvalidToken
	}

	payloadData := make([]byte, base64.RawStdEncoding.DecodedLen(len(tokenParts[1])))
	_, err := base64.RawStdEncoding.Decode(payloadData, []byte(tokenParts[1]))
	if err != nil {
		return nil, err
	}

	payload := jwtPayload{}
	if err := easyjson.Unmarshal(payloadData, &payload); err != nil {
		return nil, err
	}

	if payload.UserID == "" {
		return nil, ErrEmptyUserID
	}

	return &payload, nil
}
