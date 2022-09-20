--wrk.method = "GET"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
i=0
function request()
  i=i+1
  i = math.random(1,10000000)
  path = "/GetInfo?username=admin"..i
  return wrk.format("GET",path)
end

