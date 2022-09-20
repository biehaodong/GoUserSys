wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
function request()
  i = math.random(1,10000000)
  path = "/Login"
  body = "username=admin"..i.."&password=123456"
  return wrk.format(nil,path,nil,body)
end