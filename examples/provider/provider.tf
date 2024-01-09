provider "listmonk" {
  host     = "https://my-listmonk.example.com"
  username = "username"
  password = "password"
  headers = {
    "CF-Access-Client-Id" : "123123123123123.access",
    "CF-Access-Client-Secret" : "cn9hdf7239gqgc38wehcd08q7h30edhq8eh"
  }
}
