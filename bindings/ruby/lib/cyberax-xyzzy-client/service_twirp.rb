# Code generated by protoc-gen-twirp_ruby 1.5.0, DO NOT EDIT.
require 'twirp'
require_relative 'service_pb.rb'

module Xyzzy
  class XyzzyDataService < Twirp::Service
    package 'xyzzy'
    service 'XyzzyData'
    rpc :Ping, PingRequest, PingOk, :ruby_method => :ping
  end

  class XyzzyDataClient < Twirp::Client
    client_for XyzzyDataService
  end
end
