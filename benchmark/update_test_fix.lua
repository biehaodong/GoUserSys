wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
i=0
function request()
  i=i+1
  path = "/UpdateNickName"
  body = "username=admin"..i.."&nickname=admin"..i
  return wrk.format(nil,path,nil,body)
end

