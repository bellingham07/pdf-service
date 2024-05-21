gen-api:
	goctl api go -api ./desc/all.api -dir ./ --style "go_zero" --branch="main" --remote="https://gitee.com/hblock/go-zero-tpl"

#goctl 生成表结构对于的gorm代码
gen-model:
	goctl model mysql datasource \
   -url="@tcp(ip)/biye" \
    -table="*" \
     -dir="./api/internal/model" \
     --branch="main" --remote="https://gitee.com/hblock/go-zero-tpl"