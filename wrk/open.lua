-- open.lua
local json = require "json"
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

math.randomseed(os.time())

function init(args) 
    t = {}
    for line in io.lines ('envelopeList.json') do
        local l = json.decode(line)
        table.insert(t, l)
        -- print(t['uid'])
        -- print(t['eid'])
        -- print(line)
    end
    -- print(#t)
end

function request()
    -- print(#t)
    id = math.random(#t)
    body = "uid="..t[id]['uid'].."&envelope_id="..t[id]['eid']
    return wrk.format(nil, nil, nil, body)
end

-- function request()
--     -- print(#t)
--     id = math.random(100000)
--     body = "uid="..id.."&envelope_id="..id
--     return wrk.format(nil, nil, nil, body)
-- end