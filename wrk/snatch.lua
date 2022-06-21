-- 生成随机uid 100000以内
-- 每个id抢5次

wrk.method = "POST"
wrk.body   = "uid=123"
wrk.path   = "/snatch"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

local counter = 1
local uid_start = 0
local uid_step = 25000

local threads = {}
-- request = function()
function setup(thread)
    thread:set("id", counter)
    thread:set("uid_start", uid_start)
    uid_start = uid_start + uid_step
    table.insert(threads, thread)
    counter = counter + 1
end

function init(args)
    requests = 0
    responses = 0

    math.randomseed(id)
    local msg = "thread %d created"
    print(msg:format(id))
end

local function random_uid()
    return uid_start + math.random(0, uid_step)
end

function request()
    local body = "uid=%d"
    return wrk.format("POST", wrk.path, wrk.headers, body:format(random_uid()))
end