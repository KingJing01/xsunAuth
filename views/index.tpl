<html lang="en">
<head>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="static/css/style.css">
    <title>冷链物流云平台</title>
</head>
<body>
<script src="static/js/jquery.min.js"></script>
<div class="dl">
    <!--登录头-->
    <div class="dlt">
        <h4>
            <img src="static/img/logo.png" alt="">
        </h4>
    </div>
   <!-- 登录主体-->
    <div class="dlbody">
        <div class="dly">
            <h5></h5>
            <span>
                <i class="iconfont icon-user"></i>
                <input id="username" type="text" name="username">
            </span>
            <span>
                <i class="iconfont icon-mima"></i>
                <input id="password" type="password" name="password">
            </span>
            <span>
                <input type="button" value="登录" id="submit" onclick="login()">
          </span>


            <div style="margin: auto"><span id="msg"></span></div>
        </div>
    </div>
</div>
<script>

    $(document).ready(function () {
//        $('#username').textbox({
//            // readonly: true
//        });
//        $('#password').textbox({
//            // readonly: true
//        });
        $(function () {
            $("#username,#password,#submit").keyup(function (event) {
                code = event.which;
                if (code == 13) {
                    login();
                }
            });
        });
    });

    function login() {
        if (!$("#username").val()) {
            $("#username").focus();
            $("#msg").val("请输入用户名！");
            return false;
        }
        if (!$("#password").val()) {
            $("#password").focus();
            $("#msg").val("请输入密码！");
            return false;
        }
        $.ajax({
            url: "login.html",
            type: 'post',
            dataType: 'json',
            data: {username: $("#username").val(), password: $("#password").val()},
            success: function (data) {
                try {
                    if (data.state == "200") {
                        window.location.href = 'index.jsp';
                    } else if (data.state == "911") {
                        alert(data.msg);
                    }
                } catch (ex) {
                    alert("登录失败,请重试");
                }
            },
            error: function (err) {
                console.error(err);
                //alert("登录失败,服务器网络原因");
            }
        });
    }
</script>
</body>
</html>