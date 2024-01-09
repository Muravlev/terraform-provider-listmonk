resource "listmonk_template" "example" {
  body    = "<p>Helloworld</p>"
  name    = "test_template"
  subject = "Hello world"
  type    = "tx"
}
