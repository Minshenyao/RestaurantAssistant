# 餐厅助手

--- 

> *ps：一个在校大学生学习Java语言时书上的案例，用go的web方法重构一遍，按照课程进度更新*
> 
Go语言编译器版本、使用框架及其数据库版本`go1.23.1 + gin v1.10.0 + mysql Ver 9.0.1`

开发平台`macOS`

### 安装依赖
`go get -u github.com/gin-gonic/gin`

`go get -u github.com/go-sql-driver/mysql`

### 数据库部分
创建数据库: `create database Canteen;`

创建特殊菜品表: `create table special_dishes (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100) NOT NULL, quantily INT NOT NULL, price DECIMAL(10,2) NOT NULL, date DATE NOT NULL);`

创建用户信息表: `CREATE TABLE User (phone VARCHAR(11) NOT NULL PRIMARY KEY, vip TINYINT(1) NOT NULL DEFAULT 0, deposit DECIMAL(10, 2) NOT NULL DEFAULT 0);`


