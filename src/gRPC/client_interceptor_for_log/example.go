package main

func callInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
    invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
    if reqParam, ok := req.(*suji.SayRequest); ok {
        newMsg := strings.Replace(reqParam.Msg, "喜欢", "讨厌", 1)
        req = &suji.SayRequest{Msg: newMsg}
    }
    err := invoker(ctx, method, req, reply, cc, opts...)
    if err != nil {
        log.Println("接口调用出错", method, err)
        return err
    }
    if replyParam, ok := reply.(*suji.SayReply); ok {
        newMsg := strings.Replace(replyParam.Msg, "讨厌", "喜欢", 1)
        replyParam.Msg = newMsg
    }
    return nil
}
func main() {}

func LinkSujiServer(target string) suji.SujiClient {
    conn, err := grpc.DialContext(context.Background(), target, grpc.WithInsecure(),
        grpc.WithUnaryInterceptor(callInterceptor)) // 在这里加入拦截器
    if err != nil {
        log.Fatalln("链接至服务出错", err, target)
    }
    return suji.NewSujiClient(conn)
}