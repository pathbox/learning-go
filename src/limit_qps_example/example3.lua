local key = KEYS[1] -- 限流KEY
local limit = tonumber(ARGV[1]) -- 限流大小
local current = tonumber(redis.call('get', key) or "0")

if current + 1 > limit then 
    return 0
else 
    redis.call("INCRBY", key, "1")
    redis.call("EXPIRE", key, "2")
    return 1 
end

-- 每次取一次令牌会进行一次网络开销，至少是毫秒级，并发量很有限
-- 所以每次取一批令牌到客户端，但客户端的令牌数不足达到一个阈值的时候，客户端才向控制层取令牌并且每次取一批