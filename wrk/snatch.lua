-- snatch.lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

math.randomseed(os.time())
request = function()
    uid = math.random(100000)
    body = "uid=" .. uid
    return wrk.format(nil, nil, nil, body)
end

done = function(summary, latency, requests)
    io.write("------------------------------\n")
    for _, p in pairs({50, 90, 95, 99, 99.99}) do
        n = latency:percentile(p)
        io.write(string.format("TP%g: %0.2f ms\n", p, n / 1000))
    end
end
