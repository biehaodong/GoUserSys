wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
function request()
  i = math.random(1,200)
  path = "/UploadPic"
  body = "username=admin"..i
  return wrk.format(nil,path,nil,body)
end

