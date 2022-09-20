--wrk.method = "GET"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
function request()
  i = math.random(1,10000000)
  path = "/GetInfo?username=admin"..i
  return wrk.format("GET",path)
end

