# golang-mysql

## 连接数据库

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/tutorial")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("succ open mysql")
}
```

在这里例子里，我们需要注意以下几件事：
1. ```sql.Open()```函数的第一个参数是驱动名。这个字符串是驱动向```database/sql```注册，且一般与包名相同以防混淆。举例来说，如上例中的mysql来自于[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql).
2. 第二个参数是特定于驱动程序的语法，它告诉驱动程序如何访问底层数据存储。在本例中，我们连接到本地MySQL服务器实例中的“tutorial”数据库。
3. 您应该（几乎）始终检查和处理所有```database/sql```操作返回的错误。有一些特殊情况，我们将在后面讨论，在这些情况下这样做是没有意义的。
4. 我们常用```defer db.Close``` 来关闭```sql.DB```，将其控制在函数生命范围内。

值得注意的是，```sql.Open()```并不建立任何连接到数据库，相反，它只需准备数据库抽象以供以后使用。
第一次需要时，将懒惰地建立与基础数据存储的第一个实际连接。如果要立即检查数据库是否可用并可以访问（例如，检查是否可以建立网络连接并登录），请使用```db.ping()```来执行此操作，并记住检查错误：
```go
err = db.Ping()
if err != nil {
	// do something here
}
```


## 从mysql中查数据

我们常见有以下几种操作来从数据库中查询数据

1. 通过一个查询返回若干行。
2. 准备一个查询语句，并重复使用多次，最后销毁该语句。
3. 以一次性的方式执行语句，而无需重复使用。
4. 执行一个查询语句返回一行。

GO的```database/sql```函数名称很重要。如果函数名称包含查询，则会向数据库提出查询，即使是空的，也会返回一组行。
不返回行的语句不应使用查询功能；他们应该使用```exec()```函数。

### 从数据库中获取数据
```go
func QueryUser(db *sql.DB) {
	var (
		id   string
		name string
	)
	rows, err := db.Query("select id,name from user where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
```

下面是上述代码中发生的情况：
1. 我们使用```db.Query()```向数据库发送查询。同时，我们检查其中发生的错误。
2. 使用```rows.Close()```,这非常重要。
3. 我们```rows.Next()```对返回的多行进行迭代。
4. 我们通过```rows.Scan()```将每一行的结果输入到参数中。
5. 我们在结束迭代行后进行了错误处理。

其中有几个部分非常容易出错并带来一系列不好的结果
1. 在```for rows.Next()```循环中，要做好错误处理，不要假设循环会一直进行到最后一句。
2. 只要有一个打开的查询结果(如```rows```)，底层的连接就在被使用且无法呗其他查询使用。这意味着当前链接无法在连接池中使用。如果```rows.Next()```顺利迭代，最终到最后一行遇到internal EOF错误会调用```rows.Close()```关闭。但如果有错误让你提前退出loop循环并提早返回，那么```rows```便没有正常关闭，这是一种很容易耗尽资源的方式。
3. ```row.Close()```是一种无害化操作，所以你可以多次调用它。请注意，我们应该先处理错误，后调用Close()，避免runtime panic。
4. 不要在循环中调用```defer```语句，因为它在函数结束后才执行。假如你这样做，会缓慢消耗内存资源。如果你重复的查询和消费结果在一个循环，你应该显示的调用```rows.Close()```,而不是使用```defer```


### ```Scan()```是如何工作的

当您迭代行并将其扫描到目标变量中时，GO会隐性执行数据类型转换。它基于目标变量的类型。意识到这一点可以清理您的代码，并帮助避免重复工作。
比如表里一行是varchar(45),但内容总是数字。如果我们传一个string的指针，GO会吧bytes转成string，然后我们在用```strconv.ParseInt()```将其转成数字，这很麻烦。
我们可以这样做，传一个int指针，GO会帮我们调用```strconv.ParseInt()```，如果发生错误，```Scan()```也会返回错误帮我们错误处理，这样我们的代码就会更整洁
，这也是```database/sql```推荐我们做的。

### 准备查询





