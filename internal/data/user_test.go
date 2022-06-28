package data_test

import (
	"gchat/internal/biz"
	"gchat/internal/data"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User",

	// 测试函数
	func() {
		var ro biz.UserRepo
		var uD *biz.AddUser
		BeforeEach(func() {
			// 这里的 Db 是 data_suite_test.go 文件里面定义的
			ro = data.NewUserRepo(Db, nil)
			// 这里你可以引入外部组装好的数据
			uD = &biz.AddUser{
				Username: "test",
				Password: "test",
			}
		})

		// 设置 It 块来添加单个规格
		It("CreateUser", func() {
			u, err := ro.AddUser(ctx, uD)
			Ω(err).ShouldNot(HaveOccurred())
			// 组装的数据 mobile 为 13803881388
			Ω(u.Username).Should(Equal("ezvik666111")) // 手机号应该为创建的时候写入的手机号
		})

	},
)
