{{define "navbar"}}

<div>
    <ul class="nav nav-pills">
        <li style="display:none"{{if .IsHome}}class="active"{{end}}><a href="/">首页</a></li>
		<li  {{if .IsContract}}class="active"{{end}}><a href="/contract">客户管理</a></li>
    
	{{if  .CurUser }}
			{{if eq .CurUser.Title "Admin"}}
				<li {{if .IsAddContract}}class="active"{{end}}><a href="/contract/add">添加客户</a></li>
				<li {{if .IsUser}}class="active"{{end}}><a href="/user">帐号管理</a></li>
				<li {{if .IsAddUser}}class="active"{{end}}><a href="/user/add">添加帐号</a></li>
			{{else if ne .CurUser.Title "Manager"}}
				<li {{if .IsAddContract}}class="active"{{end}}><a href="/contract/add">添加客户</a></li>
				<li {{if .IsAddUser}}class="active"{{end}}><a href="/user/add">添加账号</a></li>
			{{else if ne .CurUser.Uname "Guest"}}
				<li {{if .IsUser}}class="active"{{end}}><a href="/user/update?uname={{.CurUser.Uname}}">账户管理</a></li>
			{{end}}
		
			{{if ne .CurUser.Uname "Guest"}}
		        <li><a href="/login?exit=true">退出登录</a></li>
			{{else}}	
				<li><a href="/login">登录</a></li>
			{{end}}
		{{end}}
    </ul>
</div>

{{end}}