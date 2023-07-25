# 框架设计文档



框架参考了以下框架及部分组织内的以往代码经验

- [GoFrame](https://goframe.org/pages/viewpage.action?pageId=1114119)

- [go-zero](https://go-zero.dev/)

## 目录设计

### 设计目的

XLab 目前的后端架构是从 MVC 模式演变而来的。

以一份古老代码的项目目录为例:

```sh
├─app
│  ├─controller     
│  ├─middleware     
│  └─response       
├─model
└─utils
```

对于一个 web 后端来说，我们一般把数据库相关连接，模型定义等部分作为 Model 层分离开。而因为 web 后端没有显示界面，我们一般没有一个明显区分的 View 层。至于 Controller 层，虽然我们的项目中有显示命名的`controller`文件夹，但是实际上有部分 Controller 代码是沉淀在`model`目录下的。

教条的 MVC 架构并不能适应 web 开发，在 goFrame 框架中提供了一种替代架构：

<img src="https://goframe.org/download/attachments/3672442/image2021-1-5_19-8-47.png?version=1&modificationDate=1609844927551&api=v2" alt="img" style="zoom:50%;" />

因为我们实际上没有 View 层，而且路由分配，请求解析等类似于 “View” 的功能大部分都是在`controller`目录下实现的，所以将其合并是非常自然的事情。在此之上，goFrame 将业务逻辑代码从`controller`部分抽离出来，将业务逻辑的精简代码从重复的数据解析逻辑中分离。被抽离出来的业务逻辑代码全部归入`model`部分，之后 goFrame 对`model`再次细分了三层，以分离业务逻辑，数据模型和数据库操作。



在实际框架设计中，还有一些 go 语言特性产生的阻碍：

- go 无法处理循环引用
- go 在 import package 时只能以最底层的 package 导出。以上述目录为例，如果`model`和`route`下都有一个名为`user`的 package，那么两者都只能以`user`为名导出。这样首先无法判断来源，可读性降低，其次遇到需要导入两个 package 的情景，就只能进行 alias。

因此我们要小心处理依赖关系，尤其是在处理类似 V->C->M 的依赖关系时小心提取跨模块可用的数据结构。并且要尽量减少嵌套的子目录数量，不能用其表示“子package”。



> 参考：
>
> [代码分层设计 - GoFrame](https://goframe.org/pages/viewpage.action?pageId=3672442)





### 框架目录设计

在考虑上述因素后，本框架的目录设计如下：

```sh
src
├── main.go
├── api           
│   ├── dto       
│   └── route     
├── cmd           
├── internal      
│   ├── controller
│   ├── dao   
│   │   └─ model
│   └── service   
└── utils            
    └─ ...
```

外层的几个 package 中，`cmd`是程序的 booting 入口，`utils` 仍然作为全局的工具包。`api`部分则是仅仅包含了程序保留的端口信息，其以非常结构化的方式声明了每个接口和相应的文档，所有的逻辑部分全部提取到了`internal/controller`中。

`internal`包是 go 提供的特殊属性，其内的 package 不能被第三方引用，这一层是为了多服务/第三方package开发而加的。

内层的 package 中，`controller`完全对应了`api`中的接口，包含相应的业务逻辑。`dao`则是项目的数据库层，对外保留各种数据库操作接口。由`model`包单独声明数据模型，方便被全局引用而不产生冲突（`api`中的`dto`单独分出一个 package 也是同理）。最后有一个可选的`service`层，用于提取可以在项目中复用的逻辑模块。

整个项目的依赖关系如下：

<img src="https://s2.loli.net/2023/07/25/bC8jnSy5rToGaKR.png" alt="image-20230725204443532" style="zoom: 80%;" />

相对于 goFrame 建议的目录结构，我们将`controller`和`api`分的更开，这是为了规范 API 文档和未来可能的代码生成做考虑。

对于最底层也做了充分简化，只保留了`service`和`model`模块，这是出于两方面考虑：

- 目前我们的单个项目都比较简单，连中间的`service`层都难以抽取，如果强行保留多层，只会出现简短的代码重复在每一层写一遍的情况，还会多出很多依赖注入和测试的工作
- 我们主要使用 ORM 来操作数据库。相对于使用 ODBC 或者 预编译sql等方式，ORM 已经提供了很方便的 CRUD 接口了，可以直接从`service`或`controller`目录进行调用。如果有 ORM 无法满足的需求，可以在`DAO`层下继承 ORM 并添加自定义方法。

每个模块的详细设计见下。



## 模块设计

### DAO

==TODO: ent. integration==

```shell
dao
├── model
│	├─ user.go
│	└─ pet.go
├── dao.go
└── complex_crud.go
```

`dao`层中包括了`model`定义和数据库接口的封装，在`dao.go`中有定义如下:

```go
type DBMS struct {
    *gorm.DB
}

var DB = func(ctx context.Context) *DBMS {
    return &DBMS{db.WithContext(ctx)}
}

var db *gorm.DB
```

重新暴露的数据库接口`DBMS`继承自 ORM，理论上可以替换各类 ORM。外层可以通过`DB(ctx)`调用到当前的`DBMS`对象并调用 ORM 原生接口，如需增加 ORM 原有接口不能实现的功能，可以将新方法绑定在`DBMS`上，新增方法可以直接放在`dao`层下（如示例中的`complex_crud.go`）。



#### ORM 原生接口使用规范

因为 ORM 大多已经提供了比较方便的 CRUD 接口，框架设计中，可以从`service`等上层直接调用 ORM 提供的接口。

一般来说，ORM 的作用是免去使用者手拼 SQL 的过程，直接使用预先定义的 Data Model 进行数据库操作，以减少代码成本和维护成本。目前 Go 比较好用的 ORM，如 **ent** 和 **gorm**，都提供了一些更能“自定义”生成 SQL 语句的方式。这些方式有些时候破坏了 ORM 减少维护成本，通过编译来检查 SQL 合法性的目的。

比如，在 gorm 中，可能经常会有这样的调用:

```go
db.Table("users").Where("name = ?",name).First(&user)
```

这样的调用产生了大量的硬编码，实际上，我们完全可以通过以下语句替代上面的查询:

```go
db.Where(&model.User{Name: name}).First(&user)
```

这样的调用才是完全依赖于 model 定义，不会产生额外的技术债的方式。**因为没有强制的代码检查方式，这些规范只能由编写者自觉遵守**。

当然，碍于 ORM 的设计，有些查询我们不能完全脱离 SQL 实现，只能使用 ORM 提供的底层功能，那应该在`dao`下编写新方法，绑定在`DBMS`上供上层调用。



#### Model

数据库的 Model 层是直接和数据库 Schema 进行对应的部分，不同的 ORM 提供了不同的 Model 定义方式，以 gorm 为例:

```go
type User struct {
    ID   int     `json:"id" gorm:"primary_key"`
    Name string  `json:"name"`
}
```

gorm 的 Model 生成很大程度上依赖于“约定”的规则，虽然具体规则可以在 gorm 初始化时设置，但是可自由操作空间较小，只能通过 tag 来设置一些属性。此外，包括 gorm 在内的几乎所有 ORM，外键声明方式都依赖于提前定义的几种 relation，虽然大多数情况足够使用，但是对于一些极端情况的外键自定义比较难以实现。

除此之外，碍于 golang 语言特性，Model 声明还有零值和空值问题，这部分的解决方案详见 [MISC#null](###null)。



需要注意的是，`Model`层的数据定义和`DTO`部分返回给前端的数据模型定义完全分离，在定义`Model`层时理应是不需要考虑接口等事宜的。不过对于一些比较直接的 CRUD 接口，我们可能会将`Model`直接作为返回发送回前端，故在`Model`层的数据定义我们也加上了 json tag。但是像 binding tag 等属于`api`部分的逻辑切忌放在`Model`层，也尽量避免`dto`继承`Model`层的数据模型。



一般来说，ORM 都提供了从数据 Model 直接生成数据库 Schema 的 migrate 接口，不过有时我们希望从数据库接口反向生成 Model 定义。目前 gorm 和 ent 都没有提供足够可靠的工具实现这一点。（==TODO==gorm有gen，不过文档和gorm一样抽象，ent的entimport非官方开发，文档也不够详细，两者的部分功能都不太符合预期）



### API

在实际开发中，前后端联调是发现 API 和预期不一致，或这前端因为后端没有开发好 stuck 住进度是非常常见的痛点。

因此，API 部分是一个项目应该在早期就固定下来的东西，并且应该尽量保证较高的准确性。

```go
api               
├── dto           
│   ├── foo.go    
│   └── general.go
├── init.go       
└── route         
    ├── foo.go    
    └── route.go
```

`dto`部分包括了 API 中的 request/response model 定义，并在`general.go`中提供了可复用的请求解析方法，和统一返回定义：

```go
// general.go

// BindReq bind request, and log info&err.
// Support json & query & header,
// if there is a conflict, the former will overwrite the latter.
func BindReq[T any](c *gin.Context, req *T) error {
    logrus.Debugf("Req.Url: %s, Req.Body: %+v", c.Request.URL, c.Request.Body)
    if err := c.ShouldBindWith(req, GENERAL_BINDER); err != nil {
        return stacktrace.PropagateWithCode(err, http.StatusBadRequest, "Failed to bind request")
    }
    return nil
}

func Response(ctx *gin.Context, code int, msg string, data any) {
    ctx.JSON(http.StatusOK, Resp{
        Code: code,
        Msg:  msg,
        Data: data,
    })
}
```

返回定义中，专门实现了`ResponseFail`，通过解析 [MISC#stacktrace](###stacktrace) 提取 code 和 error message，方便将 API 处理和实际逻辑分离:

```go
func ResponseFail(ctx *gin.Context, err error) {
    logrus.Error(err)

    code := stacktrace.GetCode(err)
    if code == stacktrace.NoCode {
        code = http.StatusInternalServerError
    }
    msg := stacktrace.Current(err).Error()
    Response(ctx, int(code), msg, nil)

}
```

`route`部分注册了路由，形式如下，即**通过一个 wrapper 将具体逻辑分离成格式化的接口**：

```go
///////// route.go  /////////

func SetupRouter(r *gin.RouterGroup) {
    r.GET("/ping", Ping)
    setupFooController(r)
}

///////// foo.go  /////////

func setupFooController(r *gin.RouterGroup) {
    cw := FooCtlWrapper{
        ctl: controller.NewFooController(),
    }
    p := r.Group("/foo")
    p.GET("/get", cw.GetFoo)
}

type FooCtlWrapper struct {
    ctl controller.IFooController
}

// >>>>>>>>>>>>>>>>>> Controller >>>>>>>>>>>>>>>>>>

// GetFoo godoc
//
//	@Summary		get foo
//	@Description	just get a foo
//	@Tags			foo
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.Resp{data=dto.GetFooResp}
//	@Router			/foo/get [get]
func (w FooCtlWrapper) GetFoo(c *gin.Context) {
    var req dto.GetFooReq
    if err := dto.BindReq(c, &req); err != nil {
        dto.ResponseFail(c, err)
        return
    }
    resp, err := w.ctl.GetFoo(c, &req) // 这里是分离出去的接口，输入，返回都是预定义好的 Model
    if err != nil {
        dto.ResponseFail(c, err)
        return
    }
    dto.ResponseSuccess(c, resp)
}

```

所有的声明都较为模板化，不包括任何实际逻辑。

`init.go`为服务初始化和启动代码。



项目使用 OpenAPI 作为 API 文档，并通过代码生成尽量保证 API 的准确性。根据项目具体情况，可能会有“从项目代码生成文档”和“通过文档生成相应代码”两种需求。

从项目代码生成文档的需求，目前框架使用的是[swaggo/swag](https://github.com/swaggo/swag)，其通过代码注释进行接口文档生成。由于其仍需要进行手动维护，为了使接口文档和实际项目尽可能一致，目前 API 部分的目录设计有考虑到尽可能将注释相关信息和对应的代码分类在同一个文件中，方便编写者维护。理论上，可以通过直接分析代码生成文档，无需进行注释。

从文档生成代码的需求，目前有部分可以参考的项目，如[deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen)，但是具体生成格式都与本项目不符。本项目的 API 部分已经特别进行了逻辑分离设计，比较容易通过 OpenAPI 文档进行生成。==TODO==



### Controller

`controller`层一一对应`api`层注册的路由。

在`api`层中，我们将具体逻辑从路由处理中分离出来，以接口的形式依赖注入进`api`层，`controller`层就是对这些接口的实现。

```go
////////// foo.go /////////////

// >>>>>>>>>>>>>>>>>> Interface  >>>>>>>>>>>>>>>>>>

type IFooController interface {
    GetFoo(*gin.Context, *dto.GetFooReq) (*dto.GetFooResp, error)
}

// >>>>>>>>>>>>>>>>>> Controller >>>>>>>>>>>>>>>>>>

// check interface implementation
var _ IFooController = (*FooController)(nil)

var NewFooController = func() *FooController {
    return &FooController{}
}

type FooController struct {
    // maybe some logic config to read from viper
    // or a service dependency
}

// ---------------------- GetFoo ----------------------

func (c *FooController) GetFoo(ctx *gin.Context, req *dto.GetFooReq) (*dto.GetFooResp, error) {
    var resp dto.GetFooResp
    dao.DB(ctx).Model(model.Foo{Name: req.Name}).First(&resp)
    return &resp, nil
}
```

如示例所示，该层应该定义相应的`IXXXController`接口，并通过一个类实现该接口，`api`层可以通过这一层提供的工厂方法注入`IXXXController`。

此处的`XXXController`类中，理论可以添加多个依赖，比如和业务逻辑相关的配置文件读入，或者`service`依赖，只需在工厂方法中进行初始化即可。



### Service

在原来(goFrame)的设计中，`service`层是一层独立的逻辑层，所有`controller`需要调用`service`来操作数据库。

但是考虑到很多`controller`实际上非常简单，徒增依赖注入和考虑如何分离`controller`和`service`逻辑的工作量，本框架将`service`作为可选依赖层。

一般来说，在`controller`中较为重复出现的代码模块，可以在开发过程中逐渐沉淀到`service`层，一些明确较为总要的模块逻辑，也可以直接提取出来。如，在开发中，可能逐渐将用户权限的逻辑提取成一个 service，或者在开发一开始，就将用户管理的功能提取成一个 service，被上层`controller`调用。



## MISC

### enums

go 没有 enums，我们用如下方式实现枚举：

```go
type Season string

const (
	Spring    Season = "spring"
	Summer    Season = "summer"
	Fall      Season = "fall"
	Winter    Season = "winter"
)
```

使用`int`作为 base 类型时，go 还提供了`iota`的语法糖。

### stacktrace

go 一层一层传递 err 并处理的过程其实非常离谱，框架预期会 fork [palantir/stacktrace](https://github.com/palantir/stacktrace) 作为错误处理模块。

### 初始化 & 依赖注入

项目采用 `cobra+viper`进行初始化和参数设置。

`cmd`是项目的 entry point，`main.go`只负责初始化`cmd`。



项目的初始化分为两种方式：

- 显示初始化

  ```go
  // init db
  dao.InitDB()
  ```

- 隐式初始化（通过 go 的 init 机制）

  ```go
  //////// cmd/inti.go /////////
  import (
      "github.com/spf13/cobra"
      "github.com/spf13/viper"
      // init logger
      _ "go-svc-tpl/utils/logger"
  )
  
  //////// utils/logger ////////
  
  func init() {
      // open file location reporter
      logrus.SetReportCaller(true)
      logrus.SetFormatter(&LogFormatter{})
  }
  ```

但是，由于 config 的读取是通过`cobra.OnInitialize()`设置的，只有当 cobra 相应命令的`Excute`被执行时，才会读取配置。所以对于依赖于配置文件，或者需要手动依赖注入的依赖初始化，可能并不适合隐式初始化。

（比如，上述的 logger 初始化后，还是在 root 命令执行中单独根据配置文件设置了 loglevel）



对于有依赖注入的组件，比如`route`，我们注意到其依赖注入方式是使用下一层的工厂方法，如下：

```go
func setupFooController(r *gin.RouterGroup) {
    cw := FooCtlWrapper{
        ctl: controller.NewFooController(),
    }
    ...
}
```

而在`controller`中，工厂方法被这样声明:

```go
var NewFooController = func() *FooController {
    return &FooController{}
}
```

这种以 lambda 表达式，而非函数定义的方式暴露工厂函数，可以方便我们更改工厂方法，从而使得整个框架组件之间，可以更换依赖，或者对依赖进行 mock 测试。





### null

可以处理 go 的零值问题，json 和 sql tag 完全支持，但是对于 gorm 和 swag 似乎还有一些预期外行为。前者似乎可以通过一些简单的修改解决，后者期望从 swag 反向解决，比如以 json tag 解析结果为准？==TODO==



### 测试

应该就是简单的 test，可能需要用到 [mockery (vektra.github.io)](https://vektra.github.io/mockery/latest/)来生成 mock 依赖

