package recover

const panicHTML = `<!DOCTYPE html>
<html>
  <head>
    <title>PANIC! at {{.URL}}</title>
    <style type="text/css">
      html { background: #f0f0f5; }
      body {
                                color: #333;
                                font-family: helvetica neue, lucida grande, sans-serif;
                                line-height: 1.5;
                                margin: 0;
                                text-shadow: 0 1px 0 rgba(255, 255, 255, 0.6);
                        }
      .clearfix::after{
                                clear: both;
                                content: ".";
                                display: block;
                                height: 0;
                                visibility: hidden;
                        }
      .container{ max-width: 100%; margin: 0 auto; }
      .content{ overflow: hidden; padding: 15px 20px; }
      pre.prettyprint{
                                background: #fff !important;
                                border:1px solid #ccc !important;
                                font-family: Menlo,'Bitstream Vera Sans Mono','DejaVu Sans Mono',Monaco,Consolas,monospace;
                                font-size: 12px;
                                line-height: 1.5;
                                margin: 0;
                                padding: 10px !important;
                        }
      .pln{ color:#333 !important }
      @media screen {
        .str{color:#d14 !important}.kwd{color:#333 !important}.com{color:#998 !important}.lit,.typ{color:#458 !important}
        .clo,.opn,.pun{color:#333 !important}.tag{color:navy !important}.atn{color:teal !important}.atv{color:#d14 !important}
        .dec{color:#333 !important}.var{color:teal !important}.fun{color:#900 !important}
      }
      @media print, projection {
        .kwd,.tag,.typ{font-weight:700}.str{color:#060 !important}.kwd{color:#006 !important}
        .com{color:#600 !important;font-style:italic}.typ{color:#404 !important}.lit{color:#044 !important}
        .clo,.opn,.pun{color:#440 !important}.tag{color:#006 !important}.atn{color:#404 !important}.atv{color:#060 !important}
      }
      ol.linenums { margin-top:0; margin-bottom:0 }
      ol.linenums li { background:#fff; }
      li.L0, li.L1, li.L2, li.L3, li.L5, li.L6, li.L7, li.L8 { list-style-type:decimal !important }
      @-webkit-keyframes highlight {
        0%   { background: rgba(220, 30, 30, 0.3); }
        100% { background: rgba(220, 30, 30, 0.1); }
      }
      @-moz-keyframes highlight {
        0%   { background: rgba(220, 30, 30, 0.3); }
        100% { background: rgba(220, 30, 30, 0.1); }
      }
      @keyframes highlight {
        0%   { background: rgba(220, 30, 30, 0.3); }
        100% { background: rgba(220, 30, 30, 0.1); }
      }
      ol.linenums li:nth-child(6) {
        background: rgba(220, 30, 30, 0.1);
        -webkit-animation: highlight 400ms linear 1;
        -moz-animation: highlight 400ms linear 1;
        animation: highlight 400ms linear 1;
      }
      header.exception {
        border-bottom: solid 3px #a33;
        padding: 18px 20px;
        height: 65px;
        min-height: 65px;
        overflow: hidden;
        background-color: #20202a;
        color: #aaa;
        text-shadow: 0 1px 0 rgba(0, 0, 0, 0.3);
        font-weight: 200;
        box-shadow: inset 0 -5px 3px -3px rgba(0, 0, 0, 0.05), inset 0 -1px 0 rgba(0, 0, 0, 0.05);
        -webkit-text-smoothing: antialiased;
      }
      header.exception h2 {
        font-weight: 200;
        font-size: 16px;
        margin: 0;
      }
      header.exception h2,
      header.exception p {
        line-height: 1.4em;
        overflow: hidden;
        white-space: pre;
        text-overflow: ellipsis;
      }
      header.exception h2 strong { font-weight: 700; color: #d55; }
      header.exception p {
        font-weight: 200;
        font-size: 20px;
        color: white;
        margin-bottom: 0;
        margin-top: 10px;
      }
      header.exception:hover { height: auto; z-index: 2; }
      header.exception:hover h2,
      header.exception:hover p {
        padding-right: 20px;
        overflow-y: auto;
        word-wrap: break-word;
        white-space: pre-wrap;
        height: auto;
        max-height: 105px;
      }
      .backtrace { padding: 20px; }
      ul.frames { margin:0;padding: 0; }
      ul.frames li {
        background-color: #f8f8f8;
        background: -webkit-linear-gradient(top, #f8f8f8 80%, #f0f0f0);
        background: -moz-linear-gradient(top, #f8f8f8 80%, #f0f0f0);
        background: linear-gradient(top, #f8f8f8 80%, #f0f0f0);
        box-shadow: inset 0 -1px 0 #e2e2e2;
        border-radius: 3px;
        margin-bottom: 5px;
        padding: 7px 20px;
        overflow: hidden;
      }
      ul.frames .name,
      ul.frames .location {
        overflow: hidden;
        height: 1.5em;
        white-space: nowrap;
        word-wrap: none;
        text-overflow: ellipsis;
      }
      ul.frames .func { color: #a33; }
      ul.frames .location {
        font-size: 0.85em;
        font-weight: 400;
        color: #999;
      }
      ul.frames .line { font-weight: bold; }
      ul.frames li:first-child {
        background: #38a;
        box-shadow: inset 0 1px 0 rgba(0, 0, 0, 0.1), inset 0 2px 0 rgba(255, 255, 255, 0.01), inset 0 -1px 0 rgba(0, 0, 0, 0.1);
      }
      ul.frames li:first-child .name,
      ul.frames li:first-child .func,
      ul.frames li:first-child .location {
        color: white;
        text-shadow: 0 1px 0 rgba(0, 0, 0, 0.2);
      }
      ul.frames li:first-child .location { opacity: 0.6; }
      .trace-info, .sidebar { margin-bottom: 15px; }
      @media (min-width: 55em) {
        .trace-info {
                                        float: left;
                                        margin-bottom: 15px;
                                        width: 55%;
                                }
        .sidebar {
                                        float: left;
                                        margin-right: 4%;
                                        margin-bottom: 15px;
                                        width: 40%;
                                }
      }
      .trace-info {
        background: #fff;
        padding: 6px;
        border-radius: 3px;
        margin-bottom: 2px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.03), 1px 1px 0 rgba(0, 0, 0, 0.05), -1px 1px 0 rgba(0, 0, 0, 0.05), 0 0 0 4px rgba(0, 0, 0, 0.04);
      }
      .trace-info .title {
        background: #f1f1f1;
        box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.3);
        overflow: hidden;
        padding: 6px 10px;
        border: solid 1px #ccc;
        border-bottom: 0;
        border-top-left-radius: 2px;
        border-top-right-radius: 2px;
      }
      .trace-info .title .name,
      .trace-info .title .location {
        font-size: 9pt;
        line-height: 26px;
        height: 26px;
        overflow: hidden;
      }
      .trace-info .title .location {
        background:1px solid #aaaaaa;
        float: left;
        font-weight: bold;
        font-size: 10pt;
      }
      .trace-info .title .name {
        float: right;
        font-weight: 200;
         margin: 0;
      }
    </style>
  </head>
  <body>
    <div class="top">
      <header class="exception">
        <h2><strong>PANIC</strong> <span>at {{.URL}}</span></h2>
        <p>{{.Err}}</p>
      </header>
    </div>
    <section class="backtrace clearfix">
      <nav class="sidebar">
        <ul class="frames">
        {{range .Frames}}
          <li>
            <div class="info">
              <div class="name">
                <strong class="func">{{.Name}}</strong>
              </div>
              <div class="location">
                <span class="filename">{{.File}}</span>, line <span class="line">{{.Line}}</span>
              </div>
            </div>
          </li>
        {{end}}
        </ul>
      </nav>
      <div class="trace-info clearfix">
          <div class="title">
              <h2 class="name">{{.Name}}</h2>
              <div class="location">
          <span class="filename">{{.File}}</span>
        </div>
          </div>
          <div class="code-block">
              <pre class="prettyprint lang-go linenums:{{.StartLine}}">{{.SourceLines}}</pre>
          </div>
      </div>
    </section>
    <script src="https://cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js"></script>
  </body>
</html>`
