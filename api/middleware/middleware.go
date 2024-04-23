package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rate-limiter/api/types"
	"rate-limiter/pkg/constants"
	client "rate-limiter/pkg/redis/client"
	"rate-limiter/pkg/redis/lua"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimitLuaSha(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		rdb := client.GetClient()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusBadRequest)
			return
		}

		now := time.Now().UnixMilli()
		luaResp, err := rdb.EvalSha(context.Background(), lua.GetSha(), []string{req.UserID}, constants.MaxReqAllowedPerWindow, now, constants.WindowDuration).Result()

		if err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusInternalServerError)
			return
		}

		parsedResponse, err := lua.ParseResponse(luaResp)

		if err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusInternalServerError)
			return
		}

		log.Printf("user: %-10v | success: %-5v | req made in window: %-5v | time before next req: %-3v", req.UserID, parsedResponse.Success, parsedResponse.Count, parsedResponse.TimeBeforeNextReq)

		ctx := context.WithValue(r.Context(), constants.CtxKey, parsedResponse)
		next(w, r.WithContext(ctx))
	}
}

func RateLimitLua(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		rdb := client.GetClient()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusBadRequest)
			return
		}

		luaScript, err := lua.LoadScript()
		if err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusInternalServerError)
			return
		}

		now := time.Now().UnixMilli()
		luaResp, err := luaScript.Run(context.Background(), rdb, []string{req.UserID}, constants.MaxReqAllowedPerWindow, now, constants.WindowDuration).Result()

		if err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusInternalServerError)
			return
		}

		parsedResponse, err := lua.ParseResponse(luaResp)

		if err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusInternalServerError)
			return
		}

		log.Printf("user: %-10v | success: %-5v | req made in window: %-5v | time before next req: %-3v", req.UserID, parsedResponse.Success, parsedResponse.Count, parsedResponse.TimeBeforeNextReq)

		ctx := context.WithValue(r.Context(), constants.CtxKey, parsedResponse)
		next(w, r.WithContext(ctx))
	}
}

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		rdb := client.GetClient()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusBadRequest)
			return
		}

		now := time.Now().UnixMilli()
		oneWindowAgo := now - constants.WindowDuration

		reqCount, err := rdb.ZCount(context.Background(), req.UserID, strconv.FormatInt(oneWindowAgo, 10), strconv.FormatInt(now, 10)).Result()
		if err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusInternalServerError)
			return
		}

		if reqCount >= constants.MaxReqAllowedPerWindow {
			oldestReqInWindowStr, err := rdb.ZRangeByScore(context.Background(), req.UserID, &redis.ZRangeBy{
				Min:    strconv.FormatInt(oneWindowAgo, 10),
				Max:    strconv.FormatInt(now, 10),
				Offset: 0,
				Count:  1,
			}).Result()
			if err != nil {
				HandleError(w, err.Error(), req.UserID, "?", http.StatusInternalServerError)
				return
			}

			oldestReqInWindowInt, err := strconv.ParseInt(oldestReqInWindowStr[0], 10, 64)
			if err != nil {
				HandleError(w, err.Error(), req.UserID, "?", http.StatusInternalServerError)
				return
			}

			timeBeforeNextReq := constants.WindowDuration + oldestReqInWindowInt - now
			timeBeforeNextReqInSec := timeBeforeNextReq / 1000
			timeBeforeNextReqInSecStr := strconv.FormatInt(timeBeforeNextReqInSec, 10)
			HandleError(w, "Rate limit exceeded", req.UserID, timeBeforeNextReqInSecStr, http.StatusTooManyRequests)
			return
		}

		_, err = rdb.ZAdd(context.Background(), req.UserID, redis.Z{
			Score:  float64(now),
			Member: strconv.FormatInt(now, 10),
		}).Result()
		if err != nil {
			HandleError(w, err.Error(), req.UserID, "?", http.StatusInternalServerError)
			return
		}

		middlewareData := types.Response{
			Success:           true,
			Count:             reqCount,
			TimeBeforeNextReq: "0",
		}

		log.Printf("user: %-10v | success: %-5v | req made in window: %-5v | time before next req: %-3v", req.UserID, middlewareData.Success, middlewareData.Count, middlewareData.TimeBeforeNextReq)

		ctx := context.WithValue(r.Context(), constants.CtxKey, middlewareData)
		next(w, r.WithContext(ctx))
	}
}
