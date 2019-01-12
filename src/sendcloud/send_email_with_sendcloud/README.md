# 从Excel中读取邮箱数据，进行自动发送

### 对Excel文件的格式要求： 现在简单支持 第一张sheet， 邮箱字段要在第一列

### 使用方式

- body.txt 保存邮件正文内容, 使用<p>...</p>来进行格式段落换行,格式要整齐，参考现在的格式
- subject.txt 保存邮件主题
- 将xlsx文件放在目录下，重命名为 email_list.xlsx
- 将要加到邮件末尾的图片，命名为image.jpg，要是jpg图片文件，放到目录下
- email_auto_send 这个文件不能删除

在iTerm中执行命令： ./email_auto_send

