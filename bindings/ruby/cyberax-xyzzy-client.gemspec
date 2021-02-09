Gem::Specification.new do |s|
  s.name        = "cyberax-xyzzy-client"
  s.version     = "1.0.0"
  s.summary     = "Design Bureau API Ruby Gem"
  s.authors     = ["Aleksei Besogonov"]
  s.homepage    = "https://github.com/cyberax/xyzzy"
  s.license     = "Unlicense"
  s.files       = ["lib/cyberax-xyzzy-client.rb", "lib/cyberax-xyzzy-client/service_pb.rb", "lib/cyberax-xyzzy-client/service_twirp.rb"]
  s.metadata    = { "github_repo" => "ssh://github.com/cyberax/xyzzy" }

  s.add_runtime_dependency "rack", "~> 1.6"
  s.add_runtime_dependency "twirp", "~> 1.5", ">= 1.5.0"
end
