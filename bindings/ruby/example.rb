require 'rack'

require_relative 'xyzzy/service_twirp.rb'

c = Xyzzy::XyzzyDataClient.new("https://xyzzy.infra.auslr.io/twirp")

resp = c.ping({})
if resp.error
    puts resp.error
else
    puts resp.data
end
