require "sidekiq"

Sidekiq::configure_server do |config|
  config.redis = { url: "redis://localhost:6379" }
end

Sidekiq.options[:queues] = ["myqueue"]
Sidekiq.options[:verbose] = true

class MyRubyWorker
  include Sidekiq::Worker

  def perform(str)
    puts "Hello from Sidekiq: #{str}"
  end

end