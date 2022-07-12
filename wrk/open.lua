-- open.lua
local json = require "json"
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

math.randomseed(os.time())

function init(args) 
    t = {}
    for line in io.lines ('../openList.json') do
        local l = json.decode(line)
        table.insert(t, l)
        -- print(t['uid'])
        -- print(t['eid'])
        -- print(line)
    end
    print(#t)
end

function request()
    -- print(#t)
    id = math.random(#t)
    body = "uid="..t[id]['uid'].."&envelope_id="..t[id]['eid']
    return wrk.format(nil, nil, nil, body)
end

done = function(summary, latency, requests)
    io.write("------------------------------\n")
    for _, p in pairs({50, 90, 95, 99, 99.99}) do
        n = latency:percentile(p)
        io.write(string.format("TP%g: %0.2f ms\n", p, n / 1000))
    end
end