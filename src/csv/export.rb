require 'csv'

CSV.open("export_ruby.csv", "wb") do |line|
  1000000.times.each do |i|
    line << ["Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World,Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World", "Hello World"]
  end
end