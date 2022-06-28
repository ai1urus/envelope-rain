local eid = "3"
local uid = "2"

local envelope = redis.call("HMGET", "EnvelopeInfo:" .. eid, "uid", "opened", "value")

-- return _uid
-- Ret 1 eid 不存在
if not envelope[1] then
    return 1
end

-- Ret 2 eid 与 uid 不匹配
if envelope[1] ~= uid then 
    return 2
end

-- Ret 3 envelope 已开启
if envelope[2] == "true" then
    return 3
end

-- Ret 0 成功打开
redis.call("HMSET", "EnvelopeInfo:"..uid, "opened", "true") 
ARGV[1] = envelope[3]
return 0