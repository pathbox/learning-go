#!/usr/bin/env ruby
  require 'rubygems'
  require 'rest_client'

  def send_mail
    response = RestClient.post "http://api.sendcloud.net/apiv2/mail/send",
    :apiUser => "apiUser",
    :apiKey => "apiKey",
    :from => "mail@domain.com",
    :fromName => "ifaxin客服支持<mail@domain>",
    :to => "test@163.com",
    :subject => "this is a title",
    :html => "你太棒了！你已成功的从SendCloud发送了一封测试邮件，接下来快登录前台去完善账户信息吧！",
    :respEmailId => "true"
    return response
  end

  response = send_mail
  puts response.code
  puts response.to_str