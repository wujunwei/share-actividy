package main

var home = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8" />
<title>博弈论小游戏</title>
</head>
<body>
<form action="/add" method="post">

	name:<input type="text" name="name" /><br>
	
	number:<input type="number" max="100" min="0" name="number" /> <br>
	
	<input type="submit" value="提交">

</form>

</body>
</html>
`

var tableHeader = `
<html>
<body>
<center>
<table border="1">
<tr>
<th>序号</th>
<th>名字</th>
<th>数字</th>
<th>时间</th>
</tr>
`

var tableTail = `
</table>
</center>
</body>
</html>
`
