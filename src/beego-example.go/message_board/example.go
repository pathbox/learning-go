package main

func main() {

}

beego.Router("/msg/", &controllers.MessageController{},"get:Index")
beego.Router("/msg/list", &controllers.MessageController{},"get:List")
beego.Router("/msg/addmsg", &controllers.MessageController{},"post:AddMsg")
beego.Router("/msg/delmsg", &controllers.MessageController{},"post:DelMsg")


type MessageController struct {
    beego.Controller
}

func (c *MessageController) List() {
	msg := models.LeaveMessage{}
	limit,_ := c.GetInt("limit")
	content := c.Input().Get("content")
	response := ResponseJson{}
	response.Message = "ok"
	response.State = 0
	if messages, err:= msg.GetList(limit,page,content) ; err != nil {
			response.Message = err.Error()
			response.State = 500
	} else {
			response.Data = messages
	}
	c.Data["json"] = response
	c.ServeJSON()
	//c.TplName = "message.tpl"
}

func (c *MessageController)AddMsg()  {
    username := c.GetSession("Username")
    content := c.Input().Get("content")
    id, _ := c.GetInt("id",0)
    response := ResponseJson{}
    response.Message = "ok"
    response.State = 0
    if content == "" {
        response.Message = "留言内容不能为空"
        response.State = 500
        c.Data["json"] = response
        c.ServeJSON()
        return
    }
    if username == "" || username == nil {
        response.Message = "当前用户尚未登录，请先登录"
        response.State = 501
        c.Data["json"] = response
        c.ServeJSON()
        return
    }
    msg := models.LeaveMessage{}
    msg.Content = content
    msg.Id = id
    if id,err :=  msg.SaveMessage(username.(string)); err != nil {
        response.Message = "保存失败，请稍后再试"
        response.State = 503
    } else {
        response.Data = id
    }
    c.Data["json"] = response
    c.ServeJSON()
    return
}

func (c *MessageController)DelMsg() {
    username := c.GetSession("Username")
    id, _ := c.GetInt("id",0);
    response := ResponseJson{}
    response.Message = "ok"
    response.State = 0
    msg := models.LeaveMessage{}
    msg.Id = id
    if username == "" || username == nil {
        response.Message = "当前用户尚未登录，请先登录"
        response.State = 501
        c.Data["json"] = response
        c.ServeJSON()
        return
    }
    if err :=  msg.DelMsg(username.(string)); err != nil {
        response.Message = "删除失败，请稍后再试"
        response.State = 503
    }
    c.Data["json"] = response
    c.ServeJSON()
    return
}

package models

import (
    "errors"
    "github.com/astaxie/beego/orm"
    "log"
    "time"
)

type LeaveMessage struct {
    Id       int
    Uid      int
    Content  string
    Status   int
    CreateAt time.Time `orm:"type(datetime)"`
    UpdateAt time.Time `orm:"type(datetime)"`
}

type MsgData struct {
    Id       int
    Content  string
    Name     string
    CreateAt time.Time
}
type MessageList struct {
    Count int
    List  []MsgData
}

func init() {
    orm.RegisterModel(new(LeaveMessage))
}

/**
  添加留言
 */
func (msg *LeaveMessage) SaveMessage(username string) (int, error) {
    o := orm.NewOrm()
    user := User{Name: username}
    if err := user.GetUserId(); err != nil {
        return 0, err
    }

    msg.Uid = user.Id
    msg.Status = 1
    msg.UpdateAt = time.Now()

    if msg.Id > 0 {
        //需要判断是否是自己的留言 注意 这里读到的可能会覆盖自己的 结构 重新开启一个
        msgr := LeaveMessage{Id:msg.Id,Uid:msg.Uid}
        if err := o.Read(&msgr, "uid", "id"); err != nil {
            log.Printf("update user %v error,error info is %v ，is not yourself \n", msg, err)
            return 0, errors.New("不能修改别人的留言")
        }
        msg.CreateAt = time.Now()
        if num, err := o.Update(msg, "content","update_at"); num == 0 || err != nil {
            log.Printf("update user %v error,error info is %v \n", msg, err)
            return 0, errors.New("保存失败，请稍后再试")
        }
    } else {
        if id, err := o.Insert(msg); err != nil || id <= 0 {
            log.Printf("insert user %v error,error info is %v \n", msg, err)
            return 0, errors.New("保存失败，请稍后再试")
        }
    }
    return msg.Id, nil
}

/**
   搜索留言 分页
   1. 联表查询记录列表
   2. 筛选符合条件的结果
   3. 加入分页
 */

func (msg LeaveMessage) GetList(limit, page int, content string) (MessageList, error) {
    qb, _ := orm.NewQueryBuilder("mysql")
    qb2, _ := orm.NewQueryBuilder("mysql")
    if limit == 0 {
        limit = 20
    }
    offset := 0
    if page > 0 {
        offset = (page - 1) * limit
    }
    qb.Select("count(*) ").
        From("leave_message").
        LeftJoin("user").On("leave_message.uid = user.id")
    if content != "" {
        qb.Where("content like '%" + content + "%' ")
    }

    qb2.Select("user.name,leave_message.id,leave_message.content,leave_message.create_at").
        From("leave_message").
        LeftJoin("user").On("leave_message.uid = user.id")
    if content != "" {
        qb2.Where("content like '%" + content + "%' ")
    }
    qb2.OrderBy("leave_message.id desc").Limit(limit).Offset(offset)

    sqlCount := qb.String()
    sqlRows := qb2.String()

    o := orm.NewOrm()
    var messageList MessageList
    var count []int
    var msgDatas []MsgData
    if num, err := o.Raw(sqlCount).QueryRows(&count); err != nil || num == 0 {
        return MessageList{}, errors.New("查询失败，请稍后再试")
    }
    messageList.Count = count[0]
    if num, err := o.Raw(sqlRows).QueryRows(&msgDatas); err != nil || num == 0 {
        return MessageList{}, errors.New("查询失败，请稍后再试")
    }
    messageList.List = msgDatas
    return messageList, nil
}

/**
  删除留言
 */
func (msg *LeaveMessage) DelMsg(username string) error {
    o := orm.NewOrm()
    user := User{Name: username}
    if err := user.GetUserId(); err != nil {
        return err
    }
    msg.Uid = user.Id
    msg.Status = 1
    msg.CreateAt = time.Now()
    msg.UpdateAt = time.Now()

    if msg.Id > 0 {
        //需要判断是否是自己的留言
        if err := o.Read(msg, "uid", "id"); err != nil {
            log.Printf("delete user %v error,error info is %v ，is not yourself \n", msg, err)
            return errors.New("不能删除别人的留言")
        }
        if num,err := o.Delete(msg,"id");err != nil || num == 0{
            log.Printf("delete user %v error,error info is %v ，is not yourself \n", msg, err)
            return errors.New("删除失败，请稍后再试")
        }
        return nil
    } else {
        return errors.New("请选择你要删除的留言")
    }
}