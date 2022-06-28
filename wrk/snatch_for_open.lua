-- snatch.lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

local json = require "json"

uidList = {}
eidList = {}

-- requests = 0
-- responses = 0

local threads = {}

function setup(thread)
    table.insert(threads, thread)
end

function init(args)
    requestcount = 0
    responsecount = 0
end

math.randomseed(os.time())
request = function()
    requestcount = requestcount+1
    uid = math.random(100000)
    body = "uid=" .. uid
    uidList [requestcount] = uid
    return wrk.format(nil, nil, nil, body)
 end

response = function(status, headers, body)
    responsecount = responsecount+1 
    if status == 200 then
        -- for k,v in pairs(headers) do
        --     print(k,v )
        -- end
        local t = json.decode(body)
        -- print(t['code'])
        if t['code'] == 0 then
            eidList[responsecount] = t['data']['envelope_id']
            -- print() 
        else
            eidList[responsecount] = 'fuck'
        end
        -- print(headers[0])
        -- print(json.decode(body))
    end
end

function done(summary, latency, requests)
    print(requests)
    for i, thread in ipairs(threads) do
        local req = thread:get("requests")
        -- local uid = thread:get("uid["..req..']')
        -- print(uidList[i], eidList[i])
        local msg = "%d"
        print(msg:format(req))
    end
end
