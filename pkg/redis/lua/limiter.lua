local userId = KEYS[1]

local limitPerWindow = tonumber(ARGV[1])
local now = tonumber(ARGV[2])     -- in milliseconds
local window = tonumber(ARGV[3])  -- in milliseconds

local oneWindowAgo = now - window -- in milliseconds

local reqCount = redis.call("ZCOUNT", userId, oneWindowAgo, now)

if reqCount >= limitPerWindow then
    local oldestReqInWindow = redis.call("ZRANGEBYSCORE", userId, oneWindowAgo, now, "LIMIT", 0, 1)
    local timeBeforeNextReq = window + oldestReqInWindow[1] - now
    return {
        0,
        reqCount,
        timeBeforeNextReq
    }
else
    redis.call("ZADD", userId, now, now)
    return {
        1,
        reqCount + 1,
        0
    }
end
