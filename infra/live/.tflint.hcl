plugin "terraform" {
  enabled = false
  version = "0.3.0"
  source  = "github.com/terraform-linters/tflint-ruleset-terraform"
}

config {
  module = true
}
