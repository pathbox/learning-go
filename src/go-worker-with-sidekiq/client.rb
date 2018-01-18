require "sidekiq"

Sidekiq::configure_client do |config|
  config.redis = { url: "redis://localhost:6379" }
end

Sidekiq::Client.push "queue" => "myqueue", "class" => "MyGoWorker", "args" => ["hello"]