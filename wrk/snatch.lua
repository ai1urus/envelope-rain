-- snatch.lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

math.randomseed(os.time())
request = function()
    uid = math.random(100000)
    body = "uid=" .. uid
    return wrk.format(nil, nil, nil, body)
 end