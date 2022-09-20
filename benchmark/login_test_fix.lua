#固定用户
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
i=0
function request()
  i=i+1
  path = "/Login"
  body = "username=admin"..i.."&password=123456"
  return wrk.format(nil,path,nil,body)
end