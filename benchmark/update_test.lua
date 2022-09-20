wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
function request()
  i = math.random(1,10000000)
  path = "/UpdateNickName"
  body = "username=admin"..i.."&nickname=admin"..i
  return wrk.format(nil,path,nil,body)
end

